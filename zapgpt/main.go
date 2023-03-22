package zapgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

func GenerateGPTTtext(query string) (string, error) {
	req := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens: 150,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer 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")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func sendMessage(c *gin.Context) {
	body := c.Request.Body
	buf := make([]byte, 1024)
	num, _ := body.Read(buf)

	text := string(buf[0:num])
	result, err := GenerateGPTTtext(text)
	if err != nil {
		panic(err)
	}

	c.String(http.StatusOK, result)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/chat", sendMessage)
	router.Run(":8080")
}
