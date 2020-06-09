package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
)

func TestNormalize(t *testing.T) {
	t.Run("Email", func(t *testing.T) {
		tests := []struct {
			name  string
			email string
			want  string
			ok    bool
		}{
			{
				"Empty (Invalid)",
				"", "", false,
			},
			{
				"Valid",
				"hello@example.com", "hello@example.com", true,
			},
			{
				"Trim space",
				"  hello.world@example.com \t ", "hello.world@example.com", true,
			},
			{
				"Inner space (invalid)",
				"hello.world @example.com", "", false,
			},
			{
				"Trim dot in gmail",
				"hello.world@gmail.com", "helloworld@gmail.com", true,
			},
			{
				"Trim space and dot",
				"uyenphuong.pollz@gmail.com ", "uyenphuongpollz@gmail.com", true,
			},
			{
				"Comment is not supported",
				"abc(hello)@gmail.com", "", false,
			},
			{
				"Case insensitive",
				"Hello@exAmple.com", "hello@example.com", true,
			},
			{
				"Test email",
				"test@example.com-test", "test@example.com-test", true,
			},
			{
				"test rule email 1",
				".Mothaibabon1234@gmail.com", "", false,
			},
			{
				"test rule email 2",
				"Mothaibabon1234.@gmail.com", "", false,
			},
			{
				"test rule email 3",
				"Mothaibabon..1234@gmail.com", "", false,
			},
			{
				"test rule email 4",
				"Mothaibabon12!34@gmail.com", "", false,
			},
			{
				"test rule email 5",
				"Mothaibabon1234@gmail.com", "mothaibabon1234@gmail.com", true,
			},
			{
				"test rule email 6",
				"Mothaibabon12!34@gmail.com", "", false,
			},
			{
				"test rule email 7",
				"Mothaibabon12.34@gmail.com", "mothaibabon1234@gmail.com", true,
			},
			{
				"test rule email 8",
				"Mothaibabon12^.34@gmail.com", "", false,
			},
			{
				"test rule email 9",
				"Mothaibabon12!34@example.com", "mothaibabon12!34@example.com", false,
			},
			{
				"test rule email 10",
				"Mothaibabon1234@example-extratest.com", "mothaibabon1234@example-extratest.com", true,
			},
			{
				"test rule email 11",
				"Mothaibabon12!34@example--extratest.com", "", false,
			},
			{
				"test rule email 12",
				"MothaibaBônTest@gmail.com", "", false,
			},
			{
				"test rule email 13",
				"Mothaibabon1234@example-extratest.com.", "", false,
			},
			{
				"test rule email 14",
				"Mothaibabon1234@example-extratest.com-", "", false,
			},
			{
				"test rule email 15",
				"Mothaibabon1234@example-extratest.plaza.com", "mothaibabon1234@example-extratest.plaza.com", true,
			},
			{
				"test rule email 16",
				"Mothaibabon1234@exam!ple-extratest.com", "", false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual, ok := NormalizeEmail(tt.email)
				require.Equal(t, tt.ok, ok)
				if ok {
					require.Equal(t, tt.want, string(actual))
				}
			})
		}
	})

	t.Run("popular email address mistake", func(t *testing.T) {
		tests := []struct {
			name  string
			email string
			ok    bool
		}{
			{
				"valid email address",
				"mothaibabon1234@gmail.com", true,
			},
			{
				"not the first character",
				"mothaibabon1234@ymail.com", true,
			},
			{
				"mail.com",
				"mothaibabon1234@mail.com", false,
			},
			{
				"swap 2 characters (invalid)",
				"mothaibabon1234@gamil.com", false,
			},
			{
				"swap more than 2 characters (valid)",
				"mothaibabon1234@gamil.cmo", true,
			},
			{
				"miss one character (invalid)",
				"mothaibabon1234@gmal.com", false,
			},
			{
				"miss more than 2 characters (valid)",
				"mothaibabon1234@gmal.co", true,
			},
			{
				"miss the last character (invalid)",
				"mothaibabon1234@gmail.co", false,
			},
			{
				"miss dot (invalid)",
				"mothaibabon1234@gmailcom", false,
			},
			{
				"different domain (valid)",
				"mothaibabon1234@yahoo.com", true,
			},
		}
		wlx := wl.Init(1, wl.EtopServer)
		var ctx = context.Background()
		ctx = wlx.WrapContext(ctx, drivers.ITopXID)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := PopularEmailAddressMistake(ctx, tt.email)
				if err != nil {
					require.Equal(t, tt.ok, false)
				} else {
					require.Equal(t, tt.ok, true)
				}
			})
		}
	})

	t.Run("Phone", func(t *testing.T) {
		tests := []struct {
			name  string
			phone string
			want  string
			ok    bool
		}{
			{
				"Empty (Invalid)",
				"", "", false,
			}, {
				"Valid",
				"0123456789", "0123456789", true,
			}, {
				"Trim dash",
				"0123-456-789 ", "0123456789", true,
			}, {
				"Too short (Invalid)",
				"123", "", false,
			}, {
				"Too short (Invalid)",
				"123", "", false,
			}, {
				"Trim tab character",
				"\t0932167545", "0932167545", true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actual, ok := NormalizePhone(tt.phone)

				require.Equal(t, tt.ok, ok)
				if ok {
					require.Equal(t, tt.want, string(actual))
				}
			})
		}
	})

	t.Run("Code", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			ok   bool
		}{
			{
				"Empty (Invalid)",
				"", false,
			}, {
				"Contain space (Invalid)",
				"A code", false,
			}, {
				"Valid",
				"0123456789_@#!?-[]", true,
			}, {
				"Too short (Invalid)",
				"12", false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ok := Code(tt.code)
				require.Equal(t, tt.ok, ok)
			})
		}
	})

	t.Run("Name", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
			ok    bool
		}{
			{
				"Empty (Invalid)",
				"", "", false,
			}, {
				"Trim space",
				" Sample  Name ", "Sample Name", true,
			}, {
				"Allow vietnamese",
				" Nguyễn Ngọc Minh An ", "Nguyễn Ngọc Minh An", true,
			}, {
				"Allow sign characters",
				` @#$%\/? (shop .)&`, `@#$%\/? (shop .)&`, true,
			}, {
				"Trim invalid characters",
				"One ≠ 1", "One 1", true,
			}, {
				"Too short after trim (Invalid)",
				" ≠≠≠≠≠≠A≠≠≠≠≠≠ ", "", false,
			}, {
				"Normalize vietnamese",
				"Nguy\u1ec5n Ng\u1ecdc", "Nguyễn Ngọc", true,
			}, {
				"Normalize vietnamese",
				"Nguy\u0065\u0302\u0303n Ng\u006f\u0323c", "Nguyễn Ngọc", true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, ok := NormalizeName(tt.input)

				require.Equal(t, tt.ok, ok)
				if ok {
					require.Equal(t, tt.want, got)
				}
			})
		}
	})

	t.Run("ExternalID", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{
				"Empty",
				"", "",
			},
			{
				"Normal (ok)",
				"ABC#1", "ABC#1",
			},
			{
				"Normal (error)",
				"ABC#1@", "",
			},
			{
				"With ~",
				"~ID Tuỳ \"ý' #123", "IDTuyy#123",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := NormalizeExternalCode(tt.input)
				require.Equal(t, tt.want, got)
			})
		}
	})
}

func TestParsePhoneInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 string
		want2 string
		ok    bool
	}{
		{
			"Empty (Invalid)",
			"", "", "", false,
		}, {
			"Valid",
			"0123456789", "0123456789", "", true,
		}, {
			"Trim dash",
			"0123-456-789 ", "0123456789", "", true,
		}, {
			"Too short (Invalid)",
			"12", "012", "", false,
		}, {
			"Ext 1",
			"011 234 5678 (ext: 16)",
			"0112345678", "", true,
		}, {
			"Ext 2",
			"11 234 5678 ext 16",
			"0112345678", "", true,
		}, {
			"Ext 3",
			"11 234 5678 (16)",
			"0112345678", "", true,
		}, {
			"2 numbers",
			"11 234 5678 - 011 234 5679",
			"0112345678", "0112345679", true,
		}, {
			"2 numbers with space",
			"112345678 0112345679",
			"0112345678", "0112345679", true,
		}, {
			"2 numbers with nothing",
			"1123456780112345679,",
			"", "", false,
		}, {
			"2 numbers with ext and hyphen",
			"11 234 5678 ext 16 - 011 234 5679 ext 2",
			"0112345678", "0112345679", true,
		}, {
			"2 numbers with ext and comma",
			"11-234-5678 ext 16, 011-234-5679",
			"0112345678", "0112345679", true,
		}, {
			"Single number with hyphen",
			"1900-6035",
			"19006035", "", true,
		}, {
			"2 numbers with hyphen",
			"123456789-123456790",
			"0123456789", "0123456790", true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, ok := ParsePhoneInput(tt.input)
			require.Equal(t, tt.ok, ok, "Parsed: %v %v", got1, got2)
			if ok {
				require.Equal(t, tt.want1, got1)
				require.Equal(t, tt.want2, got2)
			}
		})
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
		ok     bool
		test   bool
		base   string
		suffix string
	}{
		{
			"",
			"0848484841",
			"0848484841", true, false,
			"0848484841", "",
		},
		{
			"+84",
			"+84123456789",
			"0123456789", true, false,
			"+84123456789", "",
		},
		{
			"84",
			"84123456789",
			"0123456789", true, false,
			"84123456789", "",
		},
		{
			"+84 and 13 character (invalid)",
			"+84123456789000",
			"0123456789000", false, false,
			"+84123456789000", "",
		},
		{
			"",
			"0848484841",
			"0848484841", true, false,
			"0848484841", "",
		},
		{
			"Test phone number with -test",
			"0123456789-test",
			"0123456789-test", true, true,
			"0123456789", "-test",
		},
		{
			"Test phone number with -foo-test",
			"0123456789-foo-test",
			"0123456789-foo-test", true, true,
			"0123456789", "-foo-test",
		},
		{
			"Test phone number with -test but invalid (too short)",
			"012-test",
			"", false, true,
			"012", "-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := NormalizePhone(tt.input)
			require.Equal(t, tt.ok, ok, "Parsed: %v", actual)
			require.Equal(t, tt.expect, string(actual))

			a, b, ok := TrimTest(tt.input)
			require.Equal(t, tt.test, ok)
			require.Equal(t, tt.base, a)
			require.Equal(t, tt.suffix, b)
		})
	}
}

func TestIsNumberic(t *testing.T) {
	assert.Equal(t, true, isNumberic('0'))
	assert.Equal(t, true, isNumberic('9'))
	assert.Equal(t, false, isNumberic(' '))
}

func TestParseSinglePhoneNumber(t *testing.T) {
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("0123456789"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("84123456789"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("+84123456789"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("123456789"))
	// assert.Equal(t, "19001234", parseSinglePhoneNumber("1900-1234"))
	// assert.Equal(t, "1900", parseSinglePhoneNumber("1900"))
	// assert.Equal(t, "", parseSinglePhoneNumber("190"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("(84)123456789"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("123456789 ext 16"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("123-456-789 (ext: 16)"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("(84)123456789(16)"))
	// assert.Equal(t, "0123456789", parseSinglePhoneNumber("000123456789(16)"))
	assert.Equal(t, "", parseSinglePhoneNumber("0000000"))
}

func TestNormalizeSearch(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string // NormalizeSearchUnaccent
		wantu string // NormalizeSearchUnderscore
		wantq string // NormalizeSearch, quoted
		wanta string // NormalizeSearchAnd
		wanto string // NormalizeSearchOr
		wants string // NormalizedSearchToTsVector, sorted and removed duplicate
	}{
		{
			name: "Empty",
		},
		{
			"Lowercase",
			"This is An example",
			"this is an example",
			"this_is_an_example",
			"this is an example",
			"this & is & an & example",
			"this | is | an | example",
			"'an' 'example' 'is' 'this'",
		},
		{
			"Single special char",
			"%", "%", "", "'%'", "'%'", "'%'", "'%'",
		},
		{
			"Tripple special chars",
			"%%%", "%", "", "'%'", "'%'", "'%'", "'%'",
		},
		{
			"Consecutive special chars",
			"#@@#",
			"# @ #",
			"",
			"'#' '@' '#'",
			"'#' & '@' & '#'",
			"'#' | '@' | '#'",
			"'#' '@'",
		},
		{
			"Alphanumeric",
			"A1B2C",
			"a 1 b 2 c",
			"a1b2c",
			"a 1 b 2 c",
			"a & 1 & b & 2 & c",
			"a | 1 | b | 2 | c",
			"'1' '2' 'a' 'b' 'c'",
		},
		{
			"Remove vnese chars",
			"ÁO CHỮ T - KIỂU 2 - đen #6812-XL-100",
			"ao chu t - kieu 2 - den # 6812 - xl - 100",
			"ao_chu_t_kieu_2_den_6812_xl_100",
			"ao chu t '-' kieu 2 '-' den '#' 6812 '-' xl '-' 100",
			"ao & chu & t & '-' & kieu & 2 & '-' & den & '#' & 6812 & '-' & xl & '-' & 100",
			"ao | chu | t | '-' | kieu | 2 | '-' | den | '#' | 6812 | '-' | xl | '-' | 100",
			"'#' '-' '100' '2' '6812' 'ao' 'chu' 'den' 'kieu' 't' 'xl'",
		},
		{
			"Spaces, special chars",
			"% Giảm giá (%) ",
			"% giam gia ( % )",
			"giam_gia",
			"'%' giam gia '(' '%' ')'",
			"'%' & giam & gia & '(' & '%' & ')'",
			"'%' | giam | gia | '(' | '%' | ')'",
			"'%' '(' ')' 'gia' 'giam'",
		},
		{
			"Spaces, special and CAP chars",
			"  Akkj   tay phong@NANDA #6812  ",
			"akkj tay phong @ nanda # 6812",
			"akkj_tay_phong_nanda_6812",
			"akkj tay phong '@' nanda '#' 6812",
			"akkj & tay & phong & '@' & nanda & '#' & 6812",
			"akkj | tay | phong | '@' | nanda | '#' | 6812",
			"'#' '6812' '@' 'akkj' 'nanda' 'phong' 'tay'",
		},
		{
			"Single and double quote",
			`A'"\@1#`,
			`a @ 1 #`,
			"a_1",
			`a '@' 1 '#'`,
			`a & '@' & 1 & '#'`,
			`a | '@' | 1 | '#'`,
			`'#' '1' '@' 'a'`,
		},
		{
			"All special characters",
			`()[]{}<>/!@#$%^&*-_+=.,:;?|1a'"\`,
			`( ) [ ] { } < > / ! @ # $ % ^ & * - _ + = . , : ; ? | 1 a`,
			"1a",
			`'(' ')' '[' ']' '{' '}' '<' '>' '/' '!' '@' '#' '$' '%' '^' '&' '*' '-' '_' '+' '=' '.' ',' ':' ';' '?' '|' 1 a`,
			`'(' & ')' & '[' & ']' & '{' & '}' & '<' & '>' & '/' & '!' & '@' & '#' & '$' & '%' & '^' & '&' & '*' & '-' & '_' & '+' & '=' & '.' & ',' & ':' & ';' & '?' & '|' & 1 & a`,
			`'(' | ')' | '[' | ']' | '{' | '}' | '<' | '>' | '/' | '!' | '@' | '#' | '$' | '%' | '^' | '&' | '*' | '-' | '_' | '+' | '=' | '.' | ',' | ':' | ';' | '?' | '|' | 1 | a`,

			// The actual order from Postgres, we have to match this.
			`'!' '#' '$' '%' '&' '(' ')' '*' '+' ',' '-' '.' '/' '1' ':' ';' '<' '=' '>' '?' '@' '[' ']' '^' '_' 'a' '{' '|' '}'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := NormalizeUnaccent(tt.input)
			assert.Equal(t, tt.want, output)

			output = NormalizeSearch(tt.input)
			assert.Equal(t, tt.wantq, output)

			output = NormalizeUnderscore(tt.input)
			assert.Equal(t, tt.wantu, output)

			output = NormalizeSearchQueryAnd(tt.input)
			assert.Equal(t, tt.wanta, output)

			output = NormalizeSearchQueryOr(tt.input)
			assert.Equal(t, tt.wanto, output)

			output = NormalizedSearchToTsVector(NormalizeSearch(tt.input))
			assert.Equal(t, tt.wants, output)
		})
	}
}

func TestNormalizeSearchPhone(t *testing.T) {
	tests := []struct {
		input string
		want  string // NormalizeSerachPhone
	}{
		{
			"0945389709",
			"09 094 0945 09453 094538 0945389 09453897 094538970 0945389709",
		}, {
			"",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			n, output := normalizeSearchPhone(tt.input)
			assert.Equal(t, tt.want, output)
			assert.Equal(t, n, len(output))
		})
	}
}

func BenchmarkNormalizeSearchPhone(b *testing.B) {
	input := "0945389709"
	var _s string

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := NormalizeSearchPhone(input)
		_s = s
	}
	_ = _s
}
