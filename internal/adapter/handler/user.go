package handler

import (
	"amazing_gateway/internal/adapter/application"
	"amazing_gateway/internal/adapter/handler/dto/user"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCase application.UserUsecaseInterface
}

func NewUserHandler(uh application.UserUsecaseInterface) *UserHandler {
	return &UserHandler{useCase: uh}
}

func (h *UserHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	public := rg.Group("/users")
	{
		public.POST("", h.CreateUser)
	}
}
func (h *UserHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	protected := rg.Group("/users")
	{
		protected.GET("", h.GetUsers)
		protected.GET(":id", h.GetUserByID)
		protected.PUT(":id", h.UpdateUserByID)
		protected.DELETE(":id", h.DeleteUserByID)
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, _ := h.useCase.FindAll()
	c.JSON(http.StatusOK, gin.H{
		"message": user.ListFromDomain(users),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input user.CreateUserInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	domainUser, err := input.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
	}

	if err := h.useCase.Create(domainUser, input.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user.FromDomain(domainUser)})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	u, err := h.useCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.FromDomain(u)})
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	u, err := h.useCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input user.UpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	if err := validator.New().Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	if err := h.useCase.Update(u, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.useCase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted succesfully"})
}
