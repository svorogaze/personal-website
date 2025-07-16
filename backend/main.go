package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
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
var imageFolder string

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
		log.Print(err)
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

func getBlogsCount() int64 {
	var c int64
	db.Get(&c, "SELECT COUNT(*) FROM blog")
	return c
}
func createBlog(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse into multipart form"})
		log.Print(err)
		return
	}

	if len(form.Value["login"]) == 0 || len(form.Value["password"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password or login aren't provided"})
		log.Print(err)
		return
	}
	login := form.Value["login"][0]
	password := form.Value["password"][0]
	if len(login) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		log.Print(err)
		return
	}

	var user Author
	err = db.Get(&user, "SELECT * FROM author WHERE login=$1", login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		log.Print(err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		log.Print(err)
		return
	}

	if len(form.Value["title"]) == 0 || len(form.Value["description"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description or title aren't provided"})
		log.Print(err)
		return
	}
	title := form.Value["title"][0]
	description := form.Value["description"][0]
	if len(form.File["blog-text"]) != 1 || len(form.File["cover-image"]) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "attach exactly 1 text file and cover image"})
		log.Print(err)
		return
	}

	if form.File["blog-text"][0].Size > 5<<20 || form.File["cover-image"][0].Size > 5<<20 { // 5 MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Files are too big, size limit is 5 MB"})
	}

	blogContentsFile, err := form.File["blog-text"][0].Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog contents file"})
		log.Print(err)
		return
	}
	defer blogContentsFile.Close()
	blogContents := new(strings.Builder)
	_, err = io.Copy(blogContents, blogContentsFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't copy the file"})
		log.Print(err)
		return
	}

	coverImageHeader := form.File["cover-image"][0]
	filename := path.Base(coverImageHeader.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".png" && ext != ".jpg" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image extension isn't valid"})
		return
	}
	imageName := fmt.Sprintf("%d-%d-%s", user.Id, getBlogsCount(), filename)
	err = c.SaveUploadedFile(coverImageHeader, filepath.Join(imageFolder, imageName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't save cover image file"})
		log.Print(err)
		return
	}

	_, err = db.Exec("INSERT INTO blog (title, contents, description, authorid, picturelink) VALUES ($1, $2, $3, $4, $5)",
		title, blogContents.String(), description, user.Id, imageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't insert row into db"})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=blogs sslmode=disable", os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"))
	imageFolder = os.Getenv("IMAGE_FOLDER_PATH")
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/api/blogs/:id", getBlogs)
	router.GET("/api/blogs", getBlogsRange)
	router.POST("/api/blogs", createBlog)
	router.Run(":8088")
}
