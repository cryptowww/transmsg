package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// 解决日志乱码
	"github.com/mattn/go-colorable"
	// let's encrypt SSL
	//"github.com/gin-gonic/autotls"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unsafe"
)

type Nationality struct {
	Country     string
	CountryAbbr string
	Language    string
}

type Person struct {
	Name   string
	Age    int
	Gender string
	Nation Nationality
}

func main() {
	// 启用gin的日志输出带颜色
	gin.ForceConsoleColor()
	// 替换默认Writer（关键步骤）,解决日志乱码问题
	//gin.DefaultWriter = colorable.NewColorableStdout()
	// 写日志到文件和控制台
	logfile, _ := os.Create("client.gin.log")
	gin.DefaultWriter = io.MultiWriter(logfile, colorable.NewColorableStdout())
	// 定义日志格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r := gin.Default()

	r.POST("/send", func(c *gin.Context) {
		//body, _ := ioutil.ReadAll(c.Request.Body)
		//fmt.Println(string(body))

		person := &Person{
			"James",
			30,
			"male",
			Nationality{
				"China",
				"CN",
				"Chiness Mandarin",
			},
		}

		byteData, err := json.Marshal(person)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(string(byteData))

		reader := bytes.NewReader(byteData)
		fmt.Println(reader)

		url := "http://localhost:8082/msg"
		request, err := http.NewRequest("POST", url, reader)
		defer request.Body.Close()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}

		response, err := client.Do(request)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		respBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		respStr := (*string)(unsafe.Pointer(&respBytes))
		fmt.Println(*respStr)
		c.JSON(200, respStr)
	})

	r.Run(":8081")
}
