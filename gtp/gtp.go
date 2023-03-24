package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"io/ioutil"
	"log"
	"net/http"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []struct {
		Message      Message           `json:"message"`
		FinishReason string            `json:"finish_reason"`
		Index        int               `json:"index"`
	} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Messages       []Message `json:"messages"`
	Temperature      float32 `json:"temperature"`
	
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{
//     "model": "gpt-3.5-turbo",
//     "messages": [{"role": "user", "content": "Say this is a test!"}],
//     "temperature": 0.7
//   }'
func Completions(msg []Message) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            "gpt-3.5-turbo",
		Messages:         msg,
		Temperature:      0.7,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
		  if v.Message.Role == "assistant" {
    		reply = v.Message.Content
    		break
		  }
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
