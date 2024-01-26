package app

import (
	"fmt"
	"net/http"

	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	CreateUser(ctx *gin.Context)
	ReadUser(ctx *gin.Context)
	ReadUsers(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	DeleteAllUsers(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	HealthCheck(ctx *gin.Context)
}

type handler struct {
	svc       ports.UserService
	secretKey string
	logger    ports.LoggingService
}

func NewGinHandler(svc ports.UserService, logger ports.LoggingService, secretKey string) GinHandler {
	routerHandler := handler{
		svc:       svc,
		secretKey: secretKey,
	}

	return routerHandler
}

func (h handler) CreateUser(ctx *gin.Context) {
	res := domain.User{
		About:        "",
		Handle:       "",
		ProfileImage: "",
		Following:    0,
		Followers:    0,
		Articles:     []domain.Article{},
	}
	if err := ctx.ShouldBindJSON(&res); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.svc.CreateUser(&res)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h handler) ReadUser(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	user, err := h.svc.ReadUserWithId(user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h handler) ReadUsers(ctx *gin.Context) {
	users, err := h.svc.ReadUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h handler) UpdateUser(ctx *gin.Context) {
	user_id := ctx.Param("id")
	_, err := h.svc.ReadUserWithId(user_id)
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
	res.UserId = user_id
	user, err := h.svc.UpdateUser(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h handler) DeleteUser(ctx *gin.Context) {
	user_id := ctx.Param("id")
	message, err := h.svc.DeleteUser(user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h handler) Login(ctx *gin.Context) {

	var user domain.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	dbUser, err := h.svc.ReadUserWithEmail(user.Email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(dbUser)
	if dbUser.CheckPasswordHarsh(user.Password) {
		middleware := NewMiddleware(h.svc, h.logger, h.secretKey)
		tokenString, err := middleware.GenerateToken(dbUser.UserId)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("token", tokenString, 3600*24*30, "", "", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"accessToken": tokenString,
		})

		return

	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
}

func (h handler) Logout(ctx *gin.Context) {
	tokenString := ctx.GetHeader("tokenString")

	if tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is missing",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Token invalidated successfuly",
	})

}

func (h handler) DeleteAllUsers(ctx *gin.Context) {
	message, err := h.svc.DeleteAllUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h handler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Server running",
	})
}
