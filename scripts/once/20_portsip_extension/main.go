package main

import (
	"flag"

	"gopkg.in/resty.v1"
	"o.o/backend/com/etelecom/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/20_portsip_extension/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll          = l.New()
	cfg         config.Config
	DBEtelecom  *cmsql.Database
	url         = "https://sip.etelecom.vn:8900/api/extensions/update"
	accessToken = "002f6a7db4474d7c91d0bb78c70b1764"
)

// migration for English Now (tenant_id = 1183790567609863933)
var hotlineIDExtensionForward = map[dot.ID]string{
	1183927935316565615: "446598389300334592", // Eng - ext ID number 222
	1183928009676675054: "446607294596255744", // Eng 1 - ext ID number 666
}

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	ll.Info("Config", l.Object("cfg", cfg))
	if DBEtelecom, err = cmsql.Connect(cfg.PostgresEtelecom); err != nil {
		ll.Fatal("Error while loading database")
	}

	exts, err := GetExtensions()
	if err != nil {
		ll.Fatal("Get extensions failed", l.Error(err))
	}

	updated := 0
	for _, ext := range exts {
		err = updateExtensionForwardRules(ext)
		if err != nil {
			ll.Fatal("Update failed", l.Error(err))
		}
		updated++
	}

	count := len(exts)
	ll.S.Infof("Done. updated: %v/%v", updated, count)
}

func updateExtensionForwardRules(ext *model.Extension) error {
	extForwardID, ok := hotlineIDExtensionForward[ext.HotlineID]
	if !ok || extForwardID == "" {
		ll.S.Errorf("Not found extForwardID", l.ID("ext ID", ext.ID), l.ID("hotline ID", ext.HotlineID))
		return cm.Errorf(cm.FailedPrecondition, nil, "Not found extForwardID")
	}

	body := map[string]interface{}{
		"id": ext.ExternalData.ID,
		"forward_rules": map[string]interface{}{
			"available": map[string]interface{}{
				"no_answer_timeval":    15,
				"no_answer_fwdto_type": "CONNECT_TO_ANOTHER_EXTENSION",
				"no_answer_fwdto_id":   extForwardID,
				"busy_fwdto_type":      "CONNECT_TO_ANOTHER_EXTENSION",
				"busy_fwdto_id":        extForwardID,
			},
			"offline": map[string]interface{}{
				"office_hours_fwdto_type":         "CONNECT_TO_ANOTHER_EXTENSION",
				"office_hours_fwdto_id":           extForwardID,
				"outside_office_hours_fwdto_type": "CONNECT_TO_ANOTHER_EXTENSION",
				"outside_office_hours_fwdto_id":   extForwardID,
			},
		},
	}
	req := resty.R().
		SetHeader("access_token", accessToken).
		SetBody(body)
	_, err := req.Post(url)
	if err != nil {
		ll.Error("error", l.Error(err), l.ID("extension ID", ext.ID), l.String("extension number", ext.ExtensionNumber))
		return err
	}
	ll.Info("Update ext success", l.ID("extension ID", ext.ID), l.String("extension number", ext.ExtensionNumber))
	return nil
}

func GetExtensions() (exts model.Extensions, err error) {
	err = DBEtelecom.From("extension").
		Where("tenant_id = ?", 1183790567609863933).
		OrderBy("created_at ASC").
		Find(&exts)
	return
}
