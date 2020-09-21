# Cấu trúc API

## Cấu hình

Khi đề cập đến các API trong phần này, chúng tôi sẽ không kèm theo **BASE_URL**
và bạn mặc định hiểu là khi truy vấn sẽ gắn thêm chuỗi **BASE_URL** phía trước
đường dẫn cụ thể của API. Ví dụ `/v1/vnposts/ping` sẽ là
`https://vnpost-api-development.movecrop.com/v1/vnposts/ping`. Bạn cũng cần một
**TOKEN** hợp lệ để truy cập các API.

```bash
export BASE_URL=https://vnpost-api-development.movecrop.com
export TOKEN=MzRiOWY1Mj...
```

## Request & Authorization

Một lời gọi API tiêu biểu như sau:

```bash
curl $BASE_URL/v1/vnposts/ping \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}'
```

Tất cả request sử dụng giao thức **HTTPS**, phương thức **POST** và truyền giá
trị bằng body sử dụng định dạng dữ liệu `application/json`. Các header bắt buộc:

| Header | Nội dung | Mô tả |
| --- | --- | --- |
| Content-Type | application/json | |
| Authorization | Bearer MzRiOWY1Mj... | $TOKEN được cung cấp |

## HTTP Code

| Status | Mô tả |
| --- | --- |
| 200 | OK - Request đã được xử lý thành công.
| 400 | Bad Request -  Request không đúng cấu trúc hoặc thiếu dữ liệu yêu cầu. Hãy kiểm tra lại request.
| 401 | Unauthorized - Token của bạn không đúng hoặc đã hết hiệu lực.
| 403 | Forbidden - Bạn không có quyền truy cập đối tượng này. 
| 404 | Not Found - API Path không tồn tại, hoặc đối tượng được request không tồn tại.
| 405 | Method Not Allowed - Luôn sử dụng **POST** trong tất cả API của chúng tôi.
| 500 | Internal Server Error - Hệ thống đang gặp sự cố. Hãy thử lại sau hoặc liên hệ với chúng tôi.
