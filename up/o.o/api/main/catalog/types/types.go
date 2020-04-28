package types

import "o.o/common/jsonx"

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (m *Attribute) String() string { return jsonx.MustMarshalToString(m) }

type Attributes []*Attribute

func (attrs Attributes) Name() string {
	if len(attrs) == 0 {
		return ""
	}
	return attrs.ShortLabel()
}

func (attrs Attributes) Label() string {
	if len(attrs) == 0 {
		return ""
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ", "...)
		}
		b = append(b, attr.Name...)
		b = append(b, ": "...)
		b = append(b, attr.Value...)
	}
	return string(b)
}

func (attrs Attributes) LabelWithDefault(s string) string {
	if len(attrs) == 0 {
		return s
	}
	return attrs.Label()
}

func (attrs Attributes) ShortLabel() string {
	if len(attrs) == 0 {
		return ""
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, attr.Value...)
	}
	return string(b)
}

func (attrs Attributes) ShortLabelWithDefault(s string) string {
	if len(attrs) == 0 {
		return s
	}
	return attrs.ShortLabel()
}
