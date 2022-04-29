package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// webSocket返回text 格式
func textApi(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()
	//读取ws中的数据
	mt, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("error read message")
		log.Fatal(err)
	}

	//写入ws数据, pong  times
	var count = 0
	for {
		count++
		if count > 10 {
			break
		}

		message = []byte(string(message) + " " + strconv.Itoa(count))
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("error write message: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

//webSocket返回 json格式
func jsonApi(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()
	var data struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	//读取ws中的数据
	err = ws.ReadJSON(&data)
	if err != nil {
		log.Println("error read json")
		log.Fatal(err)
	}

	//写入ws数据, pong  times
	var count = 0
	for {
		count++
		if count > 10 {
			break
		}

		err = ws.WriteJSON(struct {
			A string `json:"a"`
			B int    `json:"b"`
			C int    `json:"c"`
		}{
			A: "是的",
			B: data.B,
			C: count,
		})
		if err != nil {
			log.Println("error write json: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func websocketGin() {
	r := gin.Default()
	r.GET("/json", jsonApi)
	r.GET("/text", textApi)

	// static files
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.Run(":8000")
}
func main() {
	websocketGin()
}
