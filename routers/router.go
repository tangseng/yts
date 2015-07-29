package routers

import (
	"yts/controllers"
	"github.com/astaxie/beego"
)

func init() {
	loginC := &controllers.LoginController{}
	login := beego.NewNamespace("/login",
		beego.NSRouter("/", loginC, "get:Login"),
		beego.NSRouter("/in", loginC, "post:DoLogin"),
		beego.NSRouter("/out", loginC, "get,post:LoginOut"),
	)
	tsC := &controllers.TSController{}
	ts := beego.NewNamespace("/ts",
		beego.NSRouter("/php", tsC, "get:PHP"),
		beego.NSRouter("/js", tsC, "get:JS"),
		beego.NSRouter("/", tsC, "get:Get"),
		beego.NSRouter("/more", tsC, "get:More"),
		beego.NSRouter("/ajax", tsC, "get:AjaxGet"),
		beego.NSRouter("/post", tsC, "get,post:Post"),
	)
	tsklC := &controllers.TSKLController{}
	tskl := beego.NewNamespace("/tskl",
		beego.NSRouter("/", tsklC, "get:Get"),
		beego.NSRouter("/ajax", tsklC, "get:AjaxGet"),
		beego.NSRouter("/post", tsklC, "get,post:Post"),
		beego.NSRouter("/getkl", tsklC, "get,post:GetKL"),
		beego.NSRouter("/postkl", tsklC, "get,post:PostKL"),
	)
	userC := &controllers.UserController{}
	user := beego.NewNamespace("/user",
		beego.NSRouter("/", userC, "get:Show"),
		beego.NSRouter("/create", userC, "post:Create"),
		beego.NSRouter("/update", userC, "post:Update"),
		beego.NSRouter("/delete", userC, "post:Delete"),
	)
	sqC := &controllers.SQController{}
	sq := beego.NewNamespace("/sq",
		beego.NSRouter("/", sqC, "get:Show"),
		beego.NSRouter("/post", sqC, "post:Create"),
		beego.NSRouter("/admin", sqC, "get:AdminShow"),
		beego.NSRouter("/status", sqC, "post:AdminStatus"),
	)
	beego.AddNamespace(login, user, ts, tskl, sq)
}
