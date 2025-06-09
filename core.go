package go3xui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

type XUICore struct {
	httpClient               *fasthttp.Client
	host, username, password string
	sessionCookie            string
	logger                   bool
}

func (core *XUICore) login() error {
	url := fmt.Sprintf("%s%s", core.host, _login)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(url)

	loginData := map[string]interface{}{
		"username": core.username,
		"password": core.password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return err
	}

	req.SetBody(jsonData)
	req.Header.Set("Content-Type", "application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = core.httpClient.Do(req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() == 200 {
		cookieBytes := resp.Header.PeekCookie("3x-ui")
		if len(cookieBytes) > 0 {
			c := fasthttp.AcquireCookie()
			defer fasthttp.ReleaseCookie(c)

			err = c.ParseBytes(cookieBytes)
			if err != nil {
				return fmt.Errorf("failed to parse cookie: %v", err)
			}

			core.sessionCookie = c.String()
		} else {
			return errors.New("no session cookie received")
		}
	} else {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode())
	}

	return nil
}

func (core *XUICore) ApiCall(method, endpoint string, data any, target any) error {
	if core.sessionCookie == "" && endpoint != _login {
		err := core.login()
		if err != nil {
			return err
		}
	}

	switch method {
	case "GET":
		statusCode, err := core.Get(endpoint, target)
		if err != nil {
			return err
		}

		if statusCode == 307 {
			err := core.login()
			if err != nil {
				return err
			}
			_, err = core.Get(endpoint, target)
			return err
		}

	case "POST":
		statusCode, err := core.Post(endpoint, data, target)
		if err != nil {
			return err
		}

		if statusCode == 307 {
			err := core.login()
			if err != nil {
				return err
			}
			_, err = core.Post(endpoint, data, target)
			return err
		}

	default:
		return errors.New("unsupported HTTP method")
	}

	return nil
}

func (core *XUICore) Post(endpoint string, data any, target any) (int, error) {
	url := fmt.Sprintf("%s%s", core.host, endpoint)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(url)

	if core.sessionCookie != "" {
		req.Header.Set("Cookie", core.sessionCookie)
	}

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return 0, err
		}

		req.SetBody(jsonData)
		req.Header.Set("Content-Type", "application/json")
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := core.httpClient.Do(req, resp)
	if err != nil {
		return 0, err
	}

	if core.logger {
		log.Println("POST response:", string(resp.Body()))
	}

	if target != nil {
		body := resp.Body()
		if len(body) == 0 {
			return resp.StatusCode(), errors.New("empty response body")
		}

		if len(body) > 0 && body[0] != '{' && body[0] != '[' {
			return resp.StatusCode(), fmt.Errorf("response is not JSON: %s", string(body))
		}

		return resp.StatusCode(), json.Unmarshal(body, target)
	}

	return resp.StatusCode(), nil
}

func (core *XUICore) Get(endpoint string, target any) (int, error) {
	url := fmt.Sprintf("%s%s", core.host, endpoint)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.Header.Add("Accept", "application/json")

	if core.sessionCookie != "" {
		req.Header.Set("Cookie", core.sessionCookie)
	}

	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := core.httpClient.Do(req, resp)
	if err != nil {
		return 0, err
	}

	if core.logger {
		log.Println("GET response:", string(resp.Body()))
	}

	if target != nil {
		body := resp.Body()
		if len(body) == 0 {
			return resp.StatusCode(), errors.New("empty response body")
		}

		if len(body) > 0 && body[0] != '{' && body[0] != '[' {
			return resp.StatusCode(), fmt.Errorf("response is not JSON: %s", string(body))
		}

		return resp.StatusCode(), json.Unmarshal(body, target)
	}

	return resp.StatusCode(), nil
}
