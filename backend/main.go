package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "golang.org/x/crypto/bcrypt"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "strings"
)

type Blog struct {
    Id           int64  `json:"id"`
    Title        string `json:"title"`
    Contents     string `json:"contents"`
    Description  string `json:"description"`
    CreationDate string `json:"creationDate" db:"creation_date"`
    AuthorId     int64  `json:"authorId"`
    PictureLink  string `json:"pictureLink"`
}

type Author struct {
    Id         int64  `json:"id"`
    Name       string `json:"name"`
    AvatarLink string `json:"avatarLink" db:"avatar_link"`
    Login      string
    Password   string
}

var db *sqlx.DB
var minioClient *minio.Client

const minPasswordLength = 5
const maxPasswordLength = 25
const minLoginLength = 5
const maxLoginLength = 25
const fileSizeLimit = 5 << 20 // 5 MB
const minTitleLength = 1
const maxTitleLength = 75
const minDescriptionLength = 1
const maxDescriptionLength = 75

func getBlogById(id int64) (Blog, error) {
    row := db.QueryRowx("SELECT * FROM blog WHERE id=$1", id)
    var b Blog
    if err := row.StructScan(&b); err != nil {
        return b, err
    }
    return b, nil
}

func getAuthorById(id int64) (Author, error) {
    row := db.QueryRowx("SELECT * FROM author WHERE id=$1", id)
    var a Author
    if err := row.StructScan(&a); err != nil {
        return a, err
    }
    return a, nil
}

func getBlogs(c *gin.Context) {
    id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    b, err := getBlogById(id)
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    a, err := getAuthorById(b.AuthorId)
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": b, "metadata": a})
}

type BlogCard struct {
    Id          int64  `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    PictureLink string `json:"pictureLink"`
}

func getBlogsRange(c *gin.Context) {
    offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    size, err := strconv.ParseInt(c.Query("size"), 10, 64)
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    searchQuery := c.DefaultQuery("search-query", "")
    searchQuery = strings.ToLower(searchQuery)
    var rowsCount int64
    err = db.Get(&rowsCount, "SELECT COUNT(*) FROM blog WHERE title ILIKE $1", "%"+searchQuery+"%")
    if err != nil {
        log.Printf("rows count error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }

    rows, err := db.Queryx("SELECT id, title, description, pictureLink FROM blog WHERE title ILIKE $3 ORDER by id DESC LIMIT $1 OFFSET $2",
        size, offset, "%"+searchQuery+"%")
    if err != nil {
        log.Print(err)
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    bcs := make([]BlogCard, 0)
    for rows.Next() {
        var bc BlogCard
        if err := rows.StructScan(&bc); err != nil {
            log.Print(err)
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        bcs = append(bcs, bc)
    }
    c.JSON(http.StatusOK, gin.H{"data": bcs, "metadata": gin.H{"size": rowsCount}})
}

func uploadImage(fh *multipart.FileHeader, bucketName string, filename string) error {
    cxt := context.Background()
    f, err := fh.Open()
    if err != nil {
        return err
    }
    _, err = minioClient.PutObject(cxt, bucketName, filename, f, fh.Size, minio.PutObjectOptions{})
    return err
}

func validateCredentials(form *multipart.Form) (Author, error) {
    if len(form.Value["login"]) == 0 || len(form.Value["password"]) == 0 {
        return Author{}, errors.New("login or password aren't provided")
    }
    login := form.Value["login"][0]
    password := form.Value["password"][0]
    if len(login) < minLoginLength || len(login) > maxLoginLength ||
        len(password) < minPasswordLength || len(password) > maxPasswordLength {
        return Author{}, errors.New("invalid length login or password")
    }
    var user Author
    err := db.Get(&user, "SELECT * FROM author WHERE login=$1", login)
    if err != nil {
        return Author{}, errors.New("wrong credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return Author{}, errors.New("wrong credentials")
    }
    return user, nil
}

func handleCoverImage(form *multipart.Form) (string, error) {
    if len(form.File["cover-image"]) != 1 {
        return "", errors.New("attach exactly 1 cover image")
    }
    if form.File["cover-image"][0].Size > fileSizeLimit {
        return "", errors.New("files are too big, size limit is 5 MB")
    }
    coverImageHeader := form.File["cover-image"][0]
    contentType := strings.ToLower(coverImageHeader.Header.Get("Content-Type"))
    if contentType != "image/jpg" && contentType != "image/png" && contentType != "image/webp" && contentType != "image/jpeg" {
        return "", errors.New("only jpg, png, webp are allowed")
    }
    ext := strings.ToLower(filepath.Ext(path.Base(coverImageHeader.Filename)))
    filename := uuid.New().String() + ext
    err := uploadImage(coverImageHeader, "blog-cover-images", filename)
    if err != nil {
        return "", err
    }
    return filename, nil
}

func handleBlogText(form *multipart.Form) (*strings.Builder, error) {
    if len(form.File["blog-text"]) != 1 {
        return nil, errors.New("attach exactly 1 text file")
    }
    if form.File["blog-text"][0].Size > fileSizeLimit {
        return nil, errors.New("files are too big, size limit is 5 MB")
    }
    blogContentsFile, err := form.File["blog-text"][0].Open()
    if err != nil {
        return nil, errors.New("invalid blog contents file")
    }
    defer blogContentsFile.Close()
    blogContents := new(strings.Builder)
    _, err = io.Copy(blogContents, blogContentsFile)
    if err != nil {
        return nil, errors.New("couldn't save the blog text")
    }
    return blogContents, nil
}

func handleBlogTitle(form *multipart.Form) (string, error) {
    if len(form.Value["title"]) == 0 {
        return "", errors.New("description or title aren't provided")
    }
    title := form.Value["title"][0]
    if len(title) < minTitleLength || len(title) > maxTitleLength {
        return "", errors.New("title length should be between 1 and 75 characters")
    }
    return title, nil
}

func handleBlogDescription(form *multipart.Form) (string, error) {
    if len(form.Value["description"]) == 0 {
        return "", errors.New("description or title aren't provided")
    }
    title := form.Value["description"][0]
    if len(title) < minDescriptionLength || len(title) > maxDescriptionLength {
        return "", errors.New("title length should be between 1 and 75 characters")
    }
    return title, nil
}

func createBlog(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse request into multipart form"})
        return
    }

    user, err := validateCredentials(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    title, err := handleBlogTitle(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    description, err := handleBlogTitle(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    blogContents, err := handleBlogText(form)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    imageName, err := handleCoverImage(form)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err = db.Exec("INSERT INTO blog (title, contents, description, authorid, picturelink) VALUES ($1, $2, $3, $4, $5)",
        title, blogContents.String(), description, user.Id, imageName)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't save the data of blog"})
        return
    }
    c.JSON(http.StatusOK, gin.H{})
}

func createBucket(bucketName string) error {
    ctx := context.Background()
    b, err := minioClient.BucketExists(ctx, bucketName)
    if err != nil {
        return err
    }
    if !b {
        err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return err
        }
        err = minioClient.SetBucketPolicy(ctx, bucketName,
            `{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Principal": {"AWS": ["*"]},
							"Action": ["s3:GetObject"],
							"Resource": ["arn:aws:s3:::`+bucketName+`/*"]
						}
					]
		}`)
        return err
    }
    return nil
}

func main() {
    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", "db", os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
    var err error
    db, err = sqlx.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("failed to connect to db %v", err)
    }
    defer func(db *sqlx.DB) {
        err := db.Close()
        if err != nil {
            log.Fatalf("Failed to close db: %v", err)
        }
    }(db)

    minioClient, err = minio.New("minio:9000", &minio.Options{
        Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
        Secure: false,
    })
    if err != nil {
        log.Fatalf("failed to connect to minio %v", err)
    }
    err = createBucket("blog-cover-images")
    if err != nil {
        log.Fatalf("failed to create the image bucket %v", err)
    }

    router := gin.Default()
    err = router.SetTrustedProxies([]string{"nginx"})
    if err != nil {
        log.Fatalf("failed to set trusted proxies %v", err)
    }
    router.GET("/api/blogs/:id", getBlogs)
    router.GET("/api/blogs", getBlogsRange)
    router.POST("/api/blogs", createBlog)
    err = router.Run(":8088")
    if err != nil {
        log.Fatalf("Error when running server: %v", err)
    }
}
