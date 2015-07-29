package controllers

import (
	"yts/models"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"github.com/astaxie/beego"
	"time"
)

type TSController struct {
	BaseController
}

const (
	PHPPATH = "./conf/ts.php"
	PHPBHPATH = "./conf/ts.bh.php"
	JSPATH = "./conf/ts.js"
)
func (this *TSController) PHP() {
	this.Info = &map[string]string{
		"fail" : "验证不通过",
		"user" : "用户已存在",
		"ok" : "成功啦！",
	}
	token := this.GetString("t")
	if len(token) == 0 {
		this.ERR("fail")
		return
	}
	path := PHPBHPATH
	if len(token) > 0 {
		_, err := models.NewTokenOption().Get(token)
		if err == nil {
			path = PHPPATH
		}
	}
	phptag := ""
	if this.GetString("include") == "1" {
		phptag = "<?php"
	}
	replaces := map[string]string{
		"token" : token,
		"domain" : this.domain(),
		"phptag" : phptag,
	}
	cc, _ := ioutil.ReadFile(path)
	for kk, vv := range replaces {
		re := regexp.MustCompile(fmt.Sprintf(`{%s}`, kk))
		reByte := []byte(fmt.Sprintf("%s", vv))
		cc = re.ReplaceAll(cc, reByte)
	}
	this.Ctx.Output.Body(cc)
}

func (this *TSController) JS() {
	this.Info = &map[string]string{
		"fail" : "账号或密码不对",
		"user" : "用户已存在",
		"ok" : "成功啦！",
	}
	token := this.GetString("t")
	if len(token) == 0 {
		this.ERR("fail")
		return
	}
	_, err := models.NewTokenOption().Get(token)
	if err != nil {
		this.ERR("fail")
		return
	}
	kl, _ := this.GetInt("kl")
	tskl := "0"
	if kl > 0 {
		tskl = "1"
	}
	replaces := map[string]string{
		"token" : token,
		"domain" : this.domain(),
		"tskl" : tskl,
	}
	cc, _ := ioutil.ReadFile(JSPATH)
	for kk, vv := range replaces {
		re := regexp.MustCompile(fmt.Sprintf(`{%s}`, kk))
		reByte := []byte(fmt.Sprintf("%s", vv))
		cc = re.ReplaceAll(cc, reByte)
	}
	this.Ctx.Output.Body(cc)
}

func (this *TSController) Get() {
	tss, _ := models.NewTSOption().Get(this.User.Name, 0, 1)
	var hash string
	if len(tss) > 0 {
		hash = models.MD5(tss[0].Data, "")
	}
	this.Data["User"] = this.OutUser()
	this.Data["Hash"] = hash
	this.Data["TS"] = tss
	this.Data["PHPUrl"] = fmt.Sprintf("http://%s/ts/php?t=%s", this.domain(), this.Token)
	this.Data["JSUrl"] = fmt.Sprintf("http://%s/ts/js?t=%s", this.domain(), this.Token)
	this.Data["Script"] = []string{"jquery.jsonview", "app/app", "app/directives/tip", "app/directives/jsonview", "app/services/time", "app/controllers/ts/tsCtrl"}
	this.Data["Css"] = []string{"jquery.jsonview"}
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.TplNames = "ts.html"
}

func (this *TSController) AjaxGet() {
	hash := this.GetString("hash")
	name := this.User.Name
	tsOption := models.NewTSOption()
	var tss []models.TS
	var nowHash string
	var i int8
	var max int8 = 30
	for {
		if i > max {
			break
		}
		tss, _ = tsOption.Get(name, 0, 1)
		if len(tss) > 0 {
			nowHash = models.MD5(tss[0].Data, "")
		}
		if hash != nowHash {
			break
		}
		time.Sleep(time.Second)
		i++
	}
	if i > max {
		this.OK(map[string]interface {}{
			"hash" : hash,
		})
		return
	}
	this.OK(map[string]interface {}{
		"ts" : tss,
		"hash" : nowHash,
	})
}

func (this *TSController) More() {
	tss, _ := models.NewTSOption().Get(this.User.Name, 0, 10)
	this.Data["TS"] = tss
	this.Data["Script"] = []string{"jquery.jsonview", "app/app", "app/directives/tip", "app/directives/jsonview", "app/services/time", "app/controllers/tsmore/tsmoreCtrl"}
	this.Data["Css"] = []string{"jquery.jsonview"}
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.TplNames = "tsmore.html"
}


func (this *TSController) Post() {
	token := this.GetString("token")
	time := this.GetString("time")
	data := this.GetString("data")
	te := this.GetString("type")
	fmt.Println(data)
	ip := this.Ctx.Input.IP()
	go func(token, time, data, te, ip string) {
		t, err := models.NewTokenOption().Get(token)
		if err != nil {
			return
		}
		tstime, _ := strconv.ParseInt(time, 10, 64)
		ts := models.NewTS(data, te, ip, tstime, 0)
		models.NewTSOption().Insert(t.User.Name, ts)
	}(token, time, data, te, ip)
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Ctx.WriteString("")
}

func (this *TSController) domain() string {
	return fmt.Sprintf("%s:%s", beego.AppConfig.String("Domain"), beego.AppConfig.String("httpport"))
}