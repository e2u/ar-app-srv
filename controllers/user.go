package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"e2u.io/ar-app-srv/models"
	"e2u.io/ar-app-srv/util"
	"github.com/e2u/goboot"
)

// 用户注册 UserReg
// curl http://127.0.0.1:9000/v1/user/reg -d 'username=13999999999&password=123456&nickname=姓名&school-id=1'
func (c *Controller) UserReg(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)

	var userName string
	var password string
	var nickName string
	var schoolId int

	c.params.Bind(&userName, "username")
	c.params.Bind(&password, "password")
	c.params.Bind(&nickName, "nickname")
	c.params.Bind(&schoolId, "school-id")

	if len(userName) == 0 || len(password) == 0 || len(nickName) == 0 {
		c.Error("必选参数不得为空")
		return
	}

	sc, err := models.NewSchool().FindById(c.AppContext.DB, schoolId)
	if err != nil || sc == nil {
		c.Error("非法学校id")
		return
	}

	cu, err := models.NewUser().GetUserByName(c.AppContext.DB, userName)
	if err != nil && err != sql.ErrNoRows {
		c.Error(err.Error())
		return
	}

	if cu != nil {
		c.Error("用户已经存在.")
		return
	}

	newUser := &models.User{
		Id:        util.MD5String([]byte(userName)),
		UserName:  userName,
		NickName:  nickName,
		Password:  util.MD5String([]byte(password)),
		SchoolId:  schoolId,
		Status:    models.UserActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.Save(c.AppContext.DB, newUser); err != nil {
		c.Error(err.Error())
		return
	}

	c.Success()
	c.PutResp(RemarkTag, "注册成功")

}

// 用户登录 UserSignIn,登录成功返回 access_token
// curl  http://127.0.0.1:9000/v1/user/sign-in -d "username=13999999999&password=123456"
func (c *Controller) UserSignIn(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)
	var userName string
	var password string

	c.params.Bind(&userName, "username")
	c.params.Bind(&password, "password")

	if len(userName) == 0 || len(password) == 0 {
		c.Error("必选参数不得为空")
		return
	}

	errMsg := "登录失败,请检查用户名密码."

	cu, err := models.NewUser().GetUserByName(c.AppContext.DB, userName)
	if err != nil && err == sql.ErrNoRows || cu == nil {
		goboot.Log.Errorf("登录失败 -- err: %v errMsg: %v", err.Error(), errMsg)
		c.Error(errMsg)
		return
	}

	inputPassword := util.MD5String([]byte(password))
	if cu.Password != inputPassword {
		goboot.Log.Warningf("登录失败 -- 密码错误 %v", cu.UserName)
		c.Error(errMsg)
		return
	}

	// 生成 accessToken 並返回
	token := util.MD5String([]byte(time.Now().Format(time.RFC3339Nano) + "-" + cu.Id))
	newToken := &models.AccessToken{
		UserId:    cu.Id,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.NewAccessToken().CreateOrUpdate(c.AppContext.DB, newToken); err != nil {
		goboot.Log.Warningf("登录失败 -- 生成 Token 发生错误 %v", err.Error())
		c.Error("登录失败,系统错误")
		return
	}

	c.Success()
	c.PutResp("userid", cu.Id)
	c.PutResp("username", cu.UserName)
	c.PutResp("nickname", cu.NickName)
	c.PutResp("access_token", token)

}

// UserSignOut 用户退出登录
// curl http://127.0.0.1:9000/v1/user/sign-out -d "userid=502193ee233715b4c40e172206d4dd45&access-token="
func (c *Controller) UserSignOut(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)

	var userId string
	var accessToken string

	c.params.Bind(&userId, "userid")
	c.params.Bind(&accessToken, "access-token")

	if len(userId) == 0 || len(accessToken) == 0 {
		c.Error("必选参数不得为空")
		return
	}

	if err := models.NewAccessToken().DeleteByUserIdAndAccessToken(c.AppContext.DB, userId, accessToken); err != nil {
		c.Error(err.Error())
		c.PutResp(RemarkTag, "退出登录失败")
		return
	}
	c.Success()
	c.PutResp(RemarkTag, "退出成功")
}
