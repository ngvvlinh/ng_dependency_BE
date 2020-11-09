// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package connection_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumConnectionMethodName = map[int]string{
	0: "unknown",
	1: "builtin",
	2: "direct",
}

var enumConnectionMethodValue = map[string]int{
	"unknown": 0,
	"builtin": 1,
	"topship": 1,
	"direct":  2,
}

func ParseConnectionMethod(s string) (ConnectionMethod, bool) {
	val, ok := enumConnectionMethodValue[s]
	return ConnectionMethod(val), ok
}

func ParseConnectionMethodWithDefault(s string, d ConnectionMethod) ConnectionMethod {
	val, ok := enumConnectionMethodValue[s]
	if !ok {
		return d
	}
	return ConnectionMethod(val)
}

func (e ConnectionMethod) Apply(d ConnectionMethod) ConnectionMethod {
	if e == 0 {
		return d
	}
	return e
}

func (e ConnectionMethod) Enum() int {
	return int(e)
}

func (e ConnectionMethod) Name() string {
	return enumConnectionMethodName[e.Enum()]
}

func (e ConnectionMethod) String() string {
	s, ok := enumConnectionMethodName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ConnectionMethod(%v)", e.Enum())
}

func (e ConnectionMethod) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumConnectionMethodName[e.Enum()] + "\""), nil
}

func (e *ConnectionMethod) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumConnectionMethodValue, data, "ConnectionMethod")
	if err != nil {
		return err
	}
	*e = ConnectionMethod(value)
	return nil
}

func (e ConnectionMethod) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ConnectionMethod) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumConnectionMethodValue, src, "ConnectionMethod")
	*e = (ConnectionMethod)(value)
	return err
}

func (e ConnectionMethod) Wrap() NullConnectionMethod {
	return WrapConnectionMethod(e)
}

func ParseConnectionMethodWithNull(s dot.NullString, d ConnectionMethod) NullConnectionMethod {
	if !s.Valid {
		return NullConnectionMethod{}
	}
	val, ok := enumConnectionMethodValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ConnectionMethod(val).Wrap()
}

func WrapConnectionMethod(enum ConnectionMethod) NullConnectionMethod {
	return NullConnectionMethod{Enum: enum, Valid: true}
}

func (n NullConnectionMethod) Apply(s ConnectionMethod) ConnectionMethod {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullConnectionMethod) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullConnectionMethod) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullConnectionMethod) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullConnectionMethod) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumConnectionProviderName = map[int]string{
	0: "unknown",
	1: "ghn",
	2: "ghtk",
	3: "vtpost",
	4: "partner",
	5: "ahamove",
	6: "ninjavan",
	7: "dhl",
	8: "suitecrm",
}

var enumConnectionProviderValue = map[string]int{
	"unknown":  0,
	"ghn":      1,
	"ghtk":     2,
	"vtpost":   3,
	"partner":  4,
	"ahamove":  5,
	"ninjavan": 6,
	"dhl":      7,
	"suitecrm": 8,
}

func ParseConnectionProvider(s string) (ConnectionProvider, bool) {
	val, ok := enumConnectionProviderValue[s]
	return ConnectionProvider(val), ok
}

func ParseConnectionProviderWithDefault(s string, d ConnectionProvider) ConnectionProvider {
	val, ok := enumConnectionProviderValue[s]
	if !ok {
		return d
	}
	return ConnectionProvider(val)
}

func (e ConnectionProvider) Apply(d ConnectionProvider) ConnectionProvider {
	if e == 0 {
		return d
	}
	return e
}

func (e ConnectionProvider) Enum() int {
	return int(e)
}

func (e ConnectionProvider) Name() string {
	return enumConnectionProviderName[e.Enum()]
}

func (e ConnectionProvider) String() string {
	s, ok := enumConnectionProviderName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ConnectionProvider(%v)", e.Enum())
}

func (e ConnectionProvider) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumConnectionProviderName[e.Enum()] + "\""), nil
}

func (e *ConnectionProvider) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumConnectionProviderValue, data, "ConnectionProvider")
	if err != nil {
		return err
	}
	*e = ConnectionProvider(value)
	return nil
}

func (e ConnectionProvider) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ConnectionProvider) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumConnectionProviderValue, src, "ConnectionProvider")
	*e = (ConnectionProvider)(value)
	return err
}

func (e ConnectionProvider) Wrap() NullConnectionProvider {
	return WrapConnectionProvider(e)
}

func ParseConnectionProviderWithNull(s dot.NullString, d ConnectionProvider) NullConnectionProvider {
	if !s.Valid {
		return NullConnectionProvider{}
	}
	val, ok := enumConnectionProviderValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ConnectionProvider(val).Wrap()
}

func WrapConnectionProvider(enum ConnectionProvider) NullConnectionProvider {
	return NullConnectionProvider{Enum: enum, Valid: true}
}

func (n NullConnectionProvider) Apply(s ConnectionProvider) ConnectionProvider {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullConnectionProvider) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullConnectionProvider) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullConnectionProvider) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullConnectionProvider) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumConnectionSubtypeName = map[int]string{
	0: "unknown",
	1: "shipment",
	2: "manual",
	3: "shipnow",
}

var enumConnectionSubtypeValue = map[string]int{
	"unknown":  0,
	"shipment": 1,
	"manual":   2,
	"shipnow":  3,
}

func ParseConnectionSubtype(s string) (ConnectionSubtype, bool) {
	val, ok := enumConnectionSubtypeValue[s]
	return ConnectionSubtype(val), ok
}

func ParseConnectionSubtypeWithDefault(s string, d ConnectionSubtype) ConnectionSubtype {
	val, ok := enumConnectionSubtypeValue[s]
	if !ok {
		return d
	}
	return ConnectionSubtype(val)
}

func (e ConnectionSubtype) Apply(d ConnectionSubtype) ConnectionSubtype {
	if e == 0 {
		return d
	}
	return e
}

func (e ConnectionSubtype) Enum() int {
	return int(e)
}

func (e ConnectionSubtype) Name() string {
	return enumConnectionSubtypeName[e.Enum()]
}

func (e ConnectionSubtype) String() string {
	s, ok := enumConnectionSubtypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ConnectionSubtype(%v)", e.Enum())
}

func (e ConnectionSubtype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumConnectionSubtypeName[e.Enum()] + "\""), nil
}

func (e *ConnectionSubtype) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumConnectionSubtypeValue, data, "ConnectionSubtype")
	if err != nil {
		return err
	}
	*e = ConnectionSubtype(value)
	return nil
}

func (e ConnectionSubtype) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ConnectionSubtype) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumConnectionSubtypeValue, src, "ConnectionSubtype")
	*e = (ConnectionSubtype)(value)
	return err
}

func (e ConnectionSubtype) Wrap() NullConnectionSubtype {
	return WrapConnectionSubtype(e)
}

func ParseConnectionSubtypeWithNull(s dot.NullString, d ConnectionSubtype) NullConnectionSubtype {
	if !s.Valid {
		return NullConnectionSubtype{}
	}
	val, ok := enumConnectionSubtypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ConnectionSubtype(val).Wrap()
}

func WrapConnectionSubtype(enum ConnectionSubtype) NullConnectionSubtype {
	return NullConnectionSubtype{Enum: enum, Valid: true}
}

func (n NullConnectionSubtype) Apply(s ConnectionSubtype) ConnectionSubtype {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullConnectionSubtype) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullConnectionSubtype) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullConnectionSubtype) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullConnectionSubtype) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumConnectionTypeName = map[int]string{
	0: "unknown",
	1: "shipping",
	2: "crm",
}

var enumConnectionTypeValue = map[string]int{
	"unknown":  0,
	"shipping": 1,
	"crm":      2,
}

func ParseConnectionType(s string) (ConnectionType, bool) {
	val, ok := enumConnectionTypeValue[s]
	return ConnectionType(val), ok
}

func ParseConnectionTypeWithDefault(s string, d ConnectionType) ConnectionType {
	val, ok := enumConnectionTypeValue[s]
	if !ok {
		return d
	}
	return ConnectionType(val)
}

func (e ConnectionType) Apply(d ConnectionType) ConnectionType {
	if e == 0 {
		return d
	}
	return e
}

func (e ConnectionType) Enum() int {
	return int(e)
}

func (e ConnectionType) Name() string {
	return enumConnectionTypeName[e.Enum()]
}

func (e ConnectionType) String() string {
	s, ok := enumConnectionTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ConnectionType(%v)", e.Enum())
}

func (e ConnectionType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumConnectionTypeName[e.Enum()] + "\""), nil
}

func (e *ConnectionType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumConnectionTypeValue, data, "ConnectionType")
	if err != nil {
		return err
	}
	*e = ConnectionType(value)
	return nil
}

func (e ConnectionType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ConnectionType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumConnectionTypeValue, src, "ConnectionType")
	*e = (ConnectionType)(value)
	return err
}

func (e ConnectionType) Wrap() NullConnectionType {
	return WrapConnectionType(e)
}

func ParseConnectionTypeWithNull(s dot.NullString, d ConnectionType) NullConnectionType {
	if !s.Valid {
		return NullConnectionType{}
	}
	val, ok := enumConnectionTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ConnectionType(val).Wrap()
}

func WrapConnectionType(enum ConnectionType) NullConnectionType {
	return NullConnectionType{Enum: enum, Valid: true}
}

func (n NullConnectionType) Apply(s ConnectionType) ConnectionType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullConnectionType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullConnectionType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullConnectionType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullConnectionType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
