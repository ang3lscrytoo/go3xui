package go3xui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	_login        string = "/login"
	_getInbounds  string = "/panel/api/inbounds/list"
	_getInbound   string = "/panel/api/inbounds/get/{id}"
	_addClient    string = "/panel/api/inbounds/addClient"
	_deleteClient string = "/panel/api/inbounds/{inboundId}/delClient/{uuid}"
	_updateClient string = "/panel/api/inbounds/updateClient/{uuid}"
	_getStatus    string = "/server/status"
)

type XUIClient struct {
	core *XUICore
}

func NewClient(host, username, password string, enableLogger bool) *XUIClient {
	coreClient := &fasthttp.Client{
		DialTimeout: func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, timeout)
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &XUIClient{
		core: &XUICore{httpClient: coreClient, host: host, username: username, password: password, logger: enableLogger},
	}
}

func (c *XUIClient) GetInbounds() ([]*Inbound, error) {
	var response APIResponse[[]*Inbound]

	err := c.core.ApiCall("GET", _getInbounds, nil, &response)

	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, response.Err(_getInbounds)
	}
	return response.Obj, nil
}

func (c *XUIClient) GetInbound(inboundId int) (*Inbound, error) {
	var response APIResponse[Inbound]

	endpoint := strings.Replace(_getInbound, "{id}", strconv.Itoa(inboundId), 1)
	err := c.core.ApiCall("GET", endpoint, nil, &response)

	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, response.Err(endpoint)
	}
	return &response.Obj, nil
}

func (c *XUIClient) AddClient(inboundId int, inboundClient *InboundClient) error {
	var response APIResponse[any]

	bytes, err := json.Marshal(map[string]interface{}{
		"clients": append([]InboundClient{}, *inboundClient)})

	if err != nil {
		return err
	}

	err = c.core.ApiCall("POST", _addClient, map[string]interface{}{
		"id":       inboundId,
		"settings": string(bytes),
	}, &response)

	if err != nil {
		return err
	}
	if !response.Success {
		return response.Err(_addClient)
	}
	return nil
}

func (c *XUIClient) DeleteClient(inboundId int, clientUuid string) error {
	var response APIResponse[any]
	var endpoint string

	endpoint = strings.Replace(_deleteClient, "{inboundId}", strconv.Itoa(inboundId), 1)
	endpoint = strings.Replace(endpoint, "{uuid}", clientUuid, 1)
	err := c.core.ApiCall("POST", endpoint, nil, &response)

	if err != nil {
		return err
	}
	if !response.Success {
		return response.Err(endpoint)
	}
	return nil
}

func (c *XUIClient) GetClientByEmail(email string) (*InboundClient, error) {
	inbounds, err := c.GetInbounds()
	if err != nil {
		return nil, err
	}

	for _, inbound := range inbounds {
		for _, client := range inbound.Settings.Clients {
			if client.Email == email {
				return client, nil
			}
		}
	}

	return nil, errors.New("client not found")
}

func (c *XUIClient) UpdateClient(inboundId int, clientUuid string, inboundClient *InboundClient) error {
	var response APIResponse[any]

	bytes, err := json.Marshal(map[string]interface{}{
		"clients": append([]InboundClient{}, *inboundClient)})

	if err != nil {
		return err
	}

	endpoint := strings.Replace(_updateClient, "{uuid}", clientUuid, 1)
	err = c.core.ApiCall("POST", endpoint, map[string]interface{}{
		"id":       inboundId,
		"settings": string(bytes),
	}, &response)

	if err != nil {
		return err
	}
	if !response.Success {
		return response.Err(endpoint)
	}
	return nil
}

func (c *XUIClient) GetStatus() (*ServerStatus, error) {
	var response APIResponse[ServerStatus]

	err := c.core.ApiCall("POST", _getStatus, nil, &response)

	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, response.Err(_getStatus)
	}
	return &response.Obj, nil
}

func (c *XUIClient) GetIP() (string, error) {
	parsedURL, err := url.Parse(c.core.host)
	if err != nil {
		fmt.Println("Ошибка парсинга URL:", err)
		return "", errors.New("incorrect URL host")
	}

	host := parsedURL.Hostname()
	return host, nil
}
