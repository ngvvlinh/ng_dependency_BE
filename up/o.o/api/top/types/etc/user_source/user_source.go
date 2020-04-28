package user_source

// +enum
// +enum:zero=null
type UserSource int

type NullUserSource struct {
	Enum  UserSource
	Valid bool
}

const (
	// +enum=unknown
	Unknown UserSource = 0

	// +enum=psx
	PSX UserSource = 1

	// +enum=etop
	Etop UserSource = 2

	// +enum=topship
	Topship UserSource = 3

	// +enum=ts_app_android
	TsAppAndroid UserSource = 4

	// +enum=ts_app_ios
	TsAppIOS UserSource = 5

	// +enum=ts_app_web
	TsAppWeb UserSource = 6

	// +enum=partner
	Partner UserSource = 7

	// +enum=etop_app_ios
	EtopAppIos UserSource = 8

	// +enum=etop_app_android
	EtopAppAndroid UserSource = 9
)
