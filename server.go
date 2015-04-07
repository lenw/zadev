package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	GHOST_EMOJI  = ":ghost:"
	MONKEY_EMOJI = ":monkey:"
)

type slackRequest struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji"`
}

func slackPost(channel, text string) error {
	sr := slackRequest{
		Text:      text,
		Channel:   channel,
		Username:  "monkey-bot",
		IconEmoji: MONKEY_EMOJI,
	}

	body, err := json.Marshal(&sr)
	if err != nil {
		fmt.Printf("Marshal error %v\n", err)
		return err
	}
	fmt.Printf("%+v => %v\n", sr, string(body))

	req, err := http.NewRequest("POST", REQUEST_URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(resp.Status)

	}

	fmt.Printf("Resp: %+v\n", string(rb))

	return nil
}

func handleSlackRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start the engines")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return
	}
	fmt.Printf("%+v\n", r.PostForm)
	err := slackPost("#waiting-room", fmt.Sprintf("Request to join from %s <%s>", r.PostForm["name"][0], r.PostForm["email"][0]))
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return
	}
}

var REQUEST_URL string

func main() {
	REQUEST_URL = os.Getenv("REQUEST_URL")
	if REQUEST_URL == "" {
		fmt.Println("Error : could not find REQUEST_URL in the environment. ")
		os.Exit(-1)
	}
	http.HandleFunc("/slack-request", handleSlackRequest)
	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.ListenAndServe(":8080", nil)
}
