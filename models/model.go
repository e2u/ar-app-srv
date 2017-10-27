package models

import "github.com/jinzhu/gorm"

func Save(db *gorm.DB, v interface{}) error {
	return db.Save(v).Error
}
