package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	UserTaskStatusNotUnLock = "not-unlock"
	UserTaskStatusInProcess = "in-process"
	UserTaskStatusDone      = "done"
)

var (
	UserTaskStatusSlice = []string{UserTaskStatusNotUnLock, UserTaskStatusInProcess, UserTaskStatusDone}
)

// UserTask 任务定义
type UserTask struct {
	Id        int       `json:"id,omitempty" gorm:"column:id"`
	UserId    string    `json:"user-id,omitempty" gorm:"column:user_id"`
	TaskId    string    `json:"task-id,omitempty" gorm:"column:task_id"`
	Status    string    `json:"status,omitempty" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func NewUserTask() *UserTask {
	return &UserTask{}
}

func (ut *UserTask) CheckStatus(t *UserTask) error {
	fmt.Println(ut)
	for idx := range UserTaskStatusSlice {
		if UserTaskStatusSlice[idx] == t.Status {
			return nil
		}
	}
	return errors.New("非法状态")
}

// FindUserTasksByUserId
func (ut *UserTask) FindUserTasksByUserId(db *gorm.DB, userId string) ([]*UserTask, error) {
	tl := make([]*UserTask, 0)
	if err := db.Model(&UserTask{}).Where("user_id = ?", userId).Find(&tl).Error; err != nil {
		return nil, err
	}
	return tl, nil
}

// FindUserTasksByUserIdAndTaskId
func (ut *UserTask) FindUserTasksByUserIdAndTaskId(db *gorm.DB, userId, taskId string) (*UserTask, error) {
	var t UserTask
	if db.Limit(1).Model(&UserTask{}).Where("user_id = ?  and task_id = ?", userId, taskId).Find(&t).RecordNotFound() {
		return nil, sql.ErrNoRows
	}
	return &t, nil
}

// CreateOrUpdate 创建或更新用户的任务状态
func (ut *UserTask) CreateOrUpdate(db *gorm.DB, userTask *UserTask) error {
	if err := ut.CheckStatus(userTask); err != nil {
		return err
	}
	t, err := ut.FindUserTasksByUserIdAndTaskId(db, userTask.UserId, userTask.TaskId)
	if err == sql.ErrNoRows {
		userTask.CreatedAt = time.Now()
		userTask.UpdatedAt = time.Now()
		return Save(db, userTask)
	} else if err != nil {
		return err
	} else if err == nil && t != nil && t.Status == UserTaskStatusDone {
		return errors.New("已完成的任务不能再变更状态")
	}

	userTask.UpdatedAt = time.Now()
	return db.Model(&UserTask{}).Update(userTask).Error
}
