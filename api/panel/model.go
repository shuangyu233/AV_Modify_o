package panel

import (
	"github.com/goccy/go-json"
	"time"
)

// Security type
const (
	None    = 0
	Tls     = 1
	Reality = 2
)

type NodeInfo struct {
	Id           int
	Type         string
	Security     int
	PushInterval time.Duration
	PullInterval time.Duration
	RawDNS       RawDNS
	Rules        Rules

	// origin
	VAllss      *VAllssNode
	Shadowsocks *ShadowsocksNode
	Trojan      *TrojanNode
	Hysteria    *HysteriaNode
	Common      *CommonNode
}

// VAllssNode is vmess and vless node info
type VAllssNode struct {
	CommonNode
	Tls                 int             `json:"tls"`
	TlsSettings         TlsSettings     `json:"tls_settings"`
	TlsSettingsBack     *TlsSettings    `json:"tlsSettings"`
	Network             string          `json:"network"`
	NetworkSettings     json.RawMessage `json:"network_settings"`
	NetworkSettingsBack json.RawMessage `json:"networkSettings"`
	ServerName          string          `json:"server_name"`

	// vless only
	Flow          string        `json:"flow"`
	RealityConfig RealityConfig `json:"-"`
}
type ShadowsocksNode struct {
	CommonNode
	Cipher    string `json:"cipher"`
	ServerKey string `json:"server_key"`
}
type TrojanNode CommonNode
type HysteriaNode struct {
	CommonNode
	UpMbps       int    `json:"up_mbps"`
	DownMbps     int    `json:"down_mbps"`
	Obfs         string `json:"obfs"`
	ObfsPassword string `json:"obfs_password"`
}
type CommonNode struct {
	Host       string      `json:"host"`
	ServerPort int         `json:"server_port"`
	ServerName string      `json:"server_name"`
	Routes     []Route     `json:"routes"`
	BaseConfig *BaseConfig `json:"base_config"`
}
type Route struct {
	Id          int         `json:"id"`
	Match       interface{} `json:"match"`
	Action      string      `json:"action"`
	ActionValue string      `json:"action_value"`
}
type BaseConfig struct {
	PushInterval any `json:"push_interval"`
	PullInterval any `json:"pull_interval"`
}

type TlsSettings struct {
	ServerName string `json:"server_name"`
	ServerPort string `json:"server_port"`
	ShortId    string `json:"short_id"`
	PrivateKey string `json:"private_key"`
}

type RealityConfig struct {
	Xver         uint64 `json:"Xver"`
	MinClientVer string `json:"MinClientVer"`
	MaxClientVer string `json:"MaxClientVer"`
	MaxTimeDiff  string `json:"MaxTimeDiff"`
}

type RawDNS struct {
	DNSMap  map[string]map[string]interface{}
	DNSJson []byte
}

type Rules struct {
	Regexp   []string
	Protocol []string
}
type UserInfo struct {
	Id            int    `json:"id"`
	Uuid          string `json:"uuid"`
	SpeedLimit    int    `json:"speed_limit"`
	NodeConnector int64
}
type UserTraffic struct {
	UID           int
	Upload        int64
	Download      int64
	NodeConnector int64
}
type NodeStatus struct {
	ID     uint32 `json:"id"`
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime uint64
}
type OnlineUser struct {
	UID int
	IP  string
}
