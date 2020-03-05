package login

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

func init() {
	bus.AddHandler("login", LoginUser)
}

type LoginUserQuery struct {
	UserID       dot.ID
	PhoneOrEmail string
	Password     string

	Result struct {
		User *identitymodel.User
	}
}

const MsgLoginNotFound = `Người dùng chưa đăng ký. Vui lòng kiểm tra lại thông tin đăng nhập (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ %v.`
const MsgLoginUnauthenticated = `Mật khẩu không đúng. Vui lòng kiểm tra lại thông tin đăng nhập (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ %v.`

// Login flow
// 1. Get user
// 2. Check if the user activated
// 2a. If the user is activated, verify password and login
// 2b. If the user is not activated
//   - verify the generated password
//   - activate user
//   - then login
func LoginUser(ctx context.Context, query *LoginUserQuery) error {
	userQuery := &identitymodelx.GetUserByLoginQuery{
		UserID:       query.UserID,
		PhoneOrEmail: query.PhoneOrEmail,
	}
	if err := bus.Dispatch(ctx, userQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return cm.Errorf(cm.NotFound, nil, MsgLoginNotFound, wl.X(ctx).CSEmail).
				Log("NotFound: user does not exist")
		}
		return err
	}

	user := userQuery.Result.User
	hashpwd := userQuery.Result.UserInternal.Hashpwd

	// If the user is activated, verify password and login
	if user.Status != status3.Z {
		if !VerifyPassword(query.Password, hashpwd) {
			return cm.Errorf(cm.Unauthenticated, nil, MsgLoginUnauthenticated, wl.X(ctx).CSEmail).MarkTrivial()
		}

		// The user must be activated
		if user.Status != status3.P {
			return cm.Errorf(cm.Unauthenticated, nil,
				"Tài khoản của bạn đã bị khóa. Nếu cần thêm thông tin, vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
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
