package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (api *API) GetBlog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	b, err := api.getBlogById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	a, err := api.getAuthorCardById(b.AuthorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": b, "metadata": a})
}
