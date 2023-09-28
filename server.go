package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// 解决日志乱码
	"github.com/mattn/go-colorable"
	// let's encrypt SSL
	//"github.com/gin-gonic/autotls"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 启用gin的日志输出带颜色
	gin.ForceConsoleColor()
	// 替换默认Writer（关键步骤）,解决日志乱码问题
	//gin.DefaultWriter = colorable.NewColorableStdout()
	// 写日志到文件和控制台
	logfile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logfile, colorable.NewColorableStdout())
	// 定义日志格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r := gin.Default()

	r.POST("/server", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(string(body))
		c.JSON(200, gin.H{"status": "Receive the message"})
	})
	// Query and PostForm
	// curl -d "name=jim&message=good" -X POST "http://localhost:8080/post?id=99999999&page=11111"
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})

	// SecureJSON
	// curl http://localhost:8080/sjson
	r.GET("/sjson", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}
		c.SecureJSON(http.StatusOK, names)
	})

	// 上传单个文件
	// curl -X POST http://localhost:8080/uploads -F "file=@F:/IBM_DB.pdf" -H "Content-Type: multipart/form-data"
	r.MaxMultipartMemory = 8 << 20
	r.POST("/uploads", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		//if err := c.SaveUploadedFile(file, filepath.Join("d:/",file.Filename)); err != nil {
		if err := c.SaveUploadedFile(file, "upload/"+file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err : %s ", err.Error()))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// 上传多个文件
	// curl -X POST http://localhost:8080/uploadm -F "upload[]=@F:/loki-linux-amd64" -F "upload[]=@F:/loki-linux-amd64.zip" -H "Content-Type: multipart/form-data"
	r.POST("/uploadm", func(c *gin.Context) {
		form, _ := c.MultipartForm()

		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			if err := c.SaveUploadedFile(file, filepath.Join("upload/", file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err : %s ", err.Error()))
				continue
			}
		}

		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	// 从Reader读取数据
	r.GET("/getpic", func(c *gin.Context) {
		url := "https://images.cnblogs.com/cnblogs_com/wupeixuan/1186798/o_wallhaven-4d38m0.jpg"
		urls := strings.Split(url, "/")
		filename := urls[len(urls)-1]

		response, err := http.Get(url)
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		defer reader.Close()
		// -----以下是把文件保存到服务器
		out, err1 := os.Create(filepath.Join("upload/", filename))
		if err1 != nil {
			fmt.Println("write err:", err)
		}

		defer out.Close()
		wf := bufio.NewWriter(out)

		_, err = io.Copy(wf, reader)

		if err != nil {
			fmt.Println("write err:", err)
		}
		wf.Flush()

		// ----- curl http://localhost:8080/getpic
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		fmt.Printf("len:%d, type:%s\n", contentLength, contentType)

		extraHeaders := map[string]string{
			"Content-Disposition": "attachment; filename=\"gopher.png\"",
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

	})

	r.Run()
	// SSL
	//log.Fatal(autotls.Run(r,"localhost"))
}
