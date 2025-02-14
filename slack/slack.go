package slack

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	slackAPI = "https://slack.com/api/"
)

type Client struct {
	Token string
	HTTP  http.Client
}

type slackResponse struct {
	OK    bool        `json:"ok"`
	Error string      `json:"error"`
	User  userDetails `json:"user"`
}

type userDetails struct {
	ID string `json:"id"`
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
	values := url.Values{
		"channel": {target},
		"text":    {message},
	}
	_, err := client.post("chat.postMessage", values)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) GetUserByEmail(ctx context.Context, email string) (string, error) {
	log.Printf("Looking up slack user for %s", email)
	values := url.Values{
		"email": {email},
	}
	resp, err := client.get("users.lookupByEmail", values)
	user := resp.User.ID
	return user, err
}

func (client *Client) post(method string, values url.Values) (*slackResponse, error) {
	values.Set("token", client.Token)
	res, err := client.HTTP.PostForm(slackAPI+method, values)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send slack message")
	}

	resp := &slackResponse{}
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read slack response")
	}

	if !resp.OK {
		return resp, errors.Errorf("error from slack api [%s]", resp.Error)
	}
	return resp, nil
}

func (client *Client) get(method string, values url.Values) (*slackResponse, error) {
	values.Set("token", client.Token)
	query := values.Encode()

	res, err := client.HTTP.Get(slackAPI + method + "?" + query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send slack message")
	}

	resp := &slackResponse{}
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read slack response")
	}

	if !resp.OK {
		return resp, errors.Errorf("error from slack api [%s]", resp.Error)
	}
	return resp, nil
}
