package user_otp_action

// +enum
// +enum:zero=null
type UserOTPAction int

type NullUserOTPAction struct {
	Enum  UserOTPAction
	Valid bool
}

const (
	// +enum=unknown
	// +enum:RefName:Không xác định
	UserOTPActionUnknown UserOTPAction = 0

	// +enum=reset_password
	// +enum:RefName:Reset Password
	UserOTPActionResetPassword UserOTPAction = 1

	// +enum=register
	// +enum:RefName:Register
	UserOTPActionRegiter UserOTPAction = 2

	// +enum=verify_phone
	// +enum:RefName:Xác thực số điện thoại
	UserOTPActionVerifyPhone UserOTPAction = 3

	// +enum=update_phone_first_code
	// +enum:RefName:Cập nhật số điện thoại mã 1
	UserOTPActionUpdatePhoneFirstCode UserOTPAction = 4

	// +enum=update_phone_second_code
	// +enum:RefName:Cập nhật số điện thoại mã 2
	UserOTPActionUpdatePhoneSecondCode UserOTPAction = 5

	// +enum=update_email_first_code
	// +enum:RefName:Cập nhật email mã 1
	UserOTPActionUpdateEmailFirstCode UserOTPAction = 6

	// +enum=update_email_second_code
	// +enum:RefName:Cập nhật email mã 2
	UserOTPActionUpdateEmailSecondCode UserOTPAction = 7

	// +enum=verify_email_using_otp
	// +enum:RefName:Xác thực email sử dụng OTP
	UserOTPActionVerifyEmailUsingOTP UserOTPAction = 8

	// +enum=verify_email
	// +enum:RefName:Xác thức email
	UserOTPActionVerifyEmail UserOTPAction = 9 // using token

	// +enum=stoken_update_shop
	// +enum:RefName:Cập nhật tài khoản shop
	UserOTPActionSTokenUpdateShop UserOTPAction = 10
)
