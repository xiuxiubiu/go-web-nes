package game

import (
	"strings"
)

// 消息分隔符
var sep = " "

type msgType string

const (

	// 加入游戏
	mtJoin msgType = "join"

	// 游戏数据
	mtData msgType = "data"

	// 键盘按下消息
	mtDown msgType = "down"

	// 键盘抬起消息
	mtUp msgType = "up"
)

// 有效的消息类型
var msgTypeMap = map[msgType]interface{}{
	mtJoin: nil,
	mtData: nil,
	mtDown: nil,
	mtUp:   nil,
}

type Message struct {

	// 消息类型
	mt msgType

	// 数据
	data []byte
}

// 获取消息实例
func NewMessage(message []byte) *Message {

	// 用分隔符解析消息
	ma := strings.Split(string(message), sep)

	// 数组第一个元素是消息类型
	// 第二个元素是数据
	if len(ma) != 2 {
		return nil
	}

	mt := msgType(ma[0])

	// 消息类型无效
	if _, ok := msgTypeMap[mt]; !ok {
		return nil
	}

	return &Message{
		mt:   mt,
		data: message,
	}
}
