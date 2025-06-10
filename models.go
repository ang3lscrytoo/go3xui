package go3xui

import (
	"encoding/json"
	"fmt"
)

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     T      `json:"obj,omitempty"`
}

func (response *APIResponse[T]) Err(method string) error {
	return APIError{
		Message: response.Msg,
		Method:  method,
	}
}

type APIError struct {
	Message string
	Method  string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s | Ошибка - %s", e.Method, e.Message)
}

type Inbound struct {
	ID             int             `json:"id"`
	Up             int             `json:"up"`
	Down           int             `json:"down"`
	Total          int             `json:"total"`
	Remark         string          `json:"remark"`
	Enable         bool            `json:"enable"`
	ExpiryTime     int64           `json:"expiryTime"`
	ClientStats    interface{}     `json:"clientStats"`
	Listen         string          `json:"listen"`
	Port           int             `json:"port"`
	Protocol       string          `json:"protocol"`
	Settings       *Settings       `json:"settings"`
	StreamSettings *StreamSettings `json:"streamSettings"`
	Tag            string          `json:"tag"`
	Sniffing       *Sniffing       `json:"sniffing"`
	Allocate       *Allocate       `json:"allocate"`
}

type Settings struct {
	Clients    []*InboundClient `json:"clients"`
	Decryption string           `json:"decryption"`
	Fallbacks  []interface{}    `json:"fallbacks"`
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return json.Unmarshal([]byte(str), s)
	}

	type Alias Settings
	aux := &struct{ *Alias }{Alias: (*Alias)(s)}
	return json.Unmarshal(data, aux)
}

type InboundClient struct {
	Comment    string `json:"comment"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Enable     bool   `json:"enable"`
	ExpiryTime int64  `json:"expiryTime,omitempty"`
	Flow       string `json:"flow,omitempty"`
	ID         string `json:"id,omitempty"`
	LimitIp    int    `json:"limitIp,omitempty"`
	Reset      int    `json:"reset,omitempty"`
	SubId      string `json:"subId,omitempty"`
	TgId       any    `json:"tgId,omitempty"`
	TotalGB    int    `json:"totalGB,omitempty"`
}

type StreamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	ExternalProxy   []interface{}   `json:"externalProxy"`
	RealitySettings RealitySettings `json:"realitySettings"`
	TcpSettings     TcpSettings     `json:"tcpSettings"`
}

type RealitySettings struct {
	Show        bool     `json:"show"`
	Xver        int      `json:"xver"`
	Dest        string   `json:"dest"`
	ServerNames []string `json:"serverNames"`
	PrivateKey  string   `json:"privateKey"`
	MinClient   string   `json:"minClient"`
	MaxClient   string   `json:"maxClient"`
	MaxTimediff int      `json:"maxTimediff"`
	ShortIds    []string `json:"shortIds"`
	Settings    struct {
		PublicKey   string `json:"publicKey"`
		Fingerprint string `json:"fingerprint"`
		ServerName  string `json:"serverName"`
		SpiderX     string `json:"spiderX"`
	} `json:"settings"`
}

type TcpSettings struct {
	AcceptProxyProtocol bool `json:"acceptProxyProtocol"`
	Header              struct {
		Type string `json:"type"`
	} `json:"header"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
	MetadataOnly bool     `json:"metadataOnly"`
	RouteOnly    bool     `json:"routeOnly"`
}

type Allocate struct {
	Strategy    string `json:"strategy"`
	Refresh     int    `json:"refresh"`
	Concurrency int    `json:"concurrency"`
}

type ServerStatus struct {
	Cpu         float64 `json:"cpu"`
	CpuCores    int     `json:"cpuCores"`
	CpuSpeedMhz float64 `json:"cpuSpeedMhz"`
	Mem         struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"mem"`
	Swap struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"swap"`
	Disk struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"disk"`
	Xray struct {
		State    string `json:"state"`
		ErrorMsg string `json:"errorMsg"`
		Version  string `json:"version"`
	} `json:"xray"`
	Uptime   uint64    `json:"uptime"`
	Loads    []float64 `json:"loads"`
	TcpCount int       `json:"tcpCount"`
	UdpCount int       `json:"udpCount"`
	NetIO    struct {
		Up   uint64 `json:"up"`
		Down uint64 `json:"down"`
	} `json:"netIO"`
	NetTraffic struct {
		Sent uint64 `json:"sent"`
		Recv uint64 `json:"recv"`
	} `json:"netTraffic"`
	PublicIP struct {
		IPv4 string `json:"ipv4"`
		IPv6 string `json:"ipv6"`
	} `json:"publicIP"`
	AppStats struct {
		Threads uint32 `json:"threads"`
		Mem     uint64 `json:"mem"`
		Uptime  uint64 `json:"uptime"`
	} `json:"appStats"`
}

func (ss *StreamSettings) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return json.Unmarshal([]byte(str), ss)
	}

	type Alias StreamSettings
	aux := &struct{ *Alias }{Alias: (*Alias)(ss)}
	return json.Unmarshal(data, aux)
}

func (sn *Sniffing) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return json.Unmarshal([]byte(str), sn)
	}

	type Alias Sniffing
	aux := &struct{ *Alias }{Alias: (*Alias)(sn)}
	return json.Unmarshal(data, aux)
}

func (a *Allocate) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return json.Unmarshal([]byte(str), a)
	}

	type Alias Allocate
	aux := &struct{ *Alias }{Alias: (*Alias)(a)}
	return json.Unmarshal(data, aux)
}
