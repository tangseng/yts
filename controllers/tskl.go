package controllers

import (
	"yts/models"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

type TSKLController struct {
	BaseController
}

func (this *TSKLController) Get() {
	tskl, err := models.NewTSKLOption().Get(this.User)
	tsklMap := make(map[string]string)
	if err == nil {
		tsklMap["send"] = tskl.Send
		tsklMap["back"] = tskl.Back
		tsklMap["berr"] = tskl.Berr
		tsklMap["hash"] = tskl.Hash
	}
	this.Data["User"] = this.OutUser()
	this.Data["TSKL"] = tsklMap
	this.Data["STEP"] = beego.AppConfig.String("Step")
	this.Data["JSUrl"] = fmt.Sprintf("http://%s/ts/js?t=%s&kl=1", this.domain(), this.Token)
	this.Data["Script"] = []string{"jquery.jsonview", "app/app", "app/directives/tip", "app/directives/jsonview", "app/services/time", "app/controllers/tskl/tsklCtrl"}
	this.Data["Css"] = []string{"jquery.jsonview"}
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Navbar"] = "layout/navbar.html"
	this.TplNames = "tskl.html"
}

func (this *TSKLController) AjaxGet() {
	hash := this.GetString("hash")
	tsklOption := models.NewTSKLOption()
	var tskl *models.TSKL
	var i int8
	var max int8 = 30
	var err error
	for {
		if i > max {
			break
		}
		tskl, err = tsklOption.Get(this.User)
		if err == nil && (len(tskl.Back) > 0 || len(tskl.Berr) > 0) {
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
	output := map[string]interface {}{
		"back" : tskl.Back,
		"berr" : tskl.Berr,
	}
	if hash != tskl.Hash {
		output["send"] = tskl.Send
		output["hash"] = tskl.Hash
	}
	this.OK(output)
}


func (this *TSKLController) Post() {
	this.Info = &map[string]string{
		"fail" : "Oh, No!!!",
		"ok" : "OK！请等待结果的打印！",
	}
	send := this.GetString("send")
	if len(send) == 0 {
		this.ERR("fail")
		return
	}
	tskl := models.NewTSKL(this.User)
	tskl.Send = send
	err := models.NewTSKLOption().Set(tskl)
	if err != nil {
		this.ERR(err.Error())
		return
	}
	this.OK(map[string]string{
		"hash" : tskl.Hash,
		"ok" : (*this.Info)["ok"],
	})
}

func (this *TSKLController) GetKL() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Info = &map[string]string{
		"fail" : "Oh, No!!!",
	}
	token := this.GetString("t")
	if len(token) == 0 {
		this.ERR("fail")
		return
	}
	t, err := models.NewTokenOption().Get(token)
	if err != nil {
		this.ERR("fail")
		return
	}

	hash := this.GetString("hash")
	var tskl *models.TSKL
	var klerr error
	var i = 0
	var max = 30
	for {
		if i > max {
			break
		}
		tskl, klerr = models.NewTSKLOption().Get(t.User)
		if klerr == nil && hash != tskl.Hash && len(tskl.Berr) == 0 && len(tskl.Back) == 0 {
			break
		}
		time.Sleep(time.Second)
		i++
	}
	tsklMap := make(map[string]string)
	if klerr == nil && tskl.Hash != hash {
		tsklMap["send"] = tskl.Send
		tsklMap["hash"] = tskl.Hash
	} else {
		tsklMap["hash"] = hash
	}
	this.OK(tsklMap)
}

func (this *TSKLController) PostKL() {
	token := this.GetString("t")
	hash := this.GetString("hash")
	back := this.GetString("back")
	berr := this.GetString("berr")
	go func(token, hash, back, berr string) {
		t, err := models.NewTokenOption().Get(token)
		if err != nil {
			return
		}
		tsklOption := models.NewTSKLOption()
		tskl, err2 := tsklOption.Get(t.User)
		if err2 != nil {
			return
		}
		if tskl.Hash != hash {
			return
		}
		fmt.Println(tskl)
		tskl.Back = back
		tskl.Berr = berr
		fmt.Println(tskl)
		tsklOption.SetWhich(tskl, "back")
		tsklOption.SetWhich(tskl, "berr")
	}(token, hash, back, berr)
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	this.Ctx.WriteString("")
}

func (this *TSKLController) domain() string {
	return fmt.Sprintf("%s:%s", beego.AppConfig.String("Domain"), beego.AppConfig.String("httpport"))
}

