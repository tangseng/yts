package controllers

import (
	"yts/models"
	"fmt"
)

type UserController struct {
	BaseController
}

func (this *UserController) Show() {
	users := models.NewUserOption().GetAll()
	for _, v := range users {
		fmt.Printf("%v", v)
	}
	outUsers := make([]map[string]string, len(users))
	for k, v := range users {
		outUsers[k] = map[string]string{
			"name" : v.Nick,
			"loginName" : v.Name,
			"loginPass" : v.Password,
		}
	}
	this.Data["Persons"] = outUsers
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/person/personCtrl"}
	this.TplNames = "user.html"
}

func (this *UserController) Create() {
	this.Info = &map[string]string{
		"fail" : "信息不全",
		"user" : "用户已存在",
		"ok" : "成功啦！",
	}
	name := this.GetString("loginName")
	password := this.GetString("loginPass")
	nick := this.GetString("name")
	if len(name) == 0 || len(password) == 0 || len(nick) == 0 {
		this.ERR("fail")
		return
	}
	userOption := models.NewUserOption()
	user, err := userOption.Get(name)
	if user != nil {
		this.ERR("user")
		return
	}
	salt := models.RandStr2(5)
	user = models.NewUser(name, nick, password, salt, 0)
	err = userOption.Create(user)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	this.OK("ok")
}

func (this *UserController) Update() {
	this.Info = &map[string]string{
		"fail" : "信息不全",
		"user" : "用户不存在",
		"ok" : "更新成功！",
	}
	name := this.GetString("loginName")
	password := this.GetString("loginPass")
	nick := this.GetString("name")
	if len(name) == 0 || len(password) == 0 || len(nick) == 0 {
		this.ERR("fail")
		return
	}
	userOption := models.NewUserOption()
	user, err := userOption.Get(name)
	if user == nil || err != nil {
		this.ERR("user")
		return
	}
	user.Password = password
	user.Nick = nick
	err = userOption.Update(user)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	this.OK("ok")
}

func (this *UserController) Delete() {
	this.Info = &map[string]string{
		"ok" : "删除成功！",
	}
	name := this.GetString("loginName")
	userOption := models.NewUserOption()
	err := userOption.Delete(name)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	this.OK("ok")
}
