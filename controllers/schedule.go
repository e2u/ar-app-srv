package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"e2u.io/ar-app-srv/models"
)

// 取登录用户的进度
// curl "http://127.0.0.1:9000/v1/schedules?userid=502193ee233715b4c40e172206d4dd45&access-token=51645895aba921347f173309b4938970"
func (c *Controller) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)

	var userId string
	var accessToken string

	c.params.Bind(&userId, "userid")
	c.params.Bind(&accessToken, "access-token")

	if len(userId) == 0 || len(accessToken) == 0 {
		c.Error("必选参数不得为空")
		return
	}

	// 检查用户是否登录
	_, err := models.NewAccessToken().FindByUserIdAndAccessToken(c.AppContext.DB, userId, accessToken)
	if err != nil {
		c.PutResp(RespCodeTag, ResponseError)
		c.PutResp(RemarkTag, "请先登录")
		return
	}

	// 读取任务定义列表 schedule,然后再读用户任务列表,然后遍历定义表数据,再从用户已经变更状态的记录中合并数据
	taskDefs, err := models.NewTask().FindAll(c.AppContext.DB)
	if err != nil {
		c.Error(err.Error())
		return
	}
	userTasks, err := models.NewUserTask().FindUserTasksByUserId(c.AppContext.DB, userId)
	if err != nil {
		c.Error(err.Error())
		return
	}

	// 把 userTasks 转换成 map,便于查找操作
	userTasksMap := make(map[string]*models.UserTask)
	for idx := range userTasks {
		ut := userTasks[idx]
		userTasksMap[ut.TaskId] = ut
	}

	type outTask struct {
		Id        string    `json:"id,omitempty"` //  任务 id
		Name      string    `json:"name,omitempty"`
		Descript  string    `json:"descript,omitempty"`
		Status    string    `json:"status,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	outTasks := make([]*outTask, 0)
	for idx := range taskDefs {
		task := taskDefs[idx]
		outTasks = append(outTasks, func(_t *models.Task) *outTask {
			// 如用户已经有这个任务,就从用户记录中取数据
			if v, ok := userTasksMap[_t.Id]; ok {
				return &outTask{
					Id:        _t.Id,
					Name:      _t.Name,
					Descript:  _t.Descript,
					Status:    v.Status,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				}
			}
			return &outTask{
				Id:       _t.Id,
				Name:     _t.Name,
				Descript: _t.Descript,
				Status:   models.UserTaskStatusNotUnLock,
			}
		}(task))
	}
	c.Success()
	c.PutResp("tasks", outTasks)
}

// UpdateUserSchedules 更新用户的任务状态
// curl http://127.0.0.1:9000/v1/schedules/update -d "task-id=c9f0f895fb98ab9159f51fd0297e236d&userid=502193ee233715b4c40e172206d4dd45&access-token=c9d11eddeb0d167884ae8510f9b0744f&task-status=in-process"
func (c *Controller) UpdateUserSchedules(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)

	var userId string
	var accessToken string
	var taskId string
	var taskStatus string

	c.params.Bind(&userId, "userid")
	c.params.Bind(&accessToken, "access-token")
	c.params.Bind(&taskId, "task-id")
	c.params.Bind(&taskStatus, "task-status")

	if len(userId) == 0 || len(accessToken) == 0 || len(taskId) == 0 || len(taskStatus) == 0 {
		c.Error("必选参数不得为空")
		return
	}

	// 检查用户是否登录
	_, err := models.NewAccessToken().FindByUserIdAndAccessToken(c.AppContext.DB, userId, accessToken)
	if err != nil {
		c.PutResp(RespCodeTag, ResponseError)
		c.PutResp(RemarkTag, "请先登录")
		return
	}

	_, err = models.NewTask().FindTaskByTaskId(c.AppContext.DB, taskId)
	if err != nil && err == sql.ErrNoRows {
		c.Error(err.Error())
		c.PutResp(RemarkTag, "任务不存在")
		return
	}

	userTask := &models.UserTask{
		Status: taskStatus,
		UserId: userId,
		TaskId: taskId,
	}

	if err := models.NewUserTask().CreateOrUpdate(c.AppContext.DB, userTask); err != nil {
		c.Error(err.Error())
		c.PutResp(RemarkTag, "更新失败")
		return
	}

	c.Success()
	c.PutResp(RemarkTag, "更新成功")
	return

}
