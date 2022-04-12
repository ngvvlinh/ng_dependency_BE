package portsip_pbx

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"o.o/common/jsonx"
	"strconv"
	"time"
)

type SessionsResponse struct {
	Sesssion []*CallLog `json:"sesssions"`
	ScrollId string     `json:"scroll_id"`
}

func (m *SessionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CallLog struct {
	AnsweredTime      time.Time     `json:"answered_time"`
	CallID            string        `json:"call_id"`
	CallStatus        string        `json:"call_status"`
	CallTargets       []*CallTarget `json:"call_targets"`
	Callee            string        `json:"callee"`
	CalleeDomain      string        `json:"callee_domain"`
	Caller            string        `json:"caller"`
	CallerDisplayName string        `json:"caller_display_name"`
	CallerDomain      string        `json:"caller_domain"`
	DidCid            string        `json:"did_cid"`
	Direction         string        `json:"direction"`
	EndedReason       string        `json:"ended_reason"`
	EndedTime         time.Time     `json:"ended_time"`
	FailCode          int           `json:"fail_code"`
	FinalDest         string        `json:"final_dest"`
	OutboundCallerID  string        `json:"outbound_caller_id"`
	RecordingFileURL  string        `json:"recording_file_url"`
	RingDuration      int           `json:"ring_duration"`
	RingTime          time.Time     `json:"ring_time"`
	SessionID         string        `json:"session_id"`
	StartTime         time.Time     `json:"start_time"`
	TalkDuration      int           `json:"talk_duration"`
	TenantID          string        `json:"tenant_id"`
	TenantName        string        `json:"tenant_name"`
}

type CallTarget struct {
	AddTime      string `json:"add_time"`
	AnsweredTime string `json:"answered_time"`
	EndReason    string `json:"end_reason"`
	EndedTime    string `json:"ended_time"`
	FailCode     int    `json:"fail_code"`
	RingDuration int    `json:"ring_duration"`
	RingTime     string `json:"ring_time"`
	Status       string `json:"status"`
	TalkDuration int    `json:"talk_duration"`
	TargetNumber string `json:"target_number"`
	TrunkName    string `json:"trunk_name"`
}

func Convert_core_SearchHit_To_api_PortSipCalllog(in *elastic.SearchHit) *CallLog {
	if in == nil {
		return nil
	}

	source, err := in.Source.MarshalJSON()
	if err != nil {
		return nil
	}

	out := &CallLog{}
	if err := json.Unmarshal(source, &out); err != nil {
		return nil
	}

	return out
}

func Convert_core_SearchHits_To_api_PortSipCallLogs(in []*elastic.SearchHit) []*CallLog {
	if in == nil {
		return nil
	}
	outs := make([]*CallLog, len(in))
	for i, v := range in {
		outs[i] = Convert_core_SearchHit_To_api_PortSipCalllog(v)
	}

	return outs
}

func ConvertStringToInt64(s string) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
