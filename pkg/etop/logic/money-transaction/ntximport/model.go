package ntximport

import "o.o/backend/pkg/common/imcsv"

var schema = imcsv.Schema{
	{
		Name:    "stt",
		Display: "STT",
		Norm:    "stt",
	},
	{
		Name:     "customer_code",
		Display:  "Mã KH",
		Norm:     "ma kh",
		Optional: true,
	},
	{
		Name:     "shop_id",
		Display:  "Shop ID",
		Norm:     "shop id",
		Optional: true,
	},
	{
		Name:     "customer_name",
		Display:  "Tên KH",
		Norm:     "ten KH",
		Optional: true,
	},
	{
		Name:     "branch",
		Display:  "Chi Nhánh",
		Norm:     "chi nhanh",
		Optional: true,
	},
	{
		Name:     "num_vd",
		Display:  "Số VĐ",
		Norm:     "so vd",
		Optional: true,
	},
	{
		Name:    "do",
		Display: "DO",
		Norm:    "do",
	},
	{
		Name:     "do_kh",
		Display:  "DO KH",
		Norm:     "do kh",
		Optional: true,
	},
	{
		Name:     "sent_date",
		Display:  "Ngày Gửi",
		Norm:     "ngay gui",
		Optional: true,
	},
	{
		Name:     "sender",
		Display:  "Người Gửi",
		Norm:     "nguoi gui",
		Optional: true,
	},
	{
		Name:     "sender_address",
		Display:  "Địa Chỉ Gừi",
		Norm:     "dia chi gui",
		Optional: true,
	},
	{
		Name:     "receiver",
		Display:  "Người Nhận",
		Norm:     "nguoi nhan",
		Optional: true,
	},
	{
		Name:     "receiver_address",
		Display:  "Địa Chỉ Nhận",
		Norm:     "dia chi nhan",
		Optional: true,
	},
	{
		Name:     "phone",
		Display:  "Số Điện Thoại",
		Norm:     "so dien thoai",
		Optional: true,
	},
	{
		Name:     "email",
		Display:  "Email",
		Norm:     "email",
		Optional: true,
	},
	{
		Name:     "from_province",
		Display:  "Từ Tỉnh",
		Norm:     "tu tinh",
		Optional: true,
	},
	{
		Name:     "from_district",
		Display:  "Từ quận/huyện",
		Norm:     "tu quan/huyen",
		Optional: true,
	},
	{
		Name:     "to_province",
		Display:  "Đến Tỉnh",
		Norm:     "den tinh",
		Optional: true,
	},
	{
		Name:     "to_district",
		Display:  "Đến quận/huyện",
		Norm:     "den quan/huyen",
		Optional: true,
	},
	{
		Name:     "signature",
		Display:  "Chữ ký người nhận hàng",
		Norm:     "chu ky nguoi nhan hang",
		Optional: true,
	},
	{
		Name:     "return_date",
		Display:  "Ngày trả hàng",
		Norm:     "ngay tra hang",
		Optional: true,
	},
	{
		Name:     "payment_method",
		Display:  "Hình Thức Th/Toán",
		Norm:     "hinh thuc th/toan",
		Optional: true,
	},
	{
		Name:     "service",
		Display:  "Dịch Vụ",
		Norm:     "dich vu",
		Optional: true,
	},
	{
		Name:     "pack_num",
		Display:  "Số Kiện",
		Norm:     "so kien",
		Optional: true,
	},
	{
		Name:     "weight",
		Display:  "Trọng Lượng",
		Norm:     "trong luong",
		Optional: true,
	},
	{
		Name:    "weight_exchange",
		Display: "TL Quy Đổi",
		Norm:    "tl quy doi",
	},
	{
		Name:    "money_control",
		Display: "Tiền Kiểm Soát",
		Norm:    "tien kiem soat",
	},
	{
		Name:     "insurance_fees",
		Display:  "Phí Bảo Hiểm",
		Norm:     "phi bao hiem",
		Optional: true,
	},
	{
		Name:    "total_charge",
		Display: "Tổng cước",
		Norm:    "tong cuoc",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	stt          int
	shippingCode int
	codAmount    int
	shippingFee  int
}

var (
	schemas = []imcsv.Schema{schema}
	idxes   = []indexes{
		initIndexes(schema),
	}
	nColumn = len(schema)
)

func initIndexes(schema imcsv.Schema) indexes {
	indexer := schema.Indexer()
	return indexes{
		stt:          indexer("stt"),
		shippingCode: indexer("do"),
		codAmount:    indexer("money_control"),
		shippingFee:  indexer("total_charge"),
	}
}
