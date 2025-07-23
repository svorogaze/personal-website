package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) GetBlogsRange(c *gin.Context) {
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
	err = api.db.Get(&rowsCount, "SELECT COUNT(*) FROM blog WHERE title ILIKE $1", "%"+searchQuery+"%")
	if err != nil {
		log.Printf("rows count error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	rows, err := api.db.Queryx("SELECT id, title, description, pictureLink FROM blog WHERE title ILIKE $3 ORDER by id DESC LIMIT $1 OFFSET $2",
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
