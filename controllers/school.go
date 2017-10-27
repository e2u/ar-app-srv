package controllers

import (
	"net/http"

	"e2u.io/ar-app-srv/models"
)


// GetAllSchools  取所有学校列表
// curl http://127.0.0.1:9000/v1/schools
func (c *Controller) GetAllSchools(w http.ResponseWriter, r *http.Request) {
	defer c.RenderJSON(w, r)
	c.Success()
	allSchools, err := models.NewSchool().FindAll(c.AppContext.DB)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.PutResp("schools", allSchools)
}
