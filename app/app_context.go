package app

import (
	"fmt"
	"time"

	"e2u.io/ar-app-srv/models"
	"e2u.io/ar-app-srv/util"
	"github.com/e2u/goboot"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// 初始化数据库连接等
type AppContext struct {
	AppName string
	RunMode string
	DevMode bool
	DB      *gorm.DB
}

func NewAppContext() *AppContext {

	app := &AppContext{
		AppName: goboot.Config.MustString("app.name", "huafei-go-noname"),
		DevMode: goboot.Config.MustBool("mode.dev", true),
	}

	if err := app.initSQLite3(); err != nil {
		goboot.Log.Errorf("初始化数据库连接错误: %v", err.Error())
		panic(fmt.Sprintf("初始化数据库连接错误: %v", err.Error()))
	}

	return app
}

func (a *AppContext) initSQLite3() error {
	var err error
	a.DB, err = gorm.Open("sqlite3", "default.db")
	if err != nil {
		return err
	}

	a.DB.LogMode(a.DevMode)
	return a.DB.DB().Ping()
}

func (a *AppContext) InitTables() error {
	goboot.Log.Info("init tables")
	create := "create table if not exists"
	// 学校表
	schoolTable := `schools (id INTEGER NOT NULL PRIMARY KEY, name TEXT,created_at DATETIME,updated_at DATETIME)`
	// 用户表
	userTable := `users (id TEXT NOT NULL PRIMARY KEY,username TEXT,password TEXT,nickname TEXT,school_id INTEGER,status INTEGER,created_at DATETIME,updated_at DATETIME)`
	// 任务定义表
	taskTable := `tasks (id TEXT NOT NULL PRIMARY KEY,name TEXT,descript TEXT,created_at DATETIME,updated_at DATETIME)`
	// 用户任务关联表
	userTasksTable := `user_tasks (id INTEGER NOT NULL PRIMARY KEY,user_id TEXT,task_id INTEGER,status TEXT,created_at DATETIME,updated_at DATETIME)`
	// 登录状态记录
	accessTokenTable := `access_tokens (access_token NOT NULL PRIMARY KEY,user_id TEXT,created_at DATETIME,updated_at DATETIME)`
	ts := []string{schoolTable, userTable, taskTable, userTasksTable, accessTokenTable}
	for idx := range ts {
		a.DB.Exec(fmt.Sprintf("%s %s", create, ts[idx]))
	}
	return nil
}

// 初始化数据
func (a *AppContext) InitData() error {

	goboot.Log.Info("init data")

	userName := "13999999999"
	nickName := "昵称"
	password := "123456"
	schoolId := 1

	tasks := make([]interface{}, 0)
	for i := 1; i <= 10; i++ {
		tasks = append(tasks, &models.Task{
			Id:       util.MD5String([]byte(fmt.Sprintf("%d", i))),
			Name:     fmt.Sprintf("任务-%d", i),
			Descript: fmt.Sprintf("任务描述-%d", i),
		})
	}

	ds := []interface{}{
		&models.School{Id: 1, Name: "学校01", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		&models.User{
			Id:        util.MD5String([]byte(userName)),
			UserName:  userName,
			NickName:  nickName,
			Password:  util.MD5String([]byte(password)),
			SchoolId:  schoolId,
			Status:    models.UserActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ds = append(ds, tasks...)

	for idx := range ds {
		models.Save(a.DB, ds[idx])
	}

	return nil
}
