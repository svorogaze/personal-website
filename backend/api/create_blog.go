package api

import (
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "path"
    "path/filepath"
    "strings"
)

func (api *API) validateCredentials(form *multipart.Form) (Author, error) {
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
    err := api.db.Get(&user, "SELECT * FROM author WHERE login=$1", login)
    if err != nil {
        return Author{}, errors.New("wrong credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return Author{}, errors.New("wrong credentials")
    }
    return user, nil
}

func (api *API) handleCoverImage(form *multipart.Form) (string, error) {
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
    err := api.uploadImage(coverImageHeader, "blog-cover-images", filename)
    if err != nil {
        return "", err
    }
    return filename, nil
}

func (api *API) handleBlogText(form *multipart.Form) (*strings.Builder, error) {
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
    defer func(blogContentsFile multipart.File) {
        err := blogContentsFile.Close()
        if err != nil {
            log.Printf("failed to close the file: %v", err)
        }
    }(blogContentsFile)
    blogContents := new(strings.Builder)
    _, err = io.Copy(blogContents, blogContentsFile)
    if err != nil {
        return nil, errors.New("couldn't save the blog text")
    }
    return blogContents, nil
}

func (api *API) handleBlogTitle(form *multipart.Form) (string, error) {
    if len(form.Value["title"]) == 0 {
        return "", errors.New("description or title aren't provided")
    }
    title := form.Value["title"][0]
    if len(title) < minTitleLength || len(title) > maxTitleLength {
        return "", errors.New("title length should be between 1 and 75 characters")
    }
    return title, nil
}

func (api *API) handleBlogDescription(form *multipart.Form) (string, error) {
    if len(form.Value["description"]) == 0 {
        return "", errors.New("description or title aren't provided")
    }
    title := form.Value["description"][0]
    if len(title) < minDescriptionLength || len(title) > maxDescriptionLength {
        return "", errors.New("title length should be between 1 and 75 characters")
    }
    return title, nil
}

func (api *API) CreateBlog(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse request into multipart form"})
        return
    }

    user, err := api.validateCredentials(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    title, err := api.handleBlogTitle(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    description, err := api.handleBlogTitle(form)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    blogContents, err := api.handleBlogText(form)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    imageName, err := api.handleCoverImage(form)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err = api.db.Exec("INSERT INTO blog (title, contents, description, authorid, picturelink) VALUES ($1, $2, $3, $4, $5)",
        title, blogContents.String(), description, user.Id, imageName)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't save the data of blog"})
        return
    }
    c.JSON(http.StatusOK, gin.H{})
}
