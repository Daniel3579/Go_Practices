package httpapi

import (
	"net/http"
	"strconv"

	"example.com/prc_integr/internal/models"
	"example.com/prc_integr/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Svc *service.Service
}

func (rt Router) Register(r *gin.Engine) {
	r.POST("/notes", rt.CreateNote)
	r.GET("/notes/:id", rt.GetNote)
	r.PUT("/notes/:id", rt.UpdateNote)
	r.DELETE("/notes/:id", rt.DeleteNote)
	r.GET("/notes", rt.ListNotes)
}

func (rt Router) CreateNote(c *gin.Context) {
	var in struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	n := models.Note{Title: in.Title, Content: in.Content}
	if err := rt.Svc.Create(c, &n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, n)
}

func (rt Router) GetNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	n, err := rt.Svc.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, n)
}

func (rt Router) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var in struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	n := models.Note{ID: id, Title: in.Title, Content: in.Content}
	if err := rt.Svc.Update(c, &n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (rt Router) DeleteNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := rt.Svc.Delete(c, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (rt Router) ListNotes(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	notes, err := rt.Svc.List(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}
