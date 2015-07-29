package main

import (
	_ "yts/routers"
	"github.com/astaxie/beego"
	"net/http"
	"html/template"
)

func main() {
	beego.Errorhandler("404", P404)
	beego.Run()
}

func P404(rw http.ResponseWriter, r *http.Request){
	t,_:= template.New("404.html").ParseFiles(beego.ViewsPath + "/other/404.html")
	t.Execute(rw, nil)
}