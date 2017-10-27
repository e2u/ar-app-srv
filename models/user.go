package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	UserInActive = 0
	UserActive   = 1
)

type User struct {
	Id        string    `json:"id,omitempty" gorm:"column:id"`
	UserName  string    `json:"username,omitempty" gorm:"column:username"`
	Password  string    `json:"-" gorm:"column:password"`
	NickName  string    `json:"nickname" gorm:"column:nickname"`
	SchoolId  int       `json:"school-id" gorm:"column:school_id"`
	Status    int       `json:"-" gorm:"status"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetUserByName(db *gorm.DB, userName string) (*User, error) {
	var user User
	if db.Limit(1).Model(&User{}).Where("username = ? ", userName).Find(&user).RecordNotFound() {
		return nil, sql.ErrNoRows
	}
	return &user, nil
}
