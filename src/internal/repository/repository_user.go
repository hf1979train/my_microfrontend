package repository

import (
	"myhomesv/internal/domain/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// UserRepositoryはユーザーに関するリポジトリのインターフェース
type IUserRepository interface {
	FindByUsername(username string) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Create(user models.User) error
	Update(user models.User) error
	RegisterToken(resetToken models.ResetToken) error
	GetEmailByToken(token string) (string, error)
	DeleteToken(token string) error
}

// FindByUsernameはユーザー名でユーザーを検索
func (r *BudgetRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

// FindByEmailはメールアドレスでユーザーを検索
func (r *BudgetRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Createは新しいユーザーを作成
func (r *BudgetRepository) Create(user models.User) error {
	return r.db.Table("users").Create(&user).Error
}

// Updateは既存のユーザー情報を更新
func (r *BudgetRepository) Update(user models.User) error {
	return r.db.Table("users").Save(&user).Error
}

// email, tokenのペアを登録する関数
func (r *BudgetRepository) RegisterToken(resetToken models.ResetToken) error {
	// トークンをハッシュ化
	hashedToken, err := hashToken(resetToken.Token)
	if err != nil {
		return err
	}
	resetToken.Token = hashedToken
	resetToken.CreatedAt = time.Now()
	// アップサート(updateとinsertの組み合わせ)のクエリを実行
	return r.db.Table("reset_tokens").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "created_at"}),
	}).Create(&resetToken).Error

	// var existingToken models.ResetToken
	// if err := r.db.Table("reset_tokens").Where("email = ?", resetToken.Email).First(&existingToken).Error; err == nil {
	// 	// 既存のレコードが見つかった場合、更新する
	// 	existingToken.Token = resetToken.Token
	// 	existingToken.CreatedAt = time.Now()
	// 	return r.db.Save(&existingToken).Error
	// } else if gorm.ErrRecordNotFound == err {
	// 	// 既存のレコードが見つからなかった場合、新しいレコードを挿入する
	// 	return r.db.Table("reset_tokens").Create(&resetToken).Error
	// } else {
	// 	// その他のエラーが発生した場合
	// 	return err
	// }
}

//	return r.db.Table("reset_tokens").Save(&resetToken).Error
//}

// tokenからemailを出力する関数
func (r *BudgetRepository) GetEmailByToken(token string) (string, error) {
	hashedToken, err := hashToken(token)
	if err != nil {
		return "", err
	}
	var resetToken models.ResetToken
	err = r.db.Table("reset_tokens").Where("token = ?", hashedToken).First(&resetToken).Error
	if err != nil {
		return "", err
	}
	return resetToken.Email, nil
}

// tokenを指定して削除する関数
func (r *BudgetRepository) DeleteToken(token string) error {
	hashedToken, err := hashToken(token)
	if err != nil {
		return err
	}
	return r.db.Table("reset_tokens").Where("token = ?", hashedToken).Delete(&models.ResetToken{}).Error
}

// tokenをハッシュ化する関数
func hashToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}
