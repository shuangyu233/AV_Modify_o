package panel

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strings"
	"time"
)

const (
	GET_NODE_INFO_PATH  = "/api/public/airgo/node/getNodeInfo"
	GET_USER_LIST_PATH  = "/api/public/airgo/user/getUserlist"
	REPORT_USER_TRAFFIC = "/api/public/airgo/user/reportUserTraffic"
	REPORT_NODE_STATUS  = "/api/public/airgo/node/reportNodeStatus"
	REPORT_ONLINE_USERS = "/api/public/airgo/user/reportNodeOnlineUsers"
)

func (c *Client) GetNodeInfo() (node *NodeInfo, err error) {
	r, err := c.client.
		R().
		SetQueryParams(map[string]string{
			"id": fmt.Sprintf("%d", c.NodeId),
		}).
		SetHeader("If-None-Match", c.nodeEtag).
		Get(GET_NODE_INFO_PATH)
	if err = c.checkResponse(r, GET_NODE_INFO_PATH, err); err != nil {
		return
	}
	if r.StatusCode() == 304 {
		return nil, nil
	}
	c.nodeEtag = r.Header().Get("ETag")
	var nodeInfoResponse NodeInfoResponse
	err = json.Unmarshal(r.Body(), &nodeInfoResponse)
	if err != nil {
		return nil, err
	}
	node, err = c.ParseAirGoNodeInfo(&nodeInfoResponse)
	return
}
func (c *Client) ParseAirGoNodeInfo(n *NodeInfoResponse) (*NodeInfo, error) {
	node := &NodeInfo{
		Id:   int(n.ID),
		Type: c.NodeType,
		//Security:     0,
		PushInterval: 30 * time.Second,
		PullInterval: 30 * time.Second,
		RawDNS: RawDNS{
			DNSMap:  make(map[string]map[string]any),
			DNSJson: []byte(""),
		},
		Rules:       Rules{},
		VAllss:      nil,
		Shadowsocks: nil,
		Trojan:      nil,
		Hysteria:    nil,
		//Common: &CommonNode{},
	}
	switch n.Security {
	case "none", "":
		node.Security = 0
	case "tls":
		node.Security = 1
	case "reality":
		node.Security = 2
	default:
		node.Security = 0
	}
	switch c.NodeType {
	case "vmess", "vless":
		request := map[string]any{
			"path": n.Path,
			"headers": map[string]any{
				"Host": n.Host,
			}}
		header, _ := json.Marshal(request)
		var sName, sPort string
		if n.Dest != "" {
			sName = n.Dest[:strings.Index(n.Dest, ":")]
			sPort = n.Dest[strings.Index(n.Dest, ":")+1:]
		}

		node.VAllss = &VAllssNode{
			CommonNode: CommonNode{
				Host:       "",
				ServerPort: int(n.Port),
				ServerName: "",
				Routes:     nil,
				BaseConfig: nil,
			},
			//Tls:                 0, //omit
			TlsSettings: TlsSettings{
				ServerName: sName,
				ServerPort: sPort,
				ShortId:    n.ShortId,
				PrivateKey: n.PrivateKey,
			},
			TlsSettingsBack: nil,
			Network:         n.Network,
			NetworkSettings: header,
			//NetworkSettingsBack: nil,//omit
			//ServerName:    "", //omit
			Flow: n.VlessFlow,
			//RealityConfig: RealityConfig{},  //omit
		}
		node.Common = &node.VAllss.CommonNode
	case "shadowsocks":
		node.Shadowsocks = &ShadowsocksNode{
			CommonNode: CommonNode{
				Host:       "",
				ServerPort: int(n.Port),
				ServerName: "",
				Routes:     nil,
				BaseConfig: nil,
			},
			Cipher:    n.Scy,
			ServerKey: n.ServerKey,
		}
		node.Common = &node.Shadowsocks.CommonNode

	case "hysteria2":
		node.Security = 1
		node.Hysteria = &HysteriaNode{
			CommonNode: CommonNode{
				Host:       n.Address,
				ServerPort: int(n.Port),
				ServerName: n.Sni,
				Routes:     nil,
				BaseConfig: nil,
			},
			UpMbps:       int(n.HyUpMbps),
			DownMbps:     int(n.HyDownMbps),
			Obfs:         n.HyObfs,
			ObfsPassword: n.HyObfsPassword,
		}
		node.Common = &node.Hysteria.CommonNode

	}
	for _, v := range n.Access {
		if strings.Index(v.Route, "Protocol") == -1 {
			node.Rules.Regexp = append(node.Rules.Regexp, strings.TrimSpace(v.Route))
		} else {
			node.Rules.Protocol = append(node.Rules.Regexp, strings.TrimSpace(v.Route))
		}
	}
	return node, nil
}
func (c *Client) GetUserList() (UserList []UserInfo, err error) {
	r, err := c.client.
		R().
		SetQueryParams(map[string]string{
			"id": fmt.Sprintf("%d", c.NodeId),
		}).
		SetHeader("If-None-Match", c.userEtag).
		Get(GET_USER_LIST_PATH)
	if err = c.checkResponse(r, GET_USER_LIST_PATH, err); err != nil {
		return
	}
	if r.StatusCode() == 304 {
		return nil, nil
	}
	c.userEtag = r.Header().Get("ETag")
	var userResponse []UserResponse
	json.Unmarshal(r.Body(), &userResponse)
	for _, v := range userResponse {
		UserList = append(UserList, UserInfo{
			Id:            int(v.ID),
			Uuid:          v.UUID,
			SpeedLimit:    int(v.NodeSpeedLimit),
			NodeConnector: v.NodeConnector,
		})
	}
	return
}
func (c *Client) ReportUserTraffic(userTraffic []UserTraffic) error {
	var userTrafficRequest = UserTrafficRequest{
		ID:          c.NodeId,
		UserTraffic: nil,
	}
	for _, v := range userTraffic {
		userTrafficRequest.UserTraffic = append(userTrafficRequest.UserTraffic, UserTrafficItem{
			UID:      v.UID,
			Email:    "",
			Upload:   v.Upload,
			Download: v.Download,
		})
	}
	r, err := c.client.
		R().
		SetBody(userTrafficRequest).
		Post(REPORT_USER_TRAFFIC)
	if r.StatusCode() == 200 {
		return nil
	}
	return err
}
func (c *Client) ReportNodeStatus() error {
	var nodeStatus NodeStatus

	infocpu, _ := cpu.Percent(time.Duration(time.Second), false)
	infomem, _ := mem.VirtualMemory()
	infodisk, _ := disk.Usage(".")

	nodeStatus.CPU = infocpu[0]
	nodeStatus.Mem = infomem.UsedPercent
	nodeStatus.Disk = infodisk.UsedPercent

	var nodeStatusRequest = NodeStatusRequest{
		ID:     c.NodeId,
		CPU:    nodeStatus.CPU,
		Mem:    nodeStatus.Mem,
		Disk:   nodeStatus.Disk,
		Uptime: nodeStatus.Uptime,
	}
	res, err := c.client.R().
		SetBody(nodeStatusRequest).
		ForceContentType("application/json").
		Post(REPORT_NODE_STATUS)
	if res.StatusCode() == 200 {
		return nil
	}
	return err
}
func (c *Client) ReportNodeOnlineUsers(data *map[int][]string) error {
	//var reqData = OnlineUserRequest{
	//	NodeID:      c.NodeId,
	//	UserNodeMap: *data,
	//}
	//res, err := c.client.R().
	//	SetBody(reqData).
	//	ForceContentType("application/json").
	//	Post(REPORT_ONLINE_USERS)
	//if res.StatusCode() == 200 {
	//	return nil
	//}
	//return err
	return nil
}
