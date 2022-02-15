package telecom

import (
	"o.o/capi/dot"
)

type TimeFrame struct {
	TenantID    dot.ID
	Start       string
	End         string
	Destination Destination
}

type Destination struct {
	Default  string
	Saturday string
	Sunday   string
}

var (
	ruleIndexTenantID  = make(map[dot.ID][]*TimeFrame)
	defaultDestination = make(map[dot.ID]string)
)

func init() {
	defaultDestination[1192042066291868827] = "800"
	for _, tf := range TimeFrames {
		ruleIndexTenantID[tf.TenantID] = append(ruleIndexTenantID[tf.TenantID], tf)
	}
}

func GetTimeFramesByTenantID(tenantID dot.ID) []*TimeFrame {
	return ruleIndexTenantID[tenantID]
}

func GetDefaultDestination(tenantID dot.ID) string {
	return defaultDestination[tenantID]
}
