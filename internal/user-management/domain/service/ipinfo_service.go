package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IPInfoClient interface {
	GetInfo(ip string) (map[string]interface{}, error)
}

type httpIPInfoClient struct {
	apiToken string
}

func NewIPInfoClient(token string) IPInfoClient {
	return &httpIPInfoClient{apiToken: token}
}

func (c *httpIPInfoClient) GetInfo(ip string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, c.apiToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
