package kiotviet

import (
	"context"
	"fmt"
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
)

func init() {
	bus.AddHandler("kiotviet", RequestWebhooks)
	bus.AddHandler("kiotviet", RequestCreateWebhook)
	bus.AddHandler("kiotviet", RequestDeleteWebhook)
	bus.AddHandler("kiotviet", EnsureWebhooks)
}

func RequestWebhooks(ctx context.Context, cmd *RequestWebhooksCommand) error {
	resp, err := rclient.R().
		SetAuthToken(cmd.Connection.TokenStr).
		SetHeader("Retailer", cmd.RetailerID).
		SetResult(&WebhooksResponse{}).
		SetError(&ErrorResponse{}).
		Get(cmd.URL("webhooks"))
	_, err = wrapRequest(resp, err)
	if err != nil {
		return err
	}
	cmd.Result = resp.Result().(*WebhooksResponse)
	return nil
}

func RequestCreateWebhook(ctx context.Context, cmd *RequestCreateWebhookCommand) error {
	if cmd.Webhook.URL == "" {
		return cm.Error(cm.InvalidArgument, "Invalid webhook", nil)
	}

	resp, err := rclient.R().
		SetAuthToken(cmd.Connection.TokenStr).
		SetHeader("Retailer", cmd.RetailerID).
		SetError(&ErrorResponse{}).
		SetBody(&WebhooksRequest{Webhook: cmd.Webhook}).
		Post(cmd.URL("webhooks"))
	_, err = wrapRequest(resp, err)
	return err
}

func RequestDeleteWebhook(ctx context.Context, cmd *RequestDeleteWebhookCommand) error {
	if cmd.ID == "" {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}

	resp, err := rclient.R().
		SetAuthToken(cmd.Connection.TokenStr).
		SetHeader("Retailer", cmd.RetailerID).
		SetError(&ErrorResponse{}).
		Delete(cmd.URL("webhooks/" + cmd.ID))
	_, err = wrapRequest(resp, err)
	return err
}

func EnsureWebhooks(ctx context.Context, cmd *EnsureWebhooksCommand) error {
	if cfg.skipWebhook {
		return nil
	}
	if cfg.WebhookBaseURL == "" {
		ll.Panic("Missing WebhookBaseURL")
	}
	if cm.Env() == "" {
		ll.Panic("Expect env")
	}

	conn, supplierID := cmd.Connection, cmd.SuppierID
	getCmd := &RequestWebhooksCommand{Connection: conn}
	if err := bus.Dispatch(ctx, getCmd); err != nil {
		return err
	}
	webhooks := getCmd.Result.Data

	var found *Webhook
	for _, wh := range webhooks {
		typ := string(wh.Type)
		if typ == "stock.update" {
			found = wh
			break
		}
	}

	// Kiotviet has a stupid decision when requires ?noecho in the URL. wtf!?
	stockURL := fmt.Sprintf("%vstock.update/%v?noecho", cfg.WebhookBaseURL, cmd.SuppierID)
	if found != nil && found.IsActive && string(found.URL) == stockURL {
		return nil
	}

	ll.Warn("Wanna update webhook",
		l.Any("old", found), l.String("new", stockURL))

	if found != nil {
		remoteEnv := strings.SplitN(string(found.Description), ":", 2)[0]
		if remoteEnv == cm.EnvProd && cm.Env() != cm.EnvProd {
			ll.Warn("Skip updating webhook, because the remote one has higher priority",
				l.String("remoteEnv", remoteEnv), l.String("env", cm.Env()))
			return nil
		}

		deleteCmd := &RequestDeleteWebhookCommand{
			Connection: conn,

			ID: string(found.ID),
		}
		if err := bus.Dispatch(ctx, deleteCmd); err != nil {
			return err
		}
	}

	updateCmd := &RequestCreateWebhookCommand{
		Connection: conn,
		Webhook: WebhookRequestItem{
			Type:        "stock.update",
			URL:         stockURL,
			Description: cm.Env() + ":stock.update",
			IsActive:    true,
		},
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return err
	}

	ll.Info("Updated webhook", l.Int64("supplier", supplierID), l.String("url", stockURL))
	return nil
}
