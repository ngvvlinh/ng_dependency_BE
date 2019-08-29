package email

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"

	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/bus"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

type SMTPConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Encrypt     string `yaml:"encrypt"` // tls, ssl
	FromAddress string `yaml:"from_address"`
}

func (c *SMTPConfig) SMTPServer() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *SMTPConfig) MustLoadEnv(prefix ...string) {
	p := "ET_SMTP"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_HOST":         &c.Host,
		p + "_PORT":         &c.Port,
		p + "_USERNAME":     &c.Username,
		p + "_PASSWORD":     &c.Password,
		p + "_ENCRYPT":      &c.Encrypt,
		p + "_FROM_ADDRESS": &c.FromAddress,
	}.MustLoad()
}

type Client struct {
	cfg SMTPConfig
}

// New ...
func New(cfg SMTPConfig) *Client {
	cfg.Encrypt = strings.ToLower(cfg.Encrypt)
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Register(bus bus.Bus) *Client {
	bus.AddHandlers(c.SendMail)
	return c
}

type SendEmailCommand struct {
	FromName    string
	ToAddresses []string
	Subject     string
	Content     string
}

func (c *Client) Ping() error {
	client, err := c.Dial()
	if err != nil {
		return err
	}

	defer func() { _ = client.Quit() }()

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	err = client.Auth(auth)
	if err != nil {
		ll.Error(err.Error())
		return err
	}
	return nil
}

func (c *Client) Dial() (*smtp.Client, error) {
	smtpServer := c.cfg.SMTPServer()
	encrypt := c.cfg.Encrypt
	if encrypt == "" {
		return smtp.Dial(smtpServer)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.cfg.Host,
	}

	switch encrypt {
	case "ssl":
		conn, err := tls.Dial("tcp", smtpServer, tlsConfig)
		if err != nil {
			return nil, err
		}

		return smtp.NewClient(conn, c.cfg.Host)

	case "tls":
		client, err := smtp.Dial(smtpServer)
		if err != nil {
			return nil, err
		}

		err = client.StartTLS(tlsConfig)
		return client, err

	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Unknown encryption: %v", encrypt)
	}
}

func (c *Client) SendMail(ctx context.Context, cmd *SendEmailCommand) error {
	if len(cmd.ToAddresses) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing email address", nil)
	}

	addrs := make([]string, len(cmd.ToAddresses))
	for i, address := range cmd.ToAddresses {
		addr, _, ok := validate.TrimTest(address)
		if cm.IsDevOrStag() && !ok {
			return cm.Errorf(cm.FailedPrecondition, nil, "Chỉ có thể gửi email đến địa chỉ test trên dev!")
		}
		addrs[i] = addr
	}

	err := c.sendMail(ctx, addrs, cmd)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "Không thể gửi email đến địa chỉ %v (%v). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", strings.Join(addrs, ", "), err)
	}
	return nil
}

func (c *Client) sendMail(ctx context.Context, addresses []string, cmd *SendEmailCommand) error {
	client, err := c.Dial()
	if err != nil {
		ll.Error(err.Error())
		return err
	}
	defer func() { _ = client.Quit() }()

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	err = client.Auth(auth)
	if err != nil {
		ll.Error(err.Error())
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "=?utf-8?B?" + base64.StdEncoding.EncodeToString([]byte(cmd.Subject)) + "?="

	var errs xerrors.Errors
	for _, email := range addresses {
		msg := []byte(fmt.Sprintf(
			"From: %s <%s> \r\nTo: %s\r\nSubject: %s\r\n%s\r\n\r\n%s\r\n",
			cmd.FromName, c.cfg.FromAddress, email, subject, mime, cmd.Content))

		err = client.Mail(c.cfg.FromAddress)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = client.Rcpt(email)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		d, err := client.Data()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, err := d.Write(msg); err != nil {
			errs = append(errs, err)
			continue
		}
		err = d.Close()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) > 0 {
		ll.Error("Can not send email", l.Any("errs", err))
	}
	return errs.Any()
}
