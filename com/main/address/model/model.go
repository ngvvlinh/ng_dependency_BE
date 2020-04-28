package model

import (
	"strconv"
	"strings"
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Address struct {
	ID        dot.ID `json:"id"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"` // Ward may be non-empty while WardCode is empty
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company     string       `json:"company"`
	Address1    string       `json:"address1"`
	Address2    string       `json:"address2"`
	Type        string       `json:"type"`
	AccountID   dot.ID       `json:"account_id"`
	Notes       *AddressNote `json:"notes"`
	CreatedAt   time.Time    `sq:"create" json:"-"`
	UpdatedAt   time.Time    `sq:"update" json:"-"`
	Coordinates *Coordinates `json:"coordinates"`

	Rid dot.ID `json:"rid"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Address) GetFullName() string {
	if m == nil {
		return ""
	}
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

func (m *Address) GetShortAddress() string {
	if m == nil {
		return ""
	}
	b := strings.Builder{}
	if m.Address1 != "" {
		b.WriteString(m.Address1)
		b.WriteByte('\n')
	}
	if m.Address2 != "" {
		b.WriteString(m.Address2)
		b.WriteByte('\n')
	}
	if m.Company != "" {
		b.WriteString(m.Company)
		b.WriteByte('\n')
	}
	s := b.String()
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

func (m *Address) GetPhone() string {
	if m == nil {
		return ""
	}
	return m.Phone
}

func (m *Address) GetProvince() string {
	if m == nil {
		return ""
	}
	return m.Province
}

func (m *Address) GetDistrict() string {
	if m == nil {
		return ""
	}
	return m.District
}

func (m *Address) GetWard() string {
	if m == nil {
		return ""
	}
	return m.Ward
}

// This function uses Ward (instead of WardCode) because WardCode may be empty
// while Ward retains raw name.
func (m *Address) GetFullAddress() string {
	b := strings.Builder{}
	if m.Address1 != "" {
		b.WriteString(m.Address1)
		b.WriteByte('\n')
	}
	if m.Address2 != "" {
		b.WriteString(m.Address2)
		b.WriteByte('\n')
	}
	if m.Company != "" {
		b.WriteString(m.Company)
		b.WriteByte('\n')
	}
	flag := false
	if m.Ward != "" {
		b.WriteString(m.Ward)
		flag = true
	}
	if m.District != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(m.District)
		flag = true
	}
	if m.Province != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(m.Province)
	}
	return b.String()
}

func (m *Address) UpdateAddress(phone string, fullname string) *Address {
	if phone != "" {
		m.Phone = phone
	}
	if fullname != "" {
		m.FullName = fullname
	}
	return m
}

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

func (m *AddressNote) GetFullNote() string {
	if m == nil {
		return ""
	}
	b := strings.Builder{}
	if m.Other != "" {
		if m.Other == "call" {
			b.WriteString("Gọi trước khi đến")
		} else if m.Other == "no-call" {
			b.WriteString("Không cần gọi trước, shop đã chuẩn bị sẵn")
		} else {
			b.WriteString(m.Other)
		}
		b.WriteString(". \n")
	}

	if m.Note != "" {
		b.WriteString(m.Note)
		if m.Note[len(m.Note)-1] != '.' {
			b.WriteByte('.')
		}
		b.WriteString(" \n")
	}
	if m.OpenTime != "" {
		b.WriteString("Giờ làm việc: ")
		b.WriteString(m.OpenTime)
		// Nếu làm việc cả buổi tối thì thêm dòng ghi chú này vào:
		// "(nếu lấy hàng không kịp vui lòng lấy buổi tối)"
		// format OpenTime: "08:00 - 21:00"
		text := strings.Split(m.OpenTime, "-")
		if len(text) > 1 {
			closedAt := strings.Split(text[1], ":")[0]
			if hour, err := strconv.Atoi(strings.TrimSpace(closedAt)); err == nil {
				if hour > 19 {
					b.WriteString(" (nếu lấy hàng không kịp vui lòng lấy buổi tối)")
				}
			}
		}
		b.WriteString(". \n")
	}
	if m.LunchBreak != "" {
		b.WriteString("giờ nghỉ trưa: ")
		b.WriteString(m.LunchBreak)
		b.WriteString(". \n")
	}
	s := b.String()
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}
