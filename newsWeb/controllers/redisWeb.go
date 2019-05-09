package controllers
import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func init(){
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		fmt.Println("连接错误:",err)
		return
	}
	defer conn.Close()

	resp,err:=conn.Do("mget","s1","s2","s3")

	result,_:=redis.Values(resp,err)
	var v1,v3 int
	var v2 string
	redis.Scan(result,&v1,&v2,&v3)
	fmt.Println("获取结果为:",v1,v2,v3)

}
