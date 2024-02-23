package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AntonyIS/notelify-users-service/config"
	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type GinHandler interface {
	CreateUser(ctx *gin.Context)
	ReadUser(ctx *gin.Context)
	ReadUsers(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	DeleteAllUsers(ctx *gin.Context)
	Login(ctx *gin.Context)
	GithubLogin(ctx *gin.Context)
	GithubCallback(ctx *gin.Context)
	Logout(ctx *gin.Context)
	HealthCheck(ctx *gin.Context)
}

type handler struct {
	svc         ports.UserService
	conf        config.Config
	logger      ports.LoggingService
	githubOauth *oauth2.Config
}

func NewGinHandler(svc ports.UserService, logger ports.LoggingService, conf config.Config) GinHandler {
	oauthConfig := &oauth2.Config{
		ClientID:     conf.GITHUB_CLIENT_ID,
		ClientSecret: conf.GITHUB_CLIENT_SECRET,
		RedirectURL:  conf.GITHUB_REDIRECT_URL,
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	routerHandler := handler{
		svc:         svc,
		conf:        conf,
		logger:      logger,
		githubOauth: oauthConfig,
	}

	return routerHandler
}

func (h handler) CreateUser(ctx *gin.Context) {
	var res domain.User
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusCreated, gin.H{
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

	if dbUser.CheckPasswordHarsh(user.Password) {
		middleware := NewMiddleware(h.svc, h.logger, h.conf.SECRET_KEY)
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

func (h handler) GithubLogin(ctx *gin.Context) {
	url := h.githubOauth.AuthCodeURL("state")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h handler) GithubCallback(ctx *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Println("ERROR 1!", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, err := h.githubOauth.Exchange(context.Background(), request.Code)

	if err != nil {
		fmt.Println("ERROR 2!", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	user, err := getUserDetails(token.AccessToken)
	if err != nil {
		fmt.Println("ERROR 3!", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	if user.GitHubId != "" {
		dbUser, err := h.svc.ReadUserWithGithubId(user.GitHubId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				newUser, err := h.svc.CreateUser(user)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusOK, newUser)
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dbUser)
		return
	}
}

func (h handler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Server running",
	})
}

func getUserDetails(accessToken string) (*domain.User, error) {
	// GitHub API endpoint for authenticated user details
	apiURL := "https://api.github.com/user"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the OAuth token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	var githubUser domain.GithubUser
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&githubUser); err != nil {
		return nil, err
	}
	githubUser.AccessToken = accessToken
	user := githubUser.InitGithubUser()
	return &user, nil
}
