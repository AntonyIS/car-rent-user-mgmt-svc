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
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}



type handler struct {
	svc       ports.UserService
	secretKey string
}

func NewGinHandler(svc ports.UserService, secretKey string) GinHandler {
	routerHandler := handler{
		svc:       svc,
		secretKey: secretKey,
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

	res.About, res.Handle, res.ProfileImage, res.Followers, res.Following = " ", " ", " ", 0, 0
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
	res.Id = id
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
	id := ctx.Param("id")
	message, err := h.svc.DeleteUser(id)
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

	if dbUser.CheckPasswordHarsh(user.Password) {
		middleware := NewMiddleware(h.svc, h.secretKey)
		tokenString, err := middleware.GenerateToken(dbUser.Id)

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
