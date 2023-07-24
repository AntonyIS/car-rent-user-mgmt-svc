package app

import (
	"net/http"

	"github.com/AntonyIS/notlify-user-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-user-svc/internal/core/ports"
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
	svc ports.UserService
}

func NewGinHandler(svc ports.UserService) GinHandler {
	routerHandler := handler{
		svc: svc,
	}

	return routerHandler
}

func (h handler) CreateUser(ctx *gin.Context) {
	var res *domain.User
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.svc.CreateUser(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (h handler) ReadUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.svc.ReadUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h handler) ReadUsers(ctx *gin.Context) {
	users, err := h.svc.ReadUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, users)
}

func (h handler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := h.svc.ReadUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	var res *domain.User
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.svc.UpdateUser(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := h.svc.ReadUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
}
