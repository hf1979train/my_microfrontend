package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"myhomesv/internal/domain/models"
	"myhomesv/pkg/utils"
	"os"

	"github.com/joho/godotenv"
)

// AuthUsecaseは認証に関するユースケースのインターフェース
type IAuthUsecase interface {
	Login(username, password string) (string, error)
	SignUp(user models.User) error
	GenerateResetToken() (string, error)
	SaveResetToken(email, token string) error
	SendResetEmail(email, token string) error
	ResetPassword(token string, password string) error
}

// Loginはユーザーのログイン処理を実行
func (bu *BudgetUsecase) Login(username, password string) (string, error) {
	user, err := bu.br.FindByUsername(username)
	if err != nil || user.Password != password {
		return "", errors.New("invalid username or password")
	}

	// トークンを生成（例として簡略化）
	token := "some_generated_token"
	return token, nil
}

// SignUpは新規ユーザー登録を実行
func (bu *BudgetUsecase) SignUp(user models.User) error {
	return bu.br.Create(user)
}

// パスワードリセットリンクを生成する関数
func (bu *BudgetUsecase) GenerateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (bu *BudgetUsecase) SaveResetToken(email, token string) error {
	// トークンをデータベースに保存
	resetToken := models.ResetToken{
		Token: token,
		Email: email,
	}
	return bu.br.RegisterToken(resetToken)
}

// パスワードリセットメールを実行
func (bu *BudgetUsecase) SendResetEmail(email, token string) error {
	// .envファイルを読み込む
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数からデータベース接続情報を取得
	server_host := os.Getenv("SERVER_HOST")
	server_port := os.Getenv("SERVER_PORT")
	gmail_credentials := os.Getenv("GMAIL_CREDENTIALS")
	gmail_token := os.Getenv("GMAIL_TOKEN")

	subject := "Password Reset"
	body := fmt.Sprintf("Click the link to reset your password: http://%s:%s/reset-password?token=%s", server_host, server_port, token)
	utils.Send_Gmail(email, subject, body, gmail_credentials, gmail_token)

	return nil

}

// ResetPasswordはパスワードリセットを実行
func (bu *BudgetUsecase) ResetPassword(token string, password string) error {
	// tokenからemailを取得
	email_db, err := bu.br.GetEmailByToken(token)
	if err != nil {
		return err
	}
	// ユーザーを取得
	user, err := bu.br.FindByEmail(email_db)
	if err != nil {
		return err
	}
	// パスワードリセットのロジック（例として簡略化）
	user.Password = password
	err = bu.br.Update(user)
	if err != nil {
		return err
	}
	// tokenを削除
	err = bu.br.DeleteToken(token)
	if err != nil {
		return err
	}
	return nil
}
