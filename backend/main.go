package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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
	err = db.Get(&rowsCount, "SELECT COUNT(*) FROM blog WHERE LOWER(title) LIKE $1", "%"+searchQuery+"%")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	begin := rowsCount - offset
	end := begin + size

	rows, err := db.Queryx("SELECT id, title, description, pictureLink FROM blog WHERE id >= $1 AND id < $2 AND LOWER(title) LIKE $3",
		begin, end, "%"+searchQuery+"%")
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

func main() {
	connStr := "user=admin password=super_secret_password dbname=blogs sslmode=disable"
	var err error
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/api/blogs/:id", getBlogs)
	router.GET("/api/blogs", getBlogsRange)
	router.Run(":8088")
}
