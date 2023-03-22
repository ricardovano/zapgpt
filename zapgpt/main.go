package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

type HttpReposnse struct {
	Status     string
	StatusCode int
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

func GenerateGPTTtext(query string) (string, error) {
	if query == "" {
		panic("Body vazio!")
	}

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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found!")
		} else {
			panic(err)
		}
	}

	key := viper.GetString("openaikey")
	token := "Bearer" + " " + key

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)

	response, err := http.DefaultClient.Do(request)
	if response.StatusCode != 200 {
		panic(response.Status)
	}
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
	fmt.Println("Service started on localhost:8080")
	router.POST("/chat", sendMessage)
	router.Run(":8080")
}
