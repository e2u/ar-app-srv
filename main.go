// Deposit 储值服务

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"e2u.io/ar-app-srv/app"
	"e2u.io/ar-app-srv/controllers"
	"e2u.io/ar-app-srv/middle"
	"github.com/e2u/goboot"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	ListenPort int
	RunEnv     string
	App        *app.AppContext
	stoping    bool
	accessLog  *os.File
)

func initApp() error {
	App = app.NewAppContext()
	App.RunMode = RunEnv
	App.AppName = goboot.Config.MustString("app.name", "ar-app-srv")

	App.InitTables()

	if App.RunMode == "dev" {
		App.InitData()
	}

	return nil
}

// 初始化所有参数
func initParamter() error {
	return nil
}

// 初始化函数
func init() {
	flag.StringVar(&RunEnv, "env", "dev", "app run env: [dev|dev-prod|prod]")
	flag.IntVar(&ListenPort, "port", 9000, "http listen port: [9000|9001]")
	flag.Parse()

	goboot.Init(RunEnv)
	goboot.OnAppStart(initApp, 10)
	goboot.OnAppStart(initParamter, 20)
	goboot.Startup()

	if al := goboot.Config.MustString("log.access", "stdout"); al == "stdout" {
		accessLog = os.Stdout
	} else {

		f, err := os.OpenFile(al, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			accessLog = os.Stdout
		} else {
			accessLog = f
		}
	}
}

func main() {
	goboot.Log.Infof("run mode: %v\n", RunEnv)

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(accessLog, "", log.LstdFlags)}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.DefaultCompress)

	r.Route("/v1", func(sr chi.Router) {

		c := &controllers.Controller{
			AppContext: App,
		}
		sr.Use(middle.InitController(c))

		// 取所有学校列表
		sr.Get("/schools", c.GetAllSchools)

		// 用户注册
		sr.Post("/user/reg", c.UserReg)
		// 用户登录
		sr.Post("/user/sign-in", c.UserSignIn)
		// 用户退出
		sr.Post("/user/sign-out", c.UserSignOut)

		// 获取进度
		sr.Get("/schedules", c.GetUserSchedules)

		// 更新进度
		sr.Post("/schedules/update", c.UpdateUserSchedules)

	})

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", ListenPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-stopChan // wait for SIGINT
	App.DB.Close()
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")

}
