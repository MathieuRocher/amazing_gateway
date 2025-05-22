package handler

import (
	"amazing_gateway/internal/adapter/application"
	domain "github.com/MathieuRocher/amazing_domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassGroupHandler struct {
	useCase application.ClassGroupUsecaseInterface
}

func NewClassGroupHandler(useCase application.ClassGroupUsecaseInterface) *ClassGroupHandler {
	return &ClassGroupHandler{useCase: useCase}
}

func (h *ClassGroupHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/class-groups")
	{
		group.GET("", h.GetClassGroups)
		group.POST("", h.CreateClassGroup)
		group.GET(":id", h.GetClassGroupByID)
		group.PUT(":id", h.UpdateClassGroupByID)
		group.DELETE(":id", h.DeleteClassGroupByID)

	}
}

func (h *ClassGroupHandler) GetClassGroups(c *gin.Context) {
	classGroups, _ := h.useCase.FindAll()
	c.JSON(http.StatusOK, gin.H{
		"message": classGroups,
	})
}

func (h *ClassGroupHandler) CreateClassGroup(c *gin.Context) {

	var payload domain.ClassGroup
	err := c.Bind(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	err = h.useCase.Create(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unknow sql error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func (h *ClassGroupHandler) GetClassGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	classGroup, err := h.useCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"classGroup": classGroup})
}

func (h *ClassGroupHandler) UpdateClassGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// est-ce qu'il existe
	classGroup, err := h.useCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "classGroup not found"})
		return
	}

	// je le remplace par le body
	err = c.Bind(&classGroup)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	// je save
	err = h.useCase.Update(classGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "classGroup Updated"})
}

func (h *ClassGroupHandler) DeleteClassGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.useCase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "classGroup not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "classGroup deleted succesfully"})
}
