package fbclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/extservice/telebot"
)

const (
	VI = "VI"
	EN = "EN"
)

type FacebookErrorService struct {
	bot *telebot.Channel
}

func NewFacebookErrorService(_bot *telebot.Channel) *FacebookErrorService {
	return &FacebookErrorService{
		bot: _bot,
	}
}

func (s *FacebookErrorService) HandleErrorFacebookAPI(body []byte, currentURL string) error {
	var facebookError struct {
		Data *struct {
			Error *model.FacebookError `json:"error"`
		} `json:"data"`
		Error *model.FacebookError `json:"error"`
	}

	if err := json.Unmarshal(body, &facebookError); err != nil {
		return err
	}

	if facebookError.Error != nil {
		return handleErrorFacebookAPI(facebookError.Error, currentURL)
	}
	if facebookError.Data != nil && facebookError.Data.Error != nil {
		return handleErrorFacebookAPI(facebookError.Data.Error, currentURL)
	}

	return nil
}

func handleErrorFacebookAPI(facebookError *model.FacebookError, currentURL string) error {
	currentURL, err := censorTokens(currentURL)
	if err != nil {
		return err
	}
	if facebookError.Code.Valid {
		if 200 <= facebookError.Code.Int && facebookError.Code.Int <= 299 {
			return cm.Errorf(cm.FailedPrecondition, nil, "Facebook API error").
				WithMeta("code", fmt.Sprintf("%v", facebookError.Code.Int)).
				WithMeta("msg", mapErrorCodeMessage[MissingPermission][VI]).
				WithMeta("orig_msg", mapErrorCodeMessage[MissingPermission][EN]).
				WithMeta("url", currentURL)
		}
		return cm.Errorf(cm.FailedPrecondition, nil, "Facebook API error").
			WithMeta("code", fmt.Sprintf("%v", facebookError.Code.Int)).
			WithMeta("msg", mapErrorCodeMessage[Code(facebookError.Code.Int)][VI]).
			WithMeta("orig_msg", mapErrorCodeMessage[Code(facebookError.Code.Int)][EN]).
			WithMeta("url", currentURL)
	}
	if facebookError.ErrorSubcode.Valid {
		return cm.Errorf(cm.FailedPrecondition, nil, "Facebook API error").
			WithMeta("sub_code", fmt.Sprintf("%v", facebookError.ErrorSubcode.Int)).
			WithMeta("msg", mapErrorSubCodeMessage[SubCode(facebookError.ErrorSubcode.Int)][VI]).
			WithMeta("orig_msg", mapErrorSubCodeMessage[SubCode(facebookError.ErrorSubcode.Int)][EN]).
			WithMeta("url", currentURL)
	}
	if facebookError.Type.Valid {
		if facebookError.Type.String == "OAuthException" {
			return cm.Errorf(cm.FailedPrecondition, nil, "Facebook API error").
				WithMeta("code", "OAuthException").
				WithMeta("msg", "Trạng thái đăng nhập hoặc Access token đã hết hạn, xoá bỏ, mặt khác có thể không hợp lệ.").
				WithMeta("orig_msg", "The login status or access token has expired, been revoked, or is otherwise invalid.").
				WithMeta("url", currentURL)
		}
	}
	return nil
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

const (
	ApiUnknown                              = Code(1)
	ApiService                              = Code(2)
	ApiTooManyCalls                         = Code(4)
	ApiPermissionDenied                     = Code(10)
	ApiUserTooManyCalls                     = Code(17)
	AccessTokenHasExpired                   = Code(190)
	MissingPermission                       = Code(200) // 200 - 299
	ApplicationLimitReached                 = Code(341)
	TemporarilyBlockedForPoliciesViolations = Code(368)
	DuplicatePost                           = Code(506)
	ErrorPostingLink                        = Code(1609005)
	PageLevelThrottling                     = Code(32)
	CustomLevelThrottling                   = Code(613)
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
