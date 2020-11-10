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

	if err := json.Unmarshal(res.Body(), &bodyJson); err != nil {
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
					WithMeta("msg", "Trạng thái đăng nhập hoặc Access token đã hết hạn, xoá bỏ, mặt khác có thể không hợp lệ.").
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
		VI: "Có lẽ đang có một vấn đề tạm thời dẫn đến hiện tượng  cho chức năng này không hoạt đông. Hãy đợi và thử lại. Nếu tình trạng vẫn xuất hiện lại, hãy kiểm tra lại danh sách API.",
		EN: "Possibly a temporary issue due to downtime. Wait and retry the operation. If it occurs again, check that you are requesting an existing API.",
	},
	ApiService: {
		VI: "Có lẽ đang có một vấn đề tạm thời dẫn đến hiện tượng  cho chức năng này không hoạt đông. Hãy đợi và thử lại.",
		EN: "Temporary issue due to downtime. Wait and retry the operation.",
	},
	ApiTooManyCalls: {
		VI: "Có một vấn đề tạm thời dẫn đến hiện tượng thắt cổ chai. Hãy đợi và thử lại, hoặc kiểm tra số lượng lần gọi của API này.",
		EN: "Temporary issue due to throttling. Wait and retry the operation, or examine your API request volume.",
	},
	ApiUserTooManyCalls: {
		VI: "Có một vấn đề tạm thời dẫn đến hiện tượng thắt cổ chai. Hãy đợi và thử lại, hoặc kiểm tra số lượng lần gọi của API này.",
		EN: "Temporary issue due to throttling. Wait and retry the operation, or examine your API request volume.",
	},
	InvalidParameter: {
		VI: "Thông số không hợp lệ.",
		EN: "Invalid parameter.",
	},
	ApiPermissionDenied: {
		VI: "Quyền có thể không được cấp phát hoặc đã xoá bỏ.",
		EN: "Permission is either not granted or has been removed.",
	},
	AccessTokenHasExpired: {
		VI: "Access token đã hết hạn",
		EN: "Access token has expired.",
	},
	MissingPermission: {
		VI: "Quyền có thể không được cấp phát hoặc đã xoá bỏ.",
		EN: "Permission is either not granted or has been removed.",
	},
	ApplicationLimitReached: {
		VI: "Có một vấn đề tạm thời dẫn đến hiện tượng thắt cổ chai. Hãy đợi và thử lại, hoặc kiểm tra số lượng lần gọi của API này.",
		EN: "Temporary issue due to downtime or throttling. Wait and retry the operation, or examine your API request volume.",
	},
	TemporarilyBlockedForPoliciesViolations: {
		VI: "Chức năng này tạm thời bị chặn vì vi phạm chính sách.",
		EN: "Temporarily blocked for policies violations.",
	},
	DuplicatePost: {
		VI: "Các bài viết giống nhau không thể được công bố liên tiếp. Thay đổi nội dung của bài viết và thử lại.",
		EN: "Duplicate posts cannot be published consecutively. Change the content of the post and try again.",
	},
	ErrorPostingLink: {
		VI: "Bài viết có link không hợp lệ. Kiểm trả link và thử lại.",
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
		VI: "Tài khoản chưa được đăng nhập vào app. Hãy xác thực lại tài khoản.",
		EN: "The User has not logged into your app. Reauthenticate the User.",
	},
	UserCheckPointed: {
		VI: "Tài khoản cần đăng nhập tại địa chỉ https://www.facebook.com hoặc https://m.facebook.com để xác định vấn đề.",
		EN: "The User needs to log in at https://www.facebook.com or https://m.facebook.com to correct an issue.",
	},
	PasswordChanged: {
		VI: "Trên phiên bản IOS 6 hoặc lớn hơn, nếu một người đang đăng nhập sử dụng flow của hệ thống, điều hường họ đến phần thiết lập của hệ thống Facebook để cập nhật mật khẩu của họ. Mặt khác, họ phải đăng nhập lại vào app.",
		EN: "On iOS 6 and above, if the person logged in using the OS-integrated flow, direct them to Facebook OS settings on the device to update their password. Otherwise, they must log in to the app again.",
	},
	Expired: {
		VI: "Trạng thái đăng nhập hoặc Access token đã hết hạn, xoá bỏ, mặt khác có thể không hợp lệ.",
		EN: "Login status or access token has expired, been revoked, or is otherwise invalid.",
	},
	UnconfirmedUser: {
		VI: "Tài khoản cần đăng nhập tại địa chỉ https://www.facebook.com hoặc https://m.facebook.com để xác định vấn đề.",
		EN: "The User needs to log in at https://www.facebook.com or https://m.facebook.com to correct an issue.",
	},
	InvalidAccessToken: {
		VI: "Access token đã hết hạn, xoá bỏ, mặt khác có thể không hợp lệ.",
		EN: "Access token has expired, been revoked, or is otherwise invalid.",
	},
	InvalidSession: {
		VI: "Tài khoản liên quan đến Page access token không có vai trò phù hợp trên page này.",
		EN: "User associated with the Page access token does not have an appropriate role on the Page.",
	},
}
