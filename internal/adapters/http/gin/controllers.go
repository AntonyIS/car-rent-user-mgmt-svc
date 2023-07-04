package gin

import (
	"net/http"

	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"
	service "github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/services"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	CreateUser(ctx *gin.Context)
	ReadUser(ctx *gin.Context)
	ReadUsers(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type handler struct {
	svc service.UserManagementService
}

func NewGinHandler(svc service.UserManagementService) GinHandler {
	return handler{
		svc: svc,
	}
}

func (h handler) CreateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	res, err := h.svc.CreateUser(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
	return
}

func (h handler) ReadUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := h.svc.ReadUser(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

func (h handler) ReadUsers(ctx *gin.Context) {
	users, err := h.svc.ReadUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
	return
}

func (h handler) UpdateUser(ctx *gin.Context) {
	var user domain.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.svc.UpdateUsers(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

func (h handler) DeleteUser(ctx *gin.Context) {
	if err := h.svc.DeleteUser(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
