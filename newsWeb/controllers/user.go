package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) HandleRegister() {
	name := this.GetString("userName")
	pwd := this.GetString("password")

	if name == "" || pwd == "" {
		fmt.Println("缺少数据内容")
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	var user models.User

	user.Name = name
	err := o.Read(&user, "Name")
	if err == nil { //没错误
		fmt.Println("用户名已存在，请重新输入")
		this.TplName = "register.html"
		return
	}

	user.Name = name
	user.Pwd = pwd
	_, err = o.Insert(&user)
	if err != nil { //有错误
		fmt.Println("注册失败")
		this.TplName = "register.html"
		return
	}
	this.TplName = "login.html"
	this.Redirect("/login", 302)

}

func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	dec, _ := base64.StdEncoding.DecodeString(userName)
	if userName != "" {
		this.Data["userName"] = string(dec)
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}

func (this *UserController) HandleLogin() {
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	if userName == "" || pwd == "" {
		fmt.Println("缺少数据内容")
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		fmt.Println("用户不存在")
		this.TplName = "login.html"
		return
	}
	if user.Pwd != pwd {
		fmt.Println("密码输入错误")
		this.TplName = "login.html"
		return
	}
	remember := this.GetString("remember")
	enc := base64.StdEncoding.EncodeToString([]byte(userName))
	if remember == "on" {
		this.Ctx.SetCookie("userName", enc, 60)
	} else {
		this.Ctx.SetCookie("userName", enc, -1)
	}

	this.SetSession("userName", userName)
	this.Redirect("/article/index", 302)

}

func (this *UserController) Logout() {
	this.DelSession("userName")

	this.Redirect("/login", 302)

}
