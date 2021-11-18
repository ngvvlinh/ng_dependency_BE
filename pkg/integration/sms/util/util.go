package util

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

var (
	empty = ""
	space = " "
	dot   = "."
)

var removeWord = []string{
	"giam han muc",
	"CTKM",
	"Cong San",
	"Cong Hoa",
	"DCS",
	"Giftcode",
	"HOLCIM",
	"sexy",
	"bom",
	"ma tau",
	"bao ve to quoc",
	"cong dan",
	"quyen cong dan",
	"Vietnam",
	"vietnam",
	"video",
	"ysl",
	"than mat",
	"giet",
	"chem",
	"ren ri",
	"phong the",
	"oralsex",
	"cau be",
	"massage",
	"ham muon",
	"kich duc",
	"goi duc",
	"live stream",
	"sung suong",
	"cuc khoai",
	"duong vat",
	"dcm",
	"cd",
	"cmm",
	"wtf",
	"kich thich",
	"21+",
	"18+",
	"xxx",
	"show",
	"fuck",
	"dit me",
	"dek",
	"dech",
	"dm",
	"sex",
	"quan he",
	"cai tao",
	"cai cach",
	"mua/ban",
	"thanh uy",
	"mat tran",
	"dan vong",
	"tu tuong",
	"truong ban",
	"to chuc",
	"noi bo",
	"thuong vu",
	"usd",
	"lam phat",
	"lam quyen",
	"tham o",
	"tham nhung",
	"hanh phap",
	"chap phap",
	"bct",
	"uy vien",
	"dang uy",
	"bi thu",
	"ban chap hanh",
	"ban chi huy",
	"ban chi dao",
	"quan uy trung uong",
	"chu tich nuoc",
	"thu tuong",
	"tw",
	"trung uong doan",
	"ban bi thu",
	"trung uong dang",
	"bch",
	"truong ban chap hanh",
	"hoang trung hai",
	"vo van thuong",
	"nguyen van binh",
	"vuong dinh hue",
	"truong thi mai",
	"pham binh minh",
	"tran quoc vuong",
	"truong hoa binh",
	"tong thi phong",
	"pham minh chinh",
	"ngo xuan lich",
	"to lam",
	"dinh the huynh",
	"nguyen thien nhan",
	"tap hop",
	"lat do",
	"khoi nghia",
	"gom quan",
	"dai doan ket",
	"lien minh",
	"tap hop luc luong",
	"viet nam cong hoa",
	"dien bien hoa binh",
	"mat tran nhan dan",
	"to quoc",
	"formosa",
	"dam phan",
	"vinh bac bo",
	"doan ket",
	"rsf",
	"khong bien gioi",
	"bbc",
	"rfa",
	"a chau tu do",
	"tin do",
	"tu do ton giao",
	"quyen con nguoi",
	"hrw",
	"phan lap",
	"tam quyen",
	"tu tri",
	"ly khai",
	"lap phap",
	"kknls",
	"quoc gia tu tri",
	"lmtdvn",
	"thong nhat",
	"mtqg",
	"da nguyen",
	"dung co",
	"dau tranh",
	"thong luan",
	"dang nhan dan",
	"chinh phu",
	"vi dan",
	"cong hoa",
	"tu chu",
	"nhan quyen",
	"ngon luan",
	"tu do",
	"doc lap",
	"dan chu",
	"mua sim",
	"8406",
	"bau cu",
	"giai phong",
	"mua ban",
	"ban dat",
	"ban sim",
	"con duong viet nam",
	"quoc hoi",
	"dai hoi",
	"chinh tri",
	"ban nuoc",
	"vung an",
	"truong sa",
	"hoang sa",
	"dinh doc lap",
	"nha tho duc ba",
	"dinh la thang",
	"nguyen thi kim ngan",
	"nguyen xuan phuc",
	"tran dai quang",
	"nguyen tuan dung",
	"nguyen phuc trong",
	"cong san",
	"nha nuoc",
	"tuan hanh",
	"xuong duong",
	"viet tan",
	"bieu tinh",
	"vat gia",
	"nha dat",
	"khoa hoc",
	"hoi nghi",
	"marketing",
	"gia re",
	"gia si",
	"sale off",
	"sale",
	"Mocha",
	"uu dai",
	"Instagram",
	"Youtube",
	"Twitter",
	"ola",
	"wechat",
	"tango",
	"viber",
	"zalo",
	"gia re nhat",
	"gia thap nhat",
	"dam bao",
	"khai truong",
	"san Y",
	"tham du",
	"hoan tien",
	"tham khao",
	"mien phi",
	"rut tham",
	"may man",
	"qua tang",
	"trung thuong",
	"tang them",
	"trao thuong",
	"ra mat",
	"san pham moi",
	"gioi thieu san pham",
	"gioi thieu",
	"tự",
	"quang cao",
}

var modifyString = []string{
	"qc",
	"km",
	"khuyen mai",
	"facebook",
	"chiet khau",
	"danh tang",
	"tang ngay",
	"qua tang",
	"hoc bong",
	"co hoi boc tham",
	"rut tham",
	"giam",
	"giam toi",
	"giam den",
	"giam gia",
	"giam ngay",
	"giam hoc phi",
	"giam hphi",
	"uu dai",
	"co hoi nhan ngay",
	"sale off",
	"sale",
	"kmai",
	"uu-dai",
	"giam-gia",
	"sale-off",
	"la ngay cuoi ap dung",
}

func modifyCensoredWord(word string) string {
	if len(word) <= 1 {
		return word
	}
	var wordSplits []string
	wordSplits = strings.Split(word, space)
	var lastWord = wordSplits[len(wordSplits)-1]
	var lastWordBytes []byte
	for k, _ := range lastWord {
		if k == 0 {
			lastWordBytes = append(lastWordBytes, lastWord[k])
		} else {
			lastWordBytes = append(lastWordBytes, dot[0])
		}
	}
	wordSplits[len(wordSplits)-1] = string(lastWordBytes)
	return strings.Join(wordSplits, space)
}

func removeExcssiveSpace(s string) string {
	var listBytes []byte
	var spaceByte = space[0]
	for k, _ := range s {
		if k == 0 && s[k] == spaceByte {
			continue
		}
		if s[k] == spaceByte && k+1 < len(s) && s[k+1] == spaceByte {
			continue
		}
		if s[k] == spaceByte && k == len(s)-1 {
			continue
		}
		listBytes = append(listBytes, s[k])
	}
	return string(listBytes)
}

func ModifyMsgPhone(s string) string {
	s = removeAccent(s)
	lowCase := strings.ToLower(s)
	for _, v := range removeWord {
		// TODO case Not compare lowcase
		if strings.Contains(s, v) {
			s = strings.ReplaceAll(s, v, empty)
		}
		// TODO case compare lowcase
		// if strings.Contains(lowCase, v) {
		// 	runes := []rune(s)
		// 	stringModify := string(runes[strings.Index(lowCase, v) : strings.Index(lowCase, v)+len(v)])
		// 	s = strings.ReplaceAll(s, stringModify, empty)
		// }
	}
	lowCase = strings.ToLower(s)
	for _, v := range modifyString {
		if strings.Contains(lowCase, v) {
			runes := []rune(s)
			stringModify := string(runes[strings.Index(lowCase, v) : strings.Index(lowCase, v)+len(v)])
			s = strings.ReplaceAll(s, stringModify, modifyCensoredWord(stringModify))
		}
	}
	return removeExcssiveSpace(s)
}

// copyright https://ereka.vn/post/chuyen-tieng-viet-co-dau-thanh-khong-dau-trong-golang-409502490910618438
// Mang cac ky tu goc co dau copy
var SOURCE_CHARACTERS, LL_LENGTH = stringToRune(`ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰựỸ`)

// Mang cac ky tu thay the khong dau
var DESTINATION_CHARACTERS, _ = stringToRune(`AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUuY`)

func stringToRune(s string) ([]string, int) {
	ll := utf8.RuneCountInString(s)
	var texts = make([]string, ll+1)
	var index = 0
	for _, runeValue := range s {
		texts[index] = string(runeValue)
		index++
	}
	return texts, ll
}

func binarySearch(sortedArray []string, key string, low int, high int) int {
	var middle = (low + high) / 2
	if high < low {
		return -1
	}
	if key == sortedArray[middle] {
		return middle
	} else if key < sortedArray[middle] {
		return binarySearch(sortedArray, key, low, middle-1)
	} else {
		return binarySearch(sortedArray, key, middle+1, high)
	}
}

func removeAccentChar(ch string) string {
	var index = binarySearch(SOURCE_CHARACTERS, ch, 0, LL_LENGTH)
	if index >= 0 {
		ch = DESTINATION_CHARACTERS[index]
	}
	return ch
}

func removeAccent(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		buffer.WriteString(removeAccentChar(string(runeValue)))
	}
	return buffer.String()

}
