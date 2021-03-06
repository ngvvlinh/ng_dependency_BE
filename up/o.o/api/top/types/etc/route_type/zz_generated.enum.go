// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package route_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumCustomRegionRouteTypeName = map[int]string{
	1: "noi_vung",
	2: "lien_vung",
}

var enumCustomRegionRouteTypeValue = map[string]int{
	"noi_vung":  1,
	"lien_vung": 2,
}

func ParseCustomRegionRouteType(s string) (CustomRegionRouteType, bool) {
	val, ok := enumCustomRegionRouteTypeValue[s]
	return CustomRegionRouteType(val), ok
}

func ParseCustomRegionRouteTypeWithDefault(s string, d CustomRegionRouteType) CustomRegionRouteType {
	val, ok := enumCustomRegionRouteTypeValue[s]
	if !ok {
		return d
	}
	return CustomRegionRouteType(val)
}

func (e CustomRegionRouteType) Apply(d CustomRegionRouteType) CustomRegionRouteType {
	if e == 0 {
		return d
	}
	return e
}

func (e CustomRegionRouteType) Enum() int {
	return int(e)
}

func (e CustomRegionRouteType) Name() string {
	return enumCustomRegionRouteTypeName[e.Enum()]
}

func (e CustomRegionRouteType) String() string {
	s, ok := enumCustomRegionRouteTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("CustomRegionRouteType(%v)", e.Enum())
}

func (e CustomRegionRouteType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumCustomRegionRouteTypeName[e.Enum()] + "\""), nil
}

func (e *CustomRegionRouteType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumCustomRegionRouteTypeValue, data, "CustomRegionRouteType")
	if err != nil {
		return err
	}
	*e = CustomRegionRouteType(value)
	return nil
}

func (e CustomRegionRouteType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *CustomRegionRouteType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumCustomRegionRouteTypeValue, src, "CustomRegionRouteType")
	*e = (CustomRegionRouteType)(value)
	return err
}

func (e CustomRegionRouteType) Wrap() NullCustomRegionRouteType {
	return WrapCustomRegionRouteType(e)
}

func ParseCustomRegionRouteTypeWithNull(s dot.NullString, d CustomRegionRouteType) NullCustomRegionRouteType {
	if !s.Valid {
		return NullCustomRegionRouteType{}
	}
	val, ok := enumCustomRegionRouteTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return CustomRegionRouteType(val).Wrap()
}

func WrapCustomRegionRouteType(enum CustomRegionRouteType) NullCustomRegionRouteType {
	return NullCustomRegionRouteType{Enum: enum, Valid: true}
}

func (n NullCustomRegionRouteType) Apply(s CustomRegionRouteType) CustomRegionRouteType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullCustomRegionRouteType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullCustomRegionRouteType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullCustomRegionRouteType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullCustomRegionRouteType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumProvinceRouteTypeName = map[int]string{
	1: "noi_tinh",
	2: "lien_tinh",
}

var enumProvinceRouteTypeValue = map[string]int{
	"noi_tinh":  1,
	"lien_tinh": 2,
}

func ParseProvinceRouteType(s string) (ProvinceRouteType, bool) {
	val, ok := enumProvinceRouteTypeValue[s]
	return ProvinceRouteType(val), ok
}

func ParseProvinceRouteTypeWithDefault(s string, d ProvinceRouteType) ProvinceRouteType {
	val, ok := enumProvinceRouteTypeValue[s]
	if !ok {
		return d
	}
	return ProvinceRouteType(val)
}

func (e ProvinceRouteType) Apply(d ProvinceRouteType) ProvinceRouteType {
	if e == 0 {
		return d
	}
	return e
}

func (e ProvinceRouteType) Enum() int {
	return int(e)
}

func (e ProvinceRouteType) Name() string {
	return enumProvinceRouteTypeName[e.Enum()]
}

func (e ProvinceRouteType) String() string {
	s, ok := enumProvinceRouteTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ProvinceRouteType(%v)", e.Enum())
}

func (e ProvinceRouteType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumProvinceRouteTypeName[e.Enum()] + "\""), nil
}

func (e *ProvinceRouteType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumProvinceRouteTypeValue, data, "ProvinceRouteType")
	if err != nil {
		return err
	}
	*e = ProvinceRouteType(value)
	return nil
}

func (e ProvinceRouteType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ProvinceRouteType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumProvinceRouteTypeValue, src, "ProvinceRouteType")
	*e = (ProvinceRouteType)(value)
	return err
}

func (e ProvinceRouteType) Wrap() NullProvinceRouteType {
	return WrapProvinceRouteType(e)
}

func ParseProvinceRouteTypeWithNull(s dot.NullString, d ProvinceRouteType) NullProvinceRouteType {
	if !s.Valid {
		return NullProvinceRouteType{}
	}
	val, ok := enumProvinceRouteTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ProvinceRouteType(val).Wrap()
}

func WrapProvinceRouteType(enum ProvinceRouteType) NullProvinceRouteType {
	return NullProvinceRouteType{Enum: enum, Valid: true}
}

func (n NullProvinceRouteType) Apply(s ProvinceRouteType) ProvinceRouteType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullProvinceRouteType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullProvinceRouteType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullProvinceRouteType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullProvinceRouteType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumRegionRouteTypeName = map[int]string{
	1: "noi_mien",
	2: "lien_mien",
	3: "can_mien",
}

var enumRegionRouteTypeValue = map[string]int{
	"noi_mien":  1,
	"lien_mien": 2,
	"can_mien":  3,
}

func ParseRegionRouteType(s string) (RegionRouteType, bool) {
	val, ok := enumRegionRouteTypeValue[s]
	return RegionRouteType(val), ok
}

func ParseRegionRouteTypeWithDefault(s string, d RegionRouteType) RegionRouteType {
	val, ok := enumRegionRouteTypeValue[s]
	if !ok {
		return d
	}
	return RegionRouteType(val)
}

func (e RegionRouteType) Apply(d RegionRouteType) RegionRouteType {
	if e == 0 {
		return d
	}
	return e
}

func (e RegionRouteType) Enum() int {
	return int(e)
}

func (e RegionRouteType) Name() string {
	return enumRegionRouteTypeName[e.Enum()]
}

func (e RegionRouteType) String() string {
	s, ok := enumRegionRouteTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("RegionRouteType(%v)", e.Enum())
}

func (e RegionRouteType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumRegionRouteTypeName[e.Enum()] + "\""), nil
}

func (e *RegionRouteType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumRegionRouteTypeValue, data, "RegionRouteType")
	if err != nil {
		return err
	}
	*e = RegionRouteType(value)
	return nil
}

func (e RegionRouteType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *RegionRouteType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumRegionRouteTypeValue, src, "RegionRouteType")
	*e = (RegionRouteType)(value)
	return err
}

func (e RegionRouteType) Wrap() NullRegionRouteType {
	return WrapRegionRouteType(e)
}

func ParseRegionRouteTypeWithNull(s dot.NullString, d RegionRouteType) NullRegionRouteType {
	if !s.Valid {
		return NullRegionRouteType{}
	}
	val, ok := enumRegionRouteTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return RegionRouteType(val).Wrap()
}

func WrapRegionRouteType(enum RegionRouteType) NullRegionRouteType {
	return NullRegionRouteType{Enum: enum, Valid: true}
}

func (n NullRegionRouteType) Apply(s RegionRouteType) RegionRouteType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullRegionRouteType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullRegionRouteType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullRegionRouteType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullRegionRouteType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumUrbanTypeName = map[int]string{
	0: "unknown",
	1: "noi_thanh",
	2: "ngoai_thanh_1",
	3: "ngoai_thanh_2",
}

var enumUrbanTypeValue = map[string]int{
	"unknown":       0,
	"noi_thanh":     1,
	"ngoai_thanh_1": 2,
	"ngoai_thanh_2": 3,
}

func ParseUrbanType(s string) (UrbanType, bool) {
	val, ok := enumUrbanTypeValue[s]
	return UrbanType(val), ok
}

func ParseUrbanTypeWithDefault(s string, d UrbanType) UrbanType {
	val, ok := enumUrbanTypeValue[s]
	if !ok {
		return d
	}
	return UrbanType(val)
}

func (e UrbanType) Apply(d UrbanType) UrbanType {
	if e == 0 {
		return d
	}
	return e
}

func (e UrbanType) Enum() int {
	return int(e)
}

func (e UrbanType) Name() string {
	return enumUrbanTypeName[e.Enum()]
}

func (e UrbanType) String() string {
	s, ok := enumUrbanTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("UrbanType(%v)", e.Enum())
}

func (e UrbanType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumUrbanTypeName[e.Enum()] + "\""), nil
}

func (e *UrbanType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumUrbanTypeValue, data, "UrbanType")
	if err != nil {
		return err
	}
	*e = UrbanType(value)
	return nil
}

func (e UrbanType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *UrbanType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumUrbanTypeValue, src, "UrbanType")
	*e = (UrbanType)(value)
	return err
}

func (e UrbanType) Wrap() NullUrbanType {
	return WrapUrbanType(e)
}

func ParseUrbanTypeWithNull(s dot.NullString, d UrbanType) NullUrbanType {
	if !s.Valid {
		return NullUrbanType{}
	}
	val, ok := enumUrbanTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return UrbanType(val).Wrap()
}

func WrapUrbanType(enum UrbanType) NullUrbanType {
	return NullUrbanType{Enum: enum, Valid: true}
}

func (n NullUrbanType) Apply(s UrbanType) UrbanType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullUrbanType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullUrbanType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullUrbanType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullUrbanType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
