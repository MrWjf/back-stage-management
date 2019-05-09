package controllers

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"newsWeb/models"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowIndex() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	this.Data["userName"] = userName.(string)
	//获取所有文章数据，展示到页面
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)
	typeName := this.GetString("select")
	var count int64
	if typeName == "" {
		//获取总记录数
		count, _ = qs.RelatedSel("ArticleType").Count()
	} else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}
	//获取总页数
	pageIndex := 2

	pageCount := math.Ceil(float64(count) / float64(pageIndex))
	//获取首页和末页数据
	//获取页码
	pageNum, err := this.GetInt("pageNum")
	if err != nil {
		pageNum = 1
	}
	beego.Info("数据总页数为:", pageNum)

	if typeName == "" {
		//获取对应页的数据   获取几条数据     起始位置
		qs.Limit(pageIndex, pageIndex*(pageNum-1)).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageIndex, pageIndex*(pageNum-1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)

	}
	var articleTypes []models.ArticleType
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		fmt.Println("连接错误",err)
		return
	}
	defer conn.Close()
	resp,err:=conn.Do("get","newsWeb")
	result,_:=redis.Bytes(resp,err)

	if len(result)==0{
	o.QueryTable("ArticleType").All(&articleTypes)


	var buffer bytes.Buffer
	enc:=gob.NewEncoder(&buffer)
	enc.Encode(articleTypes)
	conn.Do("set","newsWeb",buffer.Bytes())
	fmt.Println("从Mysql获取数据")

	}else{
		dec:=gob.NewDecoder(bytes.NewReader(result))
		dec.Decode(&articleTypes)
		fmt.Println(articleTypes)
		fmt.Println("从redis获取数据")

	}
	this.Data["articleTypes"] = articleTypes
	this.Data["typeName"] = typeName //传给后台$.TypeName
	this.Data["articles"] = articles
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = pageNum

	this.Layout = "layout.html"

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["indexJs"] = "indexJs.html"

	this.TplName = "index.html"
}

func (this *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes

	this.Layout = "layout.html"

	this.TplName = "add.html"
}

func (this *ArticleController) HandleAddArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")
	if articleName == "" || content == "" || typeName == "" {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = ("add.html")
		return
	}
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = ("add.html")
		return
	}
	defer file.Close()

	if head.Size > 5000000 {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片数据过大请重新上传"
		this.TplName = ("add.html")
		return
	}

	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != "ipeg" {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片格式错误"
		this.TplName = ("add.html")
		return
	}

	fileName := time.Now().Format("200601021504052222")
	this.SaveToFile("uploadname", "./static/img/"+fileName+ext)

	o := orm.NewOrm()

	var article models.Article
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/" + fileName + ext

	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")

	article.ArticleType = &articleType

	_, err = o.Insert(&article)
	if err != nil {
		fmt.Println("数据插入失败", err)
		this.Data["errmsg"] = "数据插入失败"
		this.TplName = ("add.html")
		return
	}
	this.Redirect("/article/index", 302)
}

func (this *ArticleController) ShowContent() {
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("查询错误", err)
		this.Redirect("/article/index", 302)
		return
	}

	o := orm.NewOrm()

	var article models.Article
	article.Id = id

	o.Read(&article)

	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)
	this.Data["users"] = users

	article.ReadCount += 1
	o.Update(&article)

	this.Data["article"] = article

	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	m2m := o.QueryM2M(&article, "Users")
	m2m.Add(user)

	this.Layout = "layout.html"

	this.TplName = "content.html"
}

func (this *ArticleController) ShowUpdate() {
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("查询错误", err)
		this.Redirect("/article/index", 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	this.Data["article"] = article
	this.TplName = "update.html"
}

func UploadFile(this *ArticleController, filePath string, errHtml string) string {
	file, head, err := this.GetFile(filePath)
	if err != nil {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = errHtml
		return ""
	}
	defer file.Close()

	if head.Size > 5000000 {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片数据过大请重新上传"
		this.TplName = errHtml
		return ""
	}

	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != "ipeg" {
		fmt.Println("获取数据错误")
		this.Data["errmsg"] = "图片格式错误"
		this.TplName = errHtml
		return ""
	}

	fileName := time.Now().Format("200601021504052222")
	this.SaveToFile(filePath, "./static/img/"+fileName+ext)
	return "/static/img/" + fileName + ext

}

func (this *ArticleController) HandleUpdate() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	savePath := UploadFile(this, "uploadname", "update.html")
	id, _ := this.GetInt("id") //隐藏域传值
	//校验数据
	if articleName == "" || content == "" || savePath == "" {
		beego.Error("获取数据失败")
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	article.Title = articleName
	article.Content = content
	article.Img = savePath
	o.Update(&article)
	this.Redirect("/article/index", 302)

}

func (this *ArticleController) HandleDelete() {
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("获取ID错误", err)
		this.Redirect("/article/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article, "Id")

	this.Redirect("/article/index", 302)

}

func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	//返回数据
	this.Data["articleTypes"] = articleTypes

	this.Layout = "layout.html"

	this.TplName = "addType.html"
}

func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Error("类型名称传输失败")
		this.Redirect("/article/addType", 302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)
	this.Redirect("/article/addType", 302)
}

func (this *ArticleController) DeleteType() {
	id,err:=this.GetInt("id")
	if err!=nil{
		fmt.Println("获取文章id失败",err)
		this.Redirect("/article/addType",302)
		return
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=id
	o.Delete(&articleType,"Id")
	this.Redirect("/article/addType",302)

}
