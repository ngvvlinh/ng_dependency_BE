package user_source

// +enum
type UserSource int

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
)
