package routers

import (
	"github.com/astaxie/beego/context"
	"newsWeb/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/article/index",&controllers.ArticleController{},"get,post:ShowIndex")
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/article/delete",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
	beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}

func filterFunc(ctx *context.Context){

	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}


}