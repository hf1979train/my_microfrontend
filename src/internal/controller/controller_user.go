package controller

import (
	"myhomesv/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Login(c *gin.Context)
	SignUp(c *gin.Context)
	RequestSentResetEmail(c *gin.Context)
	ResetPassword(c *gin.Context)
}

// Loginはユーザーのログインを処理
func (bc *budgetController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := bc.bu.Login(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// SignUpは新規ユーザー登録を処理
func (bc *budgetController) SignUp(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	err := bc.bu.SignUp(models.User{
		Username: email,
		Email:    email,
		Password: password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up"})
		// c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

// ResetPasswordEmail送信要求
func (bc *budgetController) RequestSentResetEmail(c *gin.Context) {
	email := c.PostForm("email")
	// リセットトークンを生成
	token, err := bc.bu.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		// c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// リセットトークンをDBに保存
	if err := bc.bu.SaveResetToken(email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		// c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// リセットメールを送信
	if err := bc.bu.SendResetEmail(email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		// c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// ResetPasswordはパスワードリセットを処理
func (bc *budgetController) ResetPassword(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	c_password := c.PostForm("confirm-password")
	if password != c_password {
		// c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}
	err := bc.bu.ResetPassword(token, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Reset Password"})
		// c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
