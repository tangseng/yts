package controllers

import (
	"yts/models"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) Login() {
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/login/loginCtrl"}
	this.TplNames = "login.html"
}

func (this *LoginController) DoLogin() {
	this.Info = &map[string]string{
		"success" : "登录成功",
		"login" : "登录失败",
	}
	name := this.GetString("name")
	password := this.GetString("pass")
	user, err := models.NewUserOption().Get(name)
	if err != nil || user.CheckPass(password) {
		this.ERR("login")
		return
	}
	token := models.NewToken(user, "")
	models.NewTokenOption().Add(token)
	this.SetSession("sessionVal", map[string]string{
		"name" : user.Name,
		"token" : token.Token(),
	})
	this.OK("success")
}

func (this *LoginController) LoginOut() {
	this.DelSession("sessionVal")
	models.NewTokenOption().Remove(this.Token)
	this.CustomRedirect()
}
