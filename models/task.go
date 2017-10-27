package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

// Task 任务定义
type Task struct {
	Id        string    `json:"id,omitempty" gorm:"column:id"`
	Name      string    `json:"name,omitempty" gorm:"column:name"`
	Descript  string    `json:"descript,omitempty" gorm:"column:descript"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func NewTask() *Task {
	return &Task{}
}

// 取所有任务
func (t *Task) FindAll(db *gorm.DB) ([]*Task, error) {
	tl := make([]*Task, 0)
	if err := db.Model(&Task{}).Find(&tl).Error; err != nil {
		return nil, err
	}
	return tl, nil
}

// FindTaskByTaskId
func (ut *Task) FindTaskByTaskId(db *gorm.DB, taskId string) (*Task, error) {
	var t Task
	if db.Limit(1).Model(&Task{}).Where("task_id = ?", taskId).Find(&t).RecordNotFound() {
		return nil, sql.ErrNoRows
	}
	return ut, nil
}
