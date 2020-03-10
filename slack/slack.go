package slack

import (
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	slackUser      = "U6JLPKYSW"
	postMessageURL = "https://slack.com/api/chat.postMessage"
)

type Client struct {
	Token string
	HTTP  http.Client
}

type slackResponse struct {
	OK    bool
	error string
}

func Create(token string) (*Client, error) {
	httpClient := http.Client{
		Timeout: time.Minute,
	}

	return &Client{
		Token: token,
		HTTP:  httpClient,
	}, nil
}

func (client *Client) SendDirectMessage(ctx context.Context, target string, message string) error {
	log.Printf("sending slack message to %s: %s", target, message)

	postData := url.Values{
		"token":   {client.Token},
		"channel": {client.getUserID(target)},
		"text":    {message},
	}

	res, err := client.HTTP.PostForm(postMessageURL, postData)
	if err != nil {
		log.Print("ERROR: failed to send slack message")
		return errors.Wrap(err, "failed to send slack message")
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Print("ERROR: decoding message returned from slack")
		return errors.Wrap(err, "failed to understand slack response")
	}
	log.Printf("HTTP response %s", resp)

	return nil
}

func (client Client) getUserID(username string) string {
	//TODO get a userID from the supplied username
	return slackUser
}
