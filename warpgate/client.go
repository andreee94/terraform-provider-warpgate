package warpgate

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	models "terraform-provider-warpgate/warpgate/models"
)

type WarpgateClient struct {
	*ClientWithResponses
	Address    string
	Port       int
	url        string
	httpClient *http.Client
}

func NewWarpgateClient(address string, port int, insecureSkipVerify bool) *WarpgateClient {
	jar, _ := cookiejar.New(nil)

	return &WarpgateClient{
		Address: address,
		Port:    port,
		url:     fmt.Sprintf("https://%s:%d", address, port),
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
			},
			Jar: jar,
		},
	}
}

func (c *WarpgateClient) Login(username string, password string) (err error) {

	login_url := fmt.Sprintf("%s%s", c.url, WARPGATE_ENDPOINT_LOGIN)
	admin_api_url := fmt.Sprintf("%s%s", c.url, WARPGATE_ENDPOINT_ADMIN_API)

	json, err := json.Marshal(&models.LoginData{
		Username: username,
		Password: password,
	})

	if err != nil {
		return
	}

	payload := strings.NewReader(string(json))

	req, err := http.NewRequest("POST", login_url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	if err != nil {
		return
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return
	}

	if res.StatusCode != 201 {
		return errors.New("wrong status code response")
	}

	c.ClientWithResponses, err = NewClientWithResponses(admin_api_url, WithHTTPClient(c.httpClient))

	return
}
