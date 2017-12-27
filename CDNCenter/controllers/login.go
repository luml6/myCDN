package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	sess := c.StartSession()
	username := sess.Get("userName")
	beego.Debug(username)
	if username == nil || username == "" {
		c.TplNames = "login.html"
	} else {
		c.Redirect("/download", 302)
	}
}
func (m *MainController) Post() {
	sess := m.StartSession()
	var user models.TAdminlist
	inputs := m.Input()
	user.Username = inputs.Get("userName")
	user.Pwd = inputs.Get("password")
	v, err := models.ValidateUser(user.Username)
	if err == nil {
		if v.Pwd == user.Pwd {
			sess.Set("userName", user.Username)
			beego.Debug("userName:", sess.Get("userName"))
			m.Redirect("/download", 302)
		} else {
			m.Data["IsSelect"] = true
			m.Data["ErrorMsg"] = "密码错误！"
			beego.Debug("密码错误！")
		}
	} else {
		beego.Debug(err)
		m.Data["IsSelect"] = true
		m.Data["ErrorMsg"] = "用户不存在"
		beego.Debug("用户不存在")
	}
	m.TplNames = "login.html"
}
