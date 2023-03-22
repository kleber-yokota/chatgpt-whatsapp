package chatgptwhatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("ChatWhatsApp", chatWhatsApp)
}

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
	Index   int      `json:"index"`
	Message struct { // possible to create struct in struct
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func chatWhatsApp(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Problem with Twilio message")
		return
	}
	uri := string(b)
	parameter, err := url.ParseQuery(uri)
	if err != nil {
		fmt.Fprint(w, "Problem URI")
		return
	}
	bodyMessage := parameter.Get("Body")

	req := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: bodyMessage,
			},
		},
		MaxTokens: 150,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Fprint(w, bodyMessage)
		return
	}
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Fprint(w, "Erro to call ChatGPT API")
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer <API KEY>")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Fprint(w, "Problem to Request to ChatGPT API")
		return
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprint(w, "Problem to read ChatGPT Response")
		return
	}

	var resp Response
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		fmt.Fprint(w, string(reqJson))
		return
	}
	fmt.Fprint(w, resp.Choices[0].Message.Content)

}
