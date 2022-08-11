package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/didi/nightingale/v5/src/models"
	"github.com/koding/multiconfig"
	"github.com/toolkits/pkg/logger"
)

// N9E complete
type N9EPlugin struct {
	Name        string
	Description string
	BuildAt     string
}

func (n *N9EPlugin) Descript() string {
	return fmt.Sprintf("%s: %s", n.Name, n.Description)
}

type Notice struct {
	Event *models.AlertCurEvent `json:"event"`
	Tpls  map[string]string     `json:"tpls"`
}

type AlertCurEventAbbr struct {
	Id        int64  `json:"id" gorm:"primaryKey"`
	Cluster   string `json:"cluster"`
	GroupId   int64  `json:"group_id"`   // busi group id
	GroupName string `json:"group_name"` // busi group name
	Hash      string `json:"hash"`       // rule_id + vector_key
	RuleId    int64  `json:"rule_id"`
	RuleName  string `json:"rule_name"`
	Severity  int    `json:"severity"`
	PromQl    string `json:"prom_ql"`

	TargetIdent      string      `json:"target_ident"`
	TriggerTime      int64       `json:"trigger_time"`
	TriggerValue     string      `json:"trigger_value"`
	IsRecovered      bool        `json:"is_recovered" gorm:"-"`     // for notify.py
	NotifyUsersObj   []*UserAbbr `json:"notify_users_obj" gorm:"-"` // for notify.py
	FirstTriggerTime int64       `json:"first_trigger_time"`        // 连续告警的首次告警时间
}
type UserAbbr struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin" gorm:"-"`
}

type HostConfig struct {
	AlertHostConfig string
}

func (n *N9EPlugin) Notify(bs []byte) {
	NoticeObj := Notice{}
	err := json.Unmarshal(bs, &NoticeObj)
	if err != nil {
		logger.Debugf("E! ", err)
	}
	AlertCurEventAbbrObj := AlertCurEventAbbr{}
	AlertCurEventAbbrObj.Id = NoticeObj.Event.Id
	AlertCurEventAbbrObj.Cluster = NoticeObj.Event.Cluster
	AlertCurEventAbbrObj.GroupId = NoticeObj.Event.GroupId
	AlertCurEventAbbrObj.GroupName = NoticeObj.Event.GroupName
	AlertCurEventAbbrObj.Hash = NoticeObj.Event.Hash
	AlertCurEventAbbrObj.RuleId = NoticeObj.Event.RuleId
	AlertCurEventAbbrObj.RuleName = NoticeObj.Event.RuleName
	AlertCurEventAbbrObj.Severity = NoticeObj.Event.Severity
	AlertCurEventAbbrObj.PromQl = NoticeObj.Event.PromQl
	AlertCurEventAbbrObj.TargetIdent = NoticeObj.Event.TargetIdent
	AlertCurEventAbbrObj.TriggerTime = NoticeObj.Event.TriggerTime
	AlertCurEventAbbrObj.TriggerValue = NoticeObj.Event.TriggerValue
	AlertCurEventAbbrObj.IsRecovered = NoticeObj.Event.IsRecovered
	AlertCurEventAbbrObj.FirstTriggerTime = NoticeObj.Event.FirstTriggerTime
	logger.Errorf("Cluster = %v", NoticeObj.Event.Cluster)
	logger.Errorf("GroupId = %v", NoticeObj.Event.GroupId)
	logger.Errorf("GroupName = %v", NoticeObj.Event.GroupName)
	logger.Errorf("RuleName = %v", NoticeObj.Event.RuleName)
	logger.Errorf("Severity = %v", NoticeObj.Event.Severity)
	logger.Errorf("PromQl = %v", NoticeObj.Event.PromQl)
	logger.Errorf("TriggerValue = %v", NoticeObj.Event.TriggerValue)
	logger.Errorf("TargetIdent = %v", NoticeObj.Event.TargetIdent)
	logger.Errorf("TriggerTime = %v", NoticeObj.Event.TriggerTime)
	logger.Errorf("TriggerTime = %v", AlertCurEventAbbrObj)
	NotifyUsersObj := NoticeObj.Event.NotifyUsersObj
	var users []*UserAbbr
	for _, user := range NotifyUsersObj {

		UserAbbrObj := UserAbbr{}
		UserAbbrObj.Username = user.Username
		UserAbbrObj.Nickname = user.Nickname
		UserAbbrObj.Phone = user.Phone
		UserAbbrObj.Email = user.Email
		UserAbbrObj.Admin = user.Admin
		users = append(users, &UserAbbrObj)
		logger.Errorf("NotifyUsersObj = %v", user.Username)
		logger.Errorf("NotifyUsersObj = %v", user.Nickname)
		logger.Errorf("NotifyUsersObj = %v", user.Phone)
		logger.Errorf("NotifyUsersObj = %v", user.Email)
		logger.Errorf("NotifyUsersObj = %v", user.Admin)
	}
	AlertCurEventAbbrObj.NotifyUsersObj = users

	m := multiconfig.NewWithPath("conf/notify_config.toml") // supports TOML, JSON and YAML

	hostConfig := new(HostConfig)
	err = m.Load(hostConfig) // Check for error
	if err != nil {
		fmt.Println(err)
		return
	}
	m.MustLoad(hostConfig)
	fmt.Println(hostConfig.AlertHostConfig)

	bytesData, err := json.Marshal(AlertCurEventAbbrObj)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(bytesData))

	req, err := http.NewRequest("POST", hostConfig.AlertHostConfig, bytes.NewBuffer(bytesData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)

}

func (n *N9EPlugin) NotifyMaintainer(bs []byte) {
	fmt.Println("do something... begin")
	result := string(bs)
	fmt.Println(result)
	fmt.Println("do something... end")
}

// will be loaded for alertingCall , The first letter must be capitalized to be exported
var N9eCaller = N9EPlugin{
	Name:        "N9EPlugin",
	Description: "Notify by lib",
	BuildAt:     time.Now().Local().Format("2022/08/11 15:04:05"),
}
