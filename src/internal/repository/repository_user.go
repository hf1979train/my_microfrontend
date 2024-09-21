package repository

import (
	"myhomesv/internal/domain/models"
	"time"

	"gorm.io/gorm"
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
	var existingToken models.ResetToken
	if err := r.db.Table("reset_tokens").Where("email = ?", resetToken.Email).First(&existingToken).Error; err == nil {
		// 既存のレコードが見つかった場合、更新する
		existingToken.Token = resetToken.Token
		existingToken.CreatedAt = time.Now()
		return r.db.Save(&existingToken).Error
	} else if gorm.ErrRecordNotFound == err {
		// 既存のレコードが見つからなかった場合、新しいレコードを挿入する
		return r.db.Table("reset_tokens").Create(&resetToken).Error
	} else {
		// その他のエラーが発生した場合
		return err
	}
}

//	return r.db.Table("reset_tokens").Save(&resetToken).Error
//}

// tokenからemailを出力する関数
func (r *BudgetRepository) GetEmailByToken(token string) (string, error) {
	var resetToken models.ResetToken
	err := r.db.Table("reset_tokens").Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		return "", err
	}
	return resetToken.Email, nil
}

// tokenを指定して削除する関数
func (r *BudgetRepository) DeleteToken(token string) error {
	return r.db.Table("reset_tokens").Where("token = ?", token).Delete(&models.ResetToken{}).Error
}
