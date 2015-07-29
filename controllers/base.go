package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"yts/models"
//	"fmt"
	"strings"
)

type BaseController struct {
	beego.Controller
	Info *map[string]string
	ControllerName string
	ActionName string
	User *models.User
	Token string
}

func (this *BaseController) Prepare() {
	this.Layout = "layout/layout.html"
	this.ControllerName, this.ActionName = this.GetControllerAndAction()
	this.Auth()
}

func (this *BaseController) Auth() {
	sessionVal := this.GetSession("sessionVal")
	switch this.ControllerName {
	case "LoginController" :
		switch this.ActionName {
		case "Login", "Dologin":
			if sessionVal != nil {
				this.Redirect("/ts", 302)
			}
		}
	case "TSController" :
		switch this.ActionName {
		case "PHP", "JS", "Post":
			if sessionVal != nil {
				//this.Redirect("/ts", 302)
			}
		default:
			if sessionVal == nil {
				this.CustomRedirect()
			}
		}
	case "TSKLController" :
		switch this.ActionName {
		case "GetKL", "PostKL":
			if sessionVal != nil {
				this.Redirect("/ts", 302)
			}
		default:
			if sessionVal == nil {
				this.CustomRedirect()
			}
		}
	case "SQController" :
		switch this.ActionName {
		case "Show", "Create":
			if sessionVal != nil {
				this.Redirect("/ts", 302)
			}
		default:
			if sessionVal == nil {
				this.CustomRedirect()
			}
		}
	default:
		if sessionVal == nil {
			this.CustomRedirect()
		}
	}
	if sessionVal != nil {
		sessionMap := sessionVal.(map[string]string)
		this.User, _ = models.NewUserOption().Get(sessionMap["name"])
		this.Token = sessionMap["token"]
		if this.ControllerName == "UserController" && this.User.Name != "admin" {
			this.CustomRedirect()
		}
		if this.ControllerName == "SQController" && strings.HasPrefix(this.ActionName, "Admin") && this.User.Name != "admin" {
			this.CustomRedirect()
		}
	}
}

func (this *BaseController) CustomRedirect() {
	this.Redirect("/login", 302)
}

func (this *BaseController) OK(suc interface {}) {
	if index, yes := suc.(string); yes {
		this.Data["json"] = map[string]string{"success" : (*this.Info)[index]}
	} else {
		this.Data["json"] = suc
	}
	this.ServeJson()
}

func (this *BaseController) ERR(err string) {
	errStr := err
	if tmpErrStr, ok := (*this.Info)[err]; ok {
		errStr = tmpErrStr
	}
	this.Data["json"] = map[string]string{"error" : errStr}
	this.ServeJson()
}

func (this *BaseController) OutUser() map[string]string {
	return map[string]string{"Name" : this.User.Nick}
}