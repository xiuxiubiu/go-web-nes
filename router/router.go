package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-web-nes/game"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

// 游戏房间实例
var room = game.NewRoom()

// 初始路由
func InitRouter() *gin.Engine {

	// 初始化游戏
	room = game.NewRoom()

	// 路由实例
	r := gin.Default()

	// 静态资源
	r.Static("/roms", "./roms")
	r.LoadHTMLGlob("views/*")

	// 主页
	r.GET("/", index)

	// 游戏列表
	r.GET("list", list)

	// 游戏通信
	r.GET("ws", ws)

	return r
}

// 主页
func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

// 游戏列表
func list(ctx *gin.Context) {

	paths, _ := filepath.Glob("./roms/*.nes")

	files := make([]string, len(paths))
	for k, s := range paths {
		files[k] = path.Base(s)
	}

	ctx.JSON(http.StatusOK, files)
}

// websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  32 * 1024,
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
	people := game.NewPeople(conn)

	// 进入游戏房间
	room.Enter <- people
}
