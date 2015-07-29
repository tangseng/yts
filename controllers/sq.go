package controllers

import (
	"yts/models"
	"fmt"
)

type SQController struct {
	BaseController
}

func (this *SQController) Show() {
	sqs := models.NewSQOption().GetAll()
	for _, v := range sqs {
		fmt.Printf("%v", v)
	}
	outSqs := make([]map[string]interface {}, len(sqs))
	for k, v := range sqs {
		outSqs[k] = map[string]interface {}{
			"name" : v.Nick,
			"status" : v.Status,
		}
	}
	this.Data["Sqs"] = outSqs
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/sq/sqCtrl"}
	this.TplNames = "sq.html"
}

func (this *SQController) Create() {
	this.Info = &map[string]string{
		"fail" : "信息不全",
		"sq" : "已有人申请",
		"no" : "有问题，失败了，请再来！",
		"ok" : "成功啦！",
	}
	name := this.GetString("loginName")
	password := this.GetString("loginPass")
	nick := this.GetString("name")
	if len(name) == 0 || len(password) == 0 || len(nick) == 0 {
		this.ERR("fail")
		return
	}
	sq := models.NewSQ(name, nick, password)
	sqOption := models.NewSQOption()
	ok := sqOption.Check(sq)
	if !ok {
		this.ERR("sq")
		return
	}
	err := sqOption.Set(sq)
	if err != nil {
		this.ERR("no")
		return
	}
	this.OK("ok")
}

func (this *SQController) AdminShow() {
	sqs := models.NewSQOption().GetAll()
	this.Data["Sqs"] = sqs
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/sq/sqadminCtrl"}
	this.TplNames = "sqadmin.html"
}

func (this *SQController) AdminStatus() {
	this.Info = &map[string]string{
		"user" : "已存在该申请用户！",
		"ok" : "OK！",
	}
	name := this.GetString("name")
	status, _ := this.GetInt8("status")
	sqOption := models.NewSQOption()
	sq, err := sqOption.Get(name)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	if status == 1 {
		userOption := models.NewUserOption()
		user, uerr := userOption.Get(name)
		if uerr == nil || user != nil {
			this.ERR("user")
			return
		}
		salt := models.RandStr2(5)
		user = models.NewUser(sq.Name, sq.Nick, sq.Password, salt, 0)
		err = userOption.Create(user)
		if err != nil {
			this.ERR(err.Error())
			return
		}
	}
	sq.Status = status
	err = sqOption.Set(sq)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	this.OK("ok")
}
