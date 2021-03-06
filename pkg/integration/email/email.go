package email

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
	"o.o/common/xerrors"
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
	c := &Client{
		cfg: cfg,
	}
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
		if cmenv.IsDevOrStag() && !ok {
			return cm.Errorf(cm.FailedPrecondition, nil, "Ch??? c?? th??? g???i email ?????n ?????a ch??? test tr??n dev!")
		}
		addrs[i] = addr
	}

	err := c.sendMail(ctx, addrs, cmd)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "Kh??ng th??? g???i email ?????n ?????a ch??? %v (%v). N???u c???n th??m th??ng tin, vui l??ng li??n h??? %v.", strings.Join(addrs, ", "), err, wl.X(ctx).CSEmail)
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
