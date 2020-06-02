package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net-nes/game"
	"net/http"
)

// 游戏房间实例
var room *game.Room

// 初始路由
func InitRouter() *gin.Engine {

	// 初始化游戏
	room = game.Step()

	// 路由实例
	r := gin.Default()

	// 静态资源
	r.Static("/js", "./static/js")
	r.Static("/roms", "./static/roms")
	r.LoadHTMLGlob("views/*")

	// 主页
	r.GET("/", index)

	// 游戏通信
	r.GET("ws", ws)

	return r
}

// 主页
func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

// websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize: 32 * 1024,
	WriteBufferSize: 32 * 1024,
}

// 游戏通信
func ws(ctx *gin.Context) {

	// websocket实例
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrader:", err)
		return
	}

	// 玩家实例
	people := &game.People{
		Conn: conn,
		SendChan: make(chan []byte, 256),
	}

	// 进入游戏房间
	room.EnterChan <- people
}