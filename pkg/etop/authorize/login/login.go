package login

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

var ll = l.New()

func init() {
	bus.AddHandler("login", LoginUser)
}

type LoginUserQuery struct {
	UserID       int64
	PhoneOrEmail string
	Password     string

	Result struct {
		User *model.User
	}
}

const MsgLoginNotFound = `Người dùng chưa đăng ký. Vui lòng kiểm tra lại thông tin đăng nhập (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.`
const MsgLoginUnauthenticated = `Mật khẩu không đúng. Vui lòng kiểm tra lại thông tin đăng nhập (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.`

// Login flow
// 1. Get user
// 2. Check if the user activated
// 2a. If the user is activated, verify password and login
// 2b. If the user is not activated
//   - verify the generated password
//   - activate user
//   - then login
func LoginUser(ctx context.Context, query *LoginUserQuery) error {
	userQuery := &model.GetUserByLoginQuery{
		UserID:       query.UserID,
		PhoneOrEmail: query.PhoneOrEmail,
	}
	if err := bus.Dispatch(ctx, userQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return cm.Error(cm.NotFound, MsgLoginNotFound, nil).
				Log("NotFound: user does not exist")
		}
		return err
	}

	user := userQuery.Result.User
	hashpwd := userQuery.Result.UserInternal.Hashpwd

	// If the user is activated, verify password and login
	if user.Status != model.S3Zero {
		if !VerifyPassword(query.Password, hashpwd) {
			return cm.Error(cm.Unauthenticated, MsgLoginUnauthenticated, nil).MarkTrivial()
		}

		// The user must be activated
		if user.Status != model.S3Positive {
			return cm.Error(cm.Unauthenticated,
				"Tài khoản của bạn đã bị khóa. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", nil)
		}
		query.Result.User = user
		return nil
	}

	// If the user is not activated, verify the generated password
	if hashpwd == "" {
		return cm.Error(cm.NotFound, MsgLoginNotFound, nil).
			Log("NotFound: user is not activated, no password")
	}
	if !VerifyPassword(query.Password, hashpwd) {
		return cm.Error(cm.NotFound, MsgLoginNotFound, nil).
			Log("NotFound: user is not activated, generated password does not match")
	}

	return cm.Error(cm.RegisterRequired,
		"Đăng nhập thành công. Vui lòng bổ sung thông tin người dùng.", nil)
}

// SaltSize is salt size in bytes.
const SaltSize = 16

// EncodePassword ...
func EncodePassword(password string) string {
	return hexa(saltedHashPassword([]byte(password)))
}

func VerifyPassword(password, hashpwd string) bool {
	return isPasswordMatch(dehexa(hashpwd), []byte(password))
}

func saltedHashPassword(secret []byte) []byte {
	buf := make([]byte, SaltSize, SaltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		ll.Panic("Unable to read from rand.Reader", l.Error(err))
		panic(err)
	}

	h := sha1.New()
	_, err = h.Write(buf)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	_, err = h.Write(secret)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	return h.Sum(buf)
}

func isPasswordMatch(data, secret []byte) bool {
	if len(data) != SaltSize+sha1.Size {
		panic("wrong length of data")
	}

	h := sha1.New()
	_, err := h.Write(data[:SaltSize])
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	_, err = h.Write(secret)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	return bytes.Equal(h.Sum(nil), data[SaltSize:])
}

func hexa(data []byte) string {
	return hex.EncodeToString(data)
}

func dehexa(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}