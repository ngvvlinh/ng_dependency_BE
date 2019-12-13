package relationship

import (
	"context"
	"html/template"
	"strings"
	"time"

	"etop.vn/api/top/types/etc/status3"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/common/l"
)

var ll = l.New()
var bot *telebot.Channel
var idempgroup = idemp.NewGroup()

func init() {
	bus.AddHandlers("logic",
		InviteUserToAccount,
		GetOrCreateUserStub,
		GetOrCreateInvitation,
	)
}

func Init(b *telebot.Channel) {
	bot = b
}

// - Get user by login query
// - The the user does not exist, create a stub
// - Add an invitation

func InviteUserToAccount(ctx context.Context, cmd *InviteUserToAccountCommand) error {
	_, err := inviteUserToAccount(ctx, cmd)
	return err
}

func inviteUserToAccount(ctx context.Context, cmd *InviteUserToAccountCommand) (*InviteUserToAccountCommand, error) {
	if cmd.AccountID == 0 {
		return cmd, cm.Error(cm.InvalidArgument, "Missing account id", nil)
	}
	if cmd.InviterInfo.InviterUserID == 0 {
		return cmd, cm.Error(cm.InvalidArgument, "Missing inviter id", nil)
	}
	if cmd.EmailOrPhone == "" {
		return cmd, cm.Error(cm.InvalidArgument, "Missing required phone or email", nil)
	}

	userCommand := &GetOrCreateUserStubCommand{
		InviterInfo: cmd.InviterInfo,
		UserInner: &model.UserInner{
			FullName:  cmd.Invitation.FullName,
			ShortName: cmd.Invitation.ShortName,
			Email:     "",
			Phone:     "",
		},
	}
	if validate.IsEmail(cmd.EmailOrPhone) {
		emailNorm, ok := validate.NormalizeEmail(cmd.EmailOrPhone)
		if !ok {
			return cmd, cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
		userCommand.UserInner.Email = emailNorm.String()
	} else {
		phoneNorm, ok := validate.NormalizePhone(cmd.EmailOrPhone)
		if !ok {
			return cmd, cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
		}
		userCommand.UserInner.Phone = phoneNorm.String()
	}
	if err := bus.Dispatch(ctx, userCommand); err != nil {
		return cmd, err
	}
	user := userCommand.Result.User

	invitationCmd := &GetOrCreateInvitationCommand{
		AccountID:   cmd.AccountID,
		UserID:      user.ID,
		InviterInfo: cmd.InviterInfo,
		Invitation:  cmd.Invitation,
	}
	if err := bus.Dispatch(ctx, invitationCmd); err != nil {
		return cmd, err
	}

	cmd.Result = invitationCmd.Result
	return cmd, nil
}

var smsInvitationTpl = template.Must(template.New("sms").Parse(`
Bạn vừa được mời vào {{.AccountType}} {{.AccountName}} bởi người dùng {{.InviterFullName}}. Sử dụng sdt {{.Phone}} và mã {{.GeneratedPassword}} để đăng nhập vào {{.Domain}}. Nếu làm mất tin nhắn này, bạn có thể sử dụng chức năng khôi phục mật khẩu để đăng nhập.`))

func GetOrCreateUserStub(ctx context.Context, cmd *GetOrCreateUserStubCommand) error {
	_, err := getOrCreateUserStub(ctx, cmd)
	return err
}

func getOrCreateUserStub(ctx context.Context, cmd *GetOrCreateUserStubCommand) (*GetOrCreateUserStubCommand, error) {
	loginQuery := &model.GetUserByLoginQuery{}
	if cmd.Phone != "" {
		loginQuery.PhoneOrEmail = cmd.Phone
	} else if cmd.Email != "" {
		loginQuery.PhoneOrEmail = cmd.Email
	} else {
		return cmd, cm.Error(cm.InvalidArgument, "Missing required phone or email", nil)
	}

	err := bus.Dispatch(ctx, loginQuery)
	if err == nil {
		cmd.Result.User = loginQuery.Result.User
		cmd.Result.UserInternal = loginQuery.Result.UserInternal
		return cmd, nil
	}
	if cm.ErrorCode(err) != cm.NotFound {
		return cmd, err
	}

	// We allow user account without email address (only phone number).
	// This will send a message and let user knows that she can login.
	var generatedPassword string
	if cmd.Phone != "" {
		// generate random password
		generatedPassword, err = gencode.Random8Digits()
		if err != nil {
			return cmd, err
		}
	}

	createUserCmd := &model.CreateUserCommand{
		UserInner: model.UserInner{
			FullName:  "",
			ShortName: "",
			Email:     cmd.Email,
			Phone:     cmd.Phone,
		},
		Password:       generatedPassword,
		Status:         status3.Z,
		AgreeTOS:       false,
		AgreeEmailInfo: false,
		IsTest:         false,
		IsStub:         true,
	}
	if err := bus.Dispatch(ctx, createUserCmd); err != nil {
		return cmd, err
	}
	newUserID := createUserCmd.Result.User.ID

	if cmd.Phone != "" && generatedPassword != "" {
		go func() (_err error) {
			// The old http context was cancelled.
			// We've fired a new goroutine for sending sms.
			// So we need a new context.
			ctx := context.Background()

			var content string
			phone := cmd.Phone
			defer func() {
				if _err == nil {
					ll.Info("Sent sms invitation", l.String("phone", phone))
					return
				}

				ll.Error("Can not send invitation sms", l.Error(_err))
				bot.MaySendMessagef("Can not send invitation sms: %v\n\n%v", _err, content)
			}()

			domain := "etop.vn (staging)"
			if cm.IsProd() {
				domain = "etop.vn"
			}
			smsData := map[string]interface{}{
				"InviterFullName":   cmd.InviterInfo.InviterFullName,
				"AccountType":       cmd.InviterInfo.InviterAccountType,
				"AccountName":       cmd.InviterInfo.InviterAccountName,
				"Phone":             cmd.Phone,
				"GeneratedPassword": generatedPassword,
				"Domain":            domain,
			}
			var b strings.Builder
			if err := smsInvitationTpl.Execute(&b, smsData); err != nil {
				return err
			}
			content = b.String()

			// send sms
			sendSmsCmd := &sms.SendSMSCommand{
				Phone:   phone,
				Content: content,
			}
			if err := bus.Dispatch(ctx, sendSmsCmd); err != nil {
				return err
			}

			updateCmd := &model.UpdateUserVerificationCommand{
				UserID:                  newUserID,
				PhoneVerificationSentAt: time.Now(),
			}
			return bus.Dispatch(ctx, updateCmd)
		}()
	}

	cmd.Result.User = createUserCmd.Result.User
	cmd.Result.UserInternal = createUserCmd.Result.UserInternal
	return cmd, nil
}

func GetOrCreateInvitation(ctx context.Context, cmd *GetOrCreateInvitationCommand) error {
	_, err := getOrCreateInvitation(ctx, cmd)
	return err
}

func getOrCreateInvitation(ctx context.Context, cmd *GetOrCreateInvitationCommand) (*GetOrCreateInvitationCommand, error) {
	if err := cmd.Invitation.Permission.Validate(); err != nil {
		return cmd, err
	}

	accUserQuery := &model.GetAccountUserQuery{
		AccountID: cmd.AccountID,
		UserID:    cmd.UserID,
	}
	err := bus.Dispatch(ctx, accUserQuery)
	if err == nil {
		cmd.Result.AccountUser = accUserQuery.Result
		return cmd, nil
	}
	if cm.ErrorCode(err) != cm.NotFound {
		return cmd, err
	}

	accUser := &model.AccountUser{
		AccountID:        cmd.AccountID,
		UserID:           cmd.UserID,
		Status:           status3.Z, // Not activated yet
		Permission:       cmd.Invitation.Permission,
		FullName:         cmd.Invitation.FullName,
		ShortName:        cmd.Invitation.ShortName,
		Position:         cmd.Invitation.Position,
		InvitationSentAt: time.Now(),
		InvitationSentBy: cmd.InviterInfo.InviterUserID,
	}

	accUserCmd := &model.CreateAccountUserCommand{
		AccountUser: accUser,
	}
	if err := bus.Dispatch(ctx, accUserCmd); err != nil {
		return cmd, err
	}
	cmd.Result.AccountUser = accUserCmd.Result
	return cmd, nil
}
