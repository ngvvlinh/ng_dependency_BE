package fbclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"gopkg.in/resty.v1"

	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

const (
	VI = "VI"
	EN = "EN"
)

type FacebookErrorService struct {
}

func NewFacebookErrorService() *FacebookErrorService {
	return &FacebookErrorService{}
}

func (s *FacebookErrorService) HandleErrorFacebookAPI(res *resty.Response, currentURL string) error {
	var bodyJson interface{}

	// handle "Sorry, this content isn't available right now"
	body := res.Body()
	if string(body) == "Sorry, this content isn't available right now" {
		return cm.Errorf(cm.FacebookError, nil, "").
			WithMeta("code", fmt.Sprintf("%v", 1)).
			WithMeta("sub_code", fmt.Sprintf("%d", 1357046)).
			WithMeta("message", "Received Invalid JSON reply.")
	}

	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return err
	}

	xBusinessUseCaseUsage := res.Header().Get(XBusinessUseCaseUsage)

	if facebookError := findError(bodyJson); facebookError != nil {
		return handleErrorFacebookAPI(facebookError, currentURL, xBusinessUseCaseUsage)
	}

	return nil
}

func handleErrorFacebookAPI(facebookError *model.FacebookError, currentURL, xBusinessUseCaseUsageHeader string) (_err error) {
	currentURL, err := censorTokens(currentURL)
	if err != nil {
		return err
	}

	defer func() {
		if _err != nil && xBusinessUseCaseUsageHeader != "" {
			_err = _err.(*xerrors.APIError).WithMetap(XBusinessUseCaseUsage, xBusinessUseCaseUsageHeader)
		}
	}()

	apiErrMsg := fmt.Sprintf("Facebook API error: %v", facebookError.Message.String)
	if facebookError.Code.Valid {
		if 200 <= facebookError.Code.Int && facebookError.Code.Int <= 299 {
			_err := cm.Errorf(cm.FacebookPermissionDenied, nil, apiErrMsg).
				WithMeta("code", fmt.Sprintf("%v", facebookError.Code.Int)).
				WithMeta("sub_code", fmt.Sprintf("%d", facebookError.ErrorSubcode.Int)).
				WithMeta("url", currentURL)
			if facebookError.Message.Valid {
				_err.WithMeta("msg", facebookError.Message.String)
			} else {
				_err.
					WithMeta("msg", mapErrorCodeMessage[MissingPermission][VI]).
					WithMeta("orig_msg", mapErrorCodeMessage[MissingPermission][EN])
			}
			return _err
		}
		_err := cm.Errorf(cm.FacebookError, nil, apiErrMsg).
			WithMeta("code", fmt.Sprintf("%v", facebookError.Code.Int)).
			WithMeta("sub_code", fmt.Sprintf("%d", facebookError.ErrorSubcode.Int)).
			WithMeta("url", currentURL)
		if facebookError.Message.Valid {
			_err.WithMeta("msg", facebookError.Message.String)
		} else {
			_err.
				WithMeta("msg", mapErrorCodeMessage[Code(facebookError.Code.Int)][VI]).
				WithMeta("orig_msg", mapErrorCodeMessage[Code(facebookError.Code.Int)][EN])
		}
	}
	if facebookError.ErrorSubcode.Valid {
		_err := cm.Errorf(cm.FacebookError, nil, apiErrMsg).
			WithMeta("code", fmt.Sprintf("%d", facebookError.Code.Int)).
			WithMeta("sub_code", fmt.Sprintf("%v", facebookError.ErrorSubcode.Int)).
			WithMeta("url", currentURL)
		if facebookError.Message.Valid {
			_err.WithMeta("msg", facebookError.Message.String)
		} else {
			_err.
				WithMeta("msg", mapErrorSubCodeMessage[SubCode(facebookError.ErrorSubcode.Int)][VI]).
				WithMeta("orig_msg", mapErrorSubCodeMessage[SubCode(facebookError.ErrorSubcode.Int)][EN])
		}
		return _err
	}
	if facebookError.Type.Valid {
		if facebookError.Type.String == "OAuthException" {
			_err := cm.Errorf(cm.FacebookError, nil, apiErrMsg).
				WithMeta("code", "OAuthException").
				WithMeta("sub_code", fmt.Sprintf("%d", facebookError.ErrorSubcode.Int)).
				WithMeta("url", currentURL)
			if facebookError.Message.Valid {
				_err.WithMeta("msg", facebookError.Message.String)
			} else {
				_err.
					WithMeta("msg", "Tr???ng th??i ????ng nh???p ho???c Access token ???? h???t h???n, xo?? b???, m???t kh??c c?? th??? kh??ng h???p l???.").
					WithMeta("orig_msg", "The login status or access token has expired, been revoked, or is otherwise invalid.")
			}
			return _err
		}
	}
	return nil
}

func findError(arg interface{}) *model.FacebookError {
	typ := reflect.TypeOf(arg)
	if typ.Kind() == reflect.Map {
		mapBody := arg.(map[string]interface{})
		if errorObject, ok := mapBody["error"]; ok {
			return convertFacebookError(errorObject.(map[string]interface{}))
		}

		for _, val := range mapBody {
			if reflect.TypeOf(val).Kind() == reflect.Map {
				facebookError := findError(val)
				if facebookError != nil {
					return facebookError
				}
			}
		}
	}

	return nil
}

func convertFacebookError(arg map[string]interface{}) *model.FacebookError {
	if len(arg) == 0 {
		return nil
	}
	var errorSubCode, code int
	var message, typ, errorUserTitle, fbTraceID string
	if arg["message"] != nil {
		message = arg["message"].(string)
	}
	if arg["type"] != nil {
		typ = arg["type"].(string)
	}
	if arg["error_user_title"] != nil {
		errorUserTitle = arg["error_user_title"].(string)
	}
	if arg["subcode"] != nil {
		errorSubCode = int(arg["subcode"].(float64))
	}
	if arg["error_subcode"] != nil {
		errorSubCode = int(arg["error_subcode"].(float64))
	}
	if arg["fbtrace_id"] != nil {
		fbTraceID = arg["fbtrace_id"].(string)
	}
	if arg["code"] != nil {
		code = int(arg["code"].(float64))
	}
	return &model.FacebookError{
		Message:        dot.String(message),
		Type:           dot.String(typ),
		Code:           dot.Int(code),
		ErrorSubcode:   dot.Int(errorSubCode),
		ErrorUserTitle: dot.String(errorUserTitle),
		FbtraceId:      dot.String(fbTraceID),
	}
}

func censorTokens(currentURL string) (string, error) {
	URL, err := url.Parse(currentURL)
	if err != nil {
		return "", err
	}

	queryParams, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return "", err
	}

	var keysRemoved []string

	for key := range queryParams {
		if strings.Contains(strings.ToLower(key), "token") {
			keysRemoved = append(keysRemoved, key)
		}
	}

	if len(keysRemoved) > 0 {
		for _, key := range keysRemoved {
			queryParams.Del(key)
		}
	}

	URL.RawQuery = queryParams.Encode()
	return URL.String(), nil
}

type Code int

func (c Code) String() string {
	return fmt.Sprintf("%d", c)
}

const (
	ApiUnknown                              = Code(1)
	ApiService                              = Code(2)
	ApiTooManyCalls                         = Code(4)
	ApiPermissionDenied                     = Code(10)
	ApiUserTooManyCalls                     = Code(17)
	InvalidParameter                        = Code(100)
	AccessTokenHasExpired                   = Code(190)
	MissingPermission                       = Code(200) // 200 - 299
	ApplicationLimitReached                 = Code(341)
	TemporarilyBlockedForPoliciesViolations = Code(368)
	DuplicatePost                           = Code(506)
	ErrorPostingLink                        = Code(1609005)
	HighMPS                                 = Code(1893007)
	PageLevelThrottling                     = Code(32)
	CustomLevelThrottling                   = Code(613)
	RateLimitCallWithPage                   = Code(80001)
	RateLimitCallWithMessenger              = Code(80006)
	UserCanNotReply                         = Code(10903)
	AlreadyRepliedTo                        = Code(10900)
	ReplyingTimeExpired                     = Code(10901)
)

var mapErrorCodeMessage = map[Code]map[string]string{
	ApiUnknown: {
		VI: "C?? l??? ??ang c?? m???t v???n ????? t???m th???i d???n ?????n hi???n t?????ng  cho ch???c n??ng n??y kh??ng ho???t ????ng. H??y ?????i v?? th??? l???i. N???u t??nh tr???ng v???n xu???t hi???n l???i, h??y ki???m tra l???i danh s??ch API.",
		EN: "Possibly a temporary issue due to downtime. Wait and retry the operation. If it occurs again, check that you are requesting an existing API.",
	},
	ApiService: {
		VI: "C?? l??? ??ang c?? m???t v???n ????? t???m th???i d???n ?????n hi???n t?????ng  cho ch???c n??ng n??y kh??ng ho???t ????ng. H??y ?????i v?? th??? l???i.",
		EN: "Temporary issue due to downtime. Wait and retry the operation.",
	},
	ApiTooManyCalls: {
		VI: "C?? m???t v???n ????? t???m th???i d???n ?????n hi???n t?????ng th???t c??? chai. H??y ?????i v?? th??? l???i, ho???c ki???m tra s??? l?????ng l???n g???i c???a API n??y.",
		EN: "Temporary issue due to throttling. Wait and retry the operation, or examine your API request volume.",
	},
	ApiUserTooManyCalls: {
		VI: "C?? m???t v???n ????? t???m th???i d???n ?????n hi???n t?????ng th???t c??? chai. H??y ?????i v?? th??? l???i, ho???c ki???m tra s??? l?????ng l???n g???i c???a API n??y.",
		EN: "Temporary issue due to throttling. Wait and retry the operation, or examine your API request volume.",
	},
	InvalidParameter: {
		VI: "Th??ng s??? kh??ng h???p l???.",
		EN: "Invalid parameter.",
	},
	ApiPermissionDenied: {
		VI: "Quy???n c?? th??? kh??ng ???????c c???p ph??t ho???c ???? xo?? b???.",
		EN: "Permission is either not granted or has been removed.",
	},
	AccessTokenHasExpired: {
		VI: "Access token ???? h???t h???n",
		EN: "Access token has expired.",
	},
	MissingPermission: {
		VI: "Quy???n c?? th??? kh??ng ???????c c???p ph??t ho???c ???? xo?? b???.",
		EN: "Permission is either not granted or has been removed.",
	},
	ApplicationLimitReached: {
		VI: "C?? m???t v???n ????? t???m th???i d???n ?????n hi???n t?????ng th???t c??? chai. H??y ?????i v?? th??? l???i, ho???c ki???m tra s??? l?????ng l???n g???i c???a API n??y.",
		EN: "Temporary issue due to downtime or throttling. Wait and retry the operation, or examine your API request volume.",
	},
	TemporarilyBlockedForPoliciesViolations: {
		VI: "Ch???c n??ng n??y t???m th???i b??? ch???n v?? vi ph???m ch??nh s??ch.",
		EN: "Temporarily blocked for policies violations.",
	},
	DuplicatePost: {
		VI: "C??c b??i vi???t gi???ng nhau kh??ng th??? ???????c c??ng b??? li??n ti???p. Thay ?????i n???i dung c???a b??i vi???t v?? th??? l???i.",
		EN: "Duplicate posts cannot be published consecutively. Change the content of the post and try again.",
	},
	ErrorPostingLink: {
		VI: "B??i vi???t c?? link kh??ng h???p l???. Ki???m tr??? link v?? th??? l???i.",
		EN: "There was a problem scraping data from the provided link. Check the URL and try again.",
	},
	PageLevelThrottling: {
		VI: "Page-level throttling",
		EN: "Page-level throttling",
	},
	CustomLevelThrottling: {
		VI: "Custom-level throttling",
		EN: "Custom-level throttling",
	},
}

type SubCode int

const (
	AppNotInstalled    = SubCode(458)
	UserCheckPointed   = SubCode(459)
	PasswordChanged    = SubCode(460)
	Expired            = SubCode(463)
	UnconfirmedUser    = SubCode(464)
	InvalidAccessToken = SubCode(467)
	InvalidSession     = SubCode(492)
	ObjectNotExist     = SubCode(33)
	MessageSentOutside = SubCode(2018278)
	PersonNotAvailable = SubCode(1545041)
)

var mapErrorSubCodeMessage = map[SubCode]map[string]string{
	AppNotInstalled: {
		VI: "T??i kho???n ch??a ???????c ????ng nh???p v??o app. H??y x??c th???c l???i t??i kho???n.",
		EN: "The User has not logged into your app. Reauthenticate the User.",
	},
	UserCheckPointed: {
		VI: "T??i kho???n c???n ????ng nh???p t???i ?????a ch??? https://www.facebook.com ho???c https://m.facebook.com ????? x??c ?????nh v???n ?????.",
		EN: "The User needs to log in at https://www.facebook.com or https://m.facebook.com to correct an issue.",
	},
	PasswordChanged: {
		VI: "Tr??n phi??n b???n IOS 6 ho???c l???n h??n, n???u m???t ng?????i ??ang ????ng nh???p s??? d???ng flow c???a h??? th???ng, ??i???u h?????ng h??? ?????n ph???n thi???t l???p c???a h??? th???ng Facebook ????? c???p nh???t m???t kh???u c???a h???. M???t kh??c, h??? ph???i ????ng nh???p l???i v??o app.",
		EN: "On iOS 6 and above, if the person logged in using the OS-integrated flow, direct them to Facebook OS settings on the device to update their password. Otherwise, they must log in to the app again.",
	},
	Expired: {
		VI: "Tr???ng th??i ????ng nh???p ho???c Access token ???? h???t h???n, xo?? b???, m???t kh??c c?? th??? kh??ng h???p l???.",
		EN: "Login status or access token has expired, been revoked, or is otherwise invalid.",
	},
	UnconfirmedUser: {
		VI: "T??i kho???n c???n ????ng nh???p t???i ?????a ch??? https://www.facebook.com ho???c https://m.facebook.com ????? x??c ?????nh v???n ?????.",
		EN: "The User needs to log in at https://www.facebook.com or https://m.facebook.com to correct an issue.",
	},
	InvalidAccessToken: {
		VI: "Access token ???? h???t h???n, xo?? b???, m???t kh??c c?? th??? kh??ng h???p l???.",
		EN: "Access token has expired, been revoked, or is otherwise invalid.",
	},
	InvalidSession: {
		VI: "T??i kho???n li??n quan ?????n Page access token kh??ng c?? vai tr?? ph?? h???p tr??n page n??y.",
		EN: "User associated with the Page access token does not have an appropriate role on the Page.",
	},
}
