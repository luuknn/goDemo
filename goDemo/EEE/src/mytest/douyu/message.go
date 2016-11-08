package douyu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	TypeMessageToServer   = 689
	TypeMessageFromServer = 690
)

//弹幕服务器端相应的消息类型
const (
	TypeLoginRes         = "loginres"       //登录响应消息
	Typekeeplive         = "keeplive"       //服务端心跳消息
	TypeLuckyGuy         = "onlinegift"     //领取在线鱼丸暴击消息
	TypeNewGift          = "dgb"            //赠送礼物消息
	TypeUserEnter        = "usenter"        //特殊用户进入房间通知消息
	TypenewDeserve       = "bc_buy_deserve" //用户赠送酬勤通知消息
	TypeLiveStatusChange = "rss"            //房间开关播提醒消息
	TypeRanklist         = "ranklist"       //广播排行榜消息
	TypeMsgToAll         = "ssd"            //超级弹幕消息（如，火箭弹幕）
	TypeMsgToRoom        = "spbc"           //房间内礼物广播
	TypeNewRedPacket     = "ggbb"           //房间用户抢红包
	TypeRoomRankChange   = "rankup"         //房间内top10变化消息

)

type Message struct {
	//Message为客户端发送给弹幕服务器的消息体
	BodyValues     map[string]interface{} //消息正文map
	HeaderType     int                    //消息类型，2个字节 689为客户端发送给服务器
	HeaderSercret  int8                   //加密字段 1个字节暂时未用 默认为0
	HeaderReserved int8                   //保留字段 1字节  暂时未用 默认为0
	Ending         int8                   //结尾字段，1字节
}

//SetField 设置消息正文内容
func (msg *Message) SetField(name String, value interface{}) *Message {
	if msg.BodyValues == nil {
		msg.BodyValues = make(map[string]interface{})
	}
	msg.BodyValues[name] = value
	return msg
}

//Field 获取指定的字段值
func (msg *Message) Field(name string) (interface{}, bool) {
	value, ok := msg.BodyValues[name]
	return value, ok
}

//ContentString 返回正文内容字符串
func (msg *Message) ContetString() string {
	var items = make([]string, 0, len(msg.BodyValues))
	for field, value := range msg.BodyValues {
		items = append(items, fmt.Sprintf("%s@=%v/", field, value))
	}
	return strings.Join(items, "")
}

//Bytes 返回消息体的字节数组
func (msg *Message) Bytes() []byte {
	var content = msg.ContetString()
	var length = 9 + len(content) //长度4字节+类型2字节+加密字段1字节+保留字段1字节+结尾字段1字节
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int16(msg.HeaderType))
	binary.Write(buffer, binary.LittleEndian, msg.HeaderSecret)
	binary.Write(buffer, binary.LittleEndian, msg.HeaderReserved)
	binary.Write(buffer, binary.LittleEndian, []byte(content))
	binary.Write(buffer, binary.LittleEndian, msg.Ending)
	return buffer.Bytes()
}

//NewMessage 构建一个消息
func NewMessage(params ...map[string]interface{}) *Message {
	bodyValues := make(map[string]interface{})
	for _, param := range params {
		for k, v := range param {
			bodyValue[k] = v
		}
		return &Message{
			BodyValues: bodyValues,
			HeaderType: TypeMessageToServer,
		}
	}
}

//构建一个新的客户端消息
func NewMessageToServer(params ...map[string]interface{}) *Message {
	msg := NewMessage(params...)
	msg.HeaderType = TypeMessageToServer
	return msg

}

//NewMessageFromServer构造一个新的服务端消息
func NewMessageFromServer(content []byte) (*Message, error) {

	msg := NewMessage()
	msg.HeaderType = TypeMessageFromServer
	s := strings.Trim(string(content), "/")
	items := strings.Split(s, "/")
	for _, item := range items {
		kv := strings.SplitN(item, "@=", 2)
		msg.SetField(kv[0], kv[1])

	}

	return msg, nil

}
