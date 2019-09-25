package httpreq

import (
	"encoding/json"
	"fmt"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

func IsNullJsonRaw(data json.RawMessage) bool {
	return len(data) == 0 ||
		len(data) == 4 && string(data) == "null"
}

func HandleResponse(res *RestyResponse, result interface{}, errMsg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if IsNullJsonRaw(body) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định: null response. Chúng tôi đang liên hệ để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
			}
			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định: %v. Chúng tôi đang liên hệ để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		var errJSON xerrors.ErrorJSON
		if !IsNullJsonRaw(body) {
			if err = jsonx.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errJSON.Msg = errMsg
			errJSON.Meta = meta
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi: %v. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", errJSON.Error()).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định: Invalid status (%v). Chúng tôi đang liên hệ để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", status)
	}
}
