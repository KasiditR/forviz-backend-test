package handlers

import (
	"net/http"

	"github.com/KasiditR/forviz-backend-api-test/internal/database"
	"github.com/KasiditR/forviz-backend-api-test/internal/models"
	"github.com/KasiditR/forviz-backend-api-test/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			return
		}

		if user.Username == nil || user.Password == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
			return
		}

		usernameCount, err := database.CountDocument(bson.M{"username": user.Username}, user)
		if err != nil || usernameCount > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Username already exist"})
			return
		}

		password := utils.HashPassword(*user.Password)
		user.ID = bson.NewObjectID()
		user.Password = &password

		_, err = database.InsertOne(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if input.Username == "" || input.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
			return
		}
		var user models.User
		err := database.FindOne(bson.M{"username": input.Username}, &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ok, msg := utils.VerifyPassword(input.Password, *user.Password)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}

		accessToken, refreshToken, err := utils.TokenGenerator(user.ID.Hex(), *user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"user_id":      user.ID.Hex(),
		})
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var refreshTokenRequest struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		payload, msg := utils.ValidateRefreshToken(refreshTokenRequest.RefreshToken)
		if msg != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		token, _, err := utils.TokenGenerator(payload.ID, payload.UserName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"accessToken": token})
	}
}
