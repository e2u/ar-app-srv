package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

type AccessToken struct {
	Token     string    `json:"access_token,omitempty" gorm:"column:access_token"`
	UserId    string    `json:"user_id,omitempty" gorm:"column:user_id"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func NewAccessToken() *AccessToken {
	return &AccessToken{}
}

// CreateOrUpdate 创建或更新 token
func (a *AccessToken) CreateOrUpdate(db *gorm.DB, accessToken *AccessToken) error {

	var at AccessToken

	// create new access token
	if db.Limit(1).Model(&AccessToken{}).Where("user_id = ? ", accessToken.UserId).Find(&at).RecordNotFound() {
		return Save(db, accessToken)
	}

	// update access token
	return db.Model(&AccessToken{}).Update(accessToken).Error
}

// DeleteByUserIdAndAccessToken  根据 userId 和 accessToken 删除登录记录
func (a *AccessToken) DeleteByUserIdAndAccessToken(db *gorm.DB, userId, accessToken string) error {
	var at AccessToken
	if db.Limit(1).Model(&AccessToken{}).Where("user_id = ? and access_token = ? ", userId, accessToken).Find(&at).RecordNotFound() {
		return sql.ErrNoRows
	}
	return db.Delete(AccessToken{}, "user_id = ? and access_token = ? ", userId, accessToken).Error
}

// FindByUserIdAndAccessToken
func (a *AccessToken) FindByUserIdAndAccessToken(db *gorm.DB, userId, accessToken string) (*AccessToken, error) {
	var at AccessToken
	if db.Limit(1).Model(&AccessToken{}).Where("user_id = ? and access_token = ? ", userId, accessToken).Find(&at).RecordNotFound() {
		return nil, sql.ErrNoRows
	}
	return &at, nil
}
