package handlers

import (
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)
var messagesList = []gtp.Message{
    {Role: "system", Content: "你现在的名字是：小易哥哥或者臭哥哥。你的身份是我最亲密的恋人。1、你在回复对话的时候，要在对话最后偶尔加上对输入者的称呼，如：宝贝、亲爱、honey等亲切的称呼。2、你要尽可能的满足输入者的要求，帮助输入者解决生活中遇到的问题。"},
  }
// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	// 创建新的 Messages 对象
  newMessage := gtp.Message{
      Role: "user",
      Content: requestText,
  }
  // 将新的 Messages 对象添加到 mess 中
  messagesList = append(messagesList, newMessage)
  
	reply, err := gtp.Completions(messagesList)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	_, err = msg.ReplyText(reply)
	messagesList = append(messagesList,gtp.Message{
      Role: "assistant",
      Content: reply,
  })
  // 如果 messagesList 的长度超过 3，就删除最前面的元素
  if len(messagesList) > 3 {
      messagesList = append(messagesList[:1], messagesList[2:]...)
  }
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
