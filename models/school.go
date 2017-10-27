package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

type School struct {
	Id        int       `json:"id,omitempty" gorm:"column:id"`
	Name      string    `json:"name,omitempty" gorm:"column:name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func NewSchool() *School {
	return &School{}
}

func (s *School) FindAll(db *gorm.DB) ([]*School, error) {
	sl := make([]*School, 0)
	if err := db.Model(&School{}).Find(&sl).Error; err != nil {
		return nil, err
	}
	return sl, nil
}

func (s *School) FindById(db *gorm.DB, id int) (*School, error) {
	var sc School
	if db.Limit(1).Model(&School{}).Where("id = ? ", id).Find(&sc).RecordNotFound() {
		return nil, sql.ErrNoRows
	}
	return &sc, nil
}
