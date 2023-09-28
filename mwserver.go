package main

import (
	"fmt"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"

	// 解决日志乱码
	"github.com/mattn/go-colorable"
	// let's encrypt SSL
	//"github.com/gin-gonic/autotls"
	//"bytes"
	"encoding/json"
	"io"
	//"io/ioutil"
	"log"
	//"net/http"
	"os"
	//"strings"
	//"unsafe"
)

type Nationality struct {
	Country     string
	CountryAbbr string
	Language    string
	// add a new field
	Lang string
}

type Person struct {
	Name   string
	Age    int
	Gender string
	Nation Nationality
	// add a new field
	Scholar string
}

// js处理msg payload和元数据
func main() {
	// 启用gin的日志输出带颜色
	gin.ForceConsoleColor()
	// 替换默认Writer（关键步骤）,解决日志乱码问题
	//gin.DefaultWriter = colorable.NewColorableStdout()
	// 写日志到文件和控制台
	logfile, _ := os.Create("wm.gin.log")
	gin.DefaultWriter = io.MultiWriter(logfile, colorable.NewColorableStdout())
	// 定义日志格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r := gin.Default()

	// init the rule engine
	config := rulego.NewConfig()
	metaData := types.NewMetadata()

	//js处理
	ruleEngine, err := rulego.New("rule01", []byte(chainJsonFile2), rulego.WithConfig(config))
	if err != nil {
		panic(err)
	}

	r.POST("/msg", func(c *gin.Context) {
		var person Person
		// method1
		err = c.BindJSON(&person)

		fmt.Println(person)

		// method2
		//body, err := ioutil.ReadAll(c.Request.Body)
		//if err != nil {
		//	fmt.Println("Nothing Got")
		//	c.JSON(200,gin.H{"rValue": "Noting Got"})
		//}

		//json.Unmarshal(body,&person)
		//fmt.Println(person)

		person.Scholar = "Bachelor"
		person.Nation.Lang = "zh"
		fmt.Println(person)
		fmt.Println(person.Name)
		fmt.Println(person.Nation.Country)

		// convert struct to json
		jperson, err := json.Marshal(person)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(200, gin.H{"error": err.Error()})
		}

		fmt.Println(string(jperson))

		metaData.PutValue("Transform_Msg", "From_Client")

		msg1 := types.NewMsg(0, "TEST_MSG_TYPE1", types.JSON, metaData, string(jperson))

		ruleEngine.OnMsgWithOptions(msg1, types.WithEndFunc(func(msg types.RuleMsg, err error) {
			fmt.Println("msg1处理结果=====")
			//得到规则链处理结果
			fmt.Println(msg.Data)
			fmt.Println(msg.Id)
			fmt.Println(msg.Ts)
			fmt.Println(msg.Type)
			fmt.Println(msg.DataType)
			fmt.Println(msg.Metadata)
			fmt.Println(msg, err)
		}))

		c.JSON(200, gin.H{"status": "arrive at mw"})
	})

	//time.Sleep(time.Second * 10)

	r.Run(":8082")
}

var chainJsonFile1 = `
{
  "ruleChain": {
	"id":"rule01",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
       {
        "id": "s1",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['name']='test01';\n metadata['index']=11;\n msg['addField']='addValue1'; return {'msg':msg,'metadata':metadata,'msgType':msgType};"
        }
      }
    ],
    "connections": [
    ],
    "ruleChainConnections": null
  }
}
`

var chainJsonFile2 = `
{
  "ruleChain": {
    "id":"rule02",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
       {
        "id": "s1",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['name']='test02';\n metadata['index']=22;\n msg['Mobile']='+8613800138000';msg['State']='BeiJing'; return {'msg':msg,'metadata':metadata,'msgType':msgType};"
        }
      },
      {
        "id": "s2",
        "type": "restApiCall",
        "name": "推送数据",
        "debugMode": true,
        "configuration": {
          "restEndpointUrlPattern": "http://localhost:8080/server",
          "requestMethod": "POST",
          "maxParallelRequestsCount": 200
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "s2",
        "type": "Success"
      }
    ],
    "ruleChainConnections": null
  }
}
`
