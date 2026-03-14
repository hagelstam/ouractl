package api

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		baseURL:    "https://api.ouraring.com",
		token:      token,
		httpClient: &http.Client{Timeout: 7 * time.Second},
	}
}

func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf(
			"unauthorized: token is invalid or expired, run `ouractl auth login`",
		)
	case http.StatusForbidden:
		return nil, fmt.Errorf("forbidden: insufficient permissions or expired Oura subscription")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limited: too many requests, try again later")
	default:
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}
}

type pagedResponse[T any] struct {
	Data      []T     `json:"data"`
	NextToken *string `json:"next_token"`
}

func getPaged[T any](c *Client, path string, baseParams url.Values) ([]T, error) {
	var all []T
	var nextToken string

	for {
		params := url.Values{}
		maps.Copy(params, baseParams)
		if nextToken != "" {
			params.Set("next_token", nextToken)
		}

		body, err := c.Get(path, params)
		if err != nil {
			return nil, err
		}

		var resp pagedResponse[T]
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		all = append(all, resp.Data...)

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = *resp.NextToken
	}

	return all, nil
}
