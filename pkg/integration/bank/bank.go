package bank

import (
	"fmt"

	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

var ll = l.New()

type BankType struct {
	Loai            string
	TenLoaiNganHang string
}

type Bank struct {
	MaNganHang string
	TenNH      string
	Loai       string
}

type Province struct {
	MaTinh       string
	TenTinhThanh string
	MaNganHang   string
}

type Branch struct {
	MaChiNhanh  string
	TenChiNhanh string
	MaNganHang  string
	MaTinh      string
}

var (
	bankTypeIndexCode = make(map[string]*BankType)
	bankTypeIndexNamN = make(map[string]*BankType)

	bankIndexCode     = make(map[string]*Bank)
	bankIndexNamN     = make(map[string]*Bank)
	bankIndexBankType = make(map[string][]*Bank)

	provinceIndexBankCode = make(map[string][]*Province)
	proviceIndexCode      = make(map[string]*Province)
	proviceIndexNamn      = make(map[string]*Province)

	branchIndexCode                 = make(map[string]*Branch)
	branchIndexNamN                 = make(map[string]*Branch)
	branchIndexBankCodeProvinceCode = make(map[string][]*Branch)
)

func init() {
	for _, bType := range BankTypes {
		nameNorm := validate.NormalizeSearchSimple(bType.TenLoaiNganHang)
		bankTypeIndexNamN[nameNorm] = bType
		bankTypeIndexCode[bType.Loai] = bType
	}

	for _, bank := range Banks {
		nameNorm := validate.NormalizeSearchSimple(bank.TenNH)
		bankIndexNamN[nameNorm] = bank
		bankIndexCode[bank.MaNganHang] = bank
		bankIndexBankType[bank.Loai] = append(bankIndexBankType[bank.Loai], bank)
	}

	for _, province := range Provinces {
		nameNorm := validate.NormalizeSearchSimple(province.TenTinhThanh)
		proviceIndexCode[province.MaTinh] = province
		proviceIndexNamn[nameNorm] = province
		provinceIndexBankCode[province.MaNganHang] = append(provinceIndexBankCode[province.MaNganHang], province)
	}

	for _, branch := range Branches {
		nameNorm := validate.NormalizeSearchSimple(branch.TenChiNhanh)
		branchIndexNamN[nameNorm] = branch
		branchIndexCode[branch.MaChiNhanh] = branch

		index := branch.MaNganHang + "_" + branch.MaTinh
		branchIndexBankCodeProvinceCode[index] = append(branchIndexBankCodeProvinceCode[index], branch)
	}
}

type BankQuery struct {
	Code string
	Name string
}

type ProvinceQuery struct {
	Code string
	Name string
}

func GetProvinceByBank(query *BankQuery) []*Province {
	var res []*Province
	bankCode, bankName := query.Code, query.Name
	if bankCode == "" && bankName == "" {
		return res
	}
	if bankCode != "" {
		if _, ok := bankIndexCode[bankCode]; !ok {
			return res
		}
		return provinceIndexBankCode[bankCode]
	}

	if bankName != "" {
		norm := validate.NormalizeSearchSimple(bankName)
		bank, ok := bankIndexNamN[norm]

		if !ok {
			return res
		}
		return provinceIndexBankCode[bank.MaNganHang]
	}
	return res
}

func GetBankProvinces(bankQuery *BankQuery, getAll bool) []*Province {
	var res []*Province
	bankCode, bankName := bankQuery.Code, bankQuery.Name

	if getAll {
		return Provinces
	}

	if bankCode != "" {
		if result, ok := provinceIndexBankCode[bankCode]; ok {
			return result
		}
	} else if bankName != "" {
		norm := validate.NormalizeSearchSimple(bankName)
		nameBank, exist := bankIndexNamN[norm]
		if !exist {
			return res
		}
		if result, ok := provinceIndexBankCode[nameBank.MaNganHang]; ok {
			return result
		}
	}

	return res
}

func GetBranchByBankProvince(bankQuery *BankQuery, provinceQuery *ProvinceQuery) []*Branch {
	var res []*Branch
	bankCode, bankName := bankQuery.Code, bankQuery.Name
	provinceCode, provinceName := provinceQuery.Code, provinceQuery.Name

	if bankCode == "" && bankName == "" {
		return res
	}
	if provinceCode == "" && provinceName == "" {
		return res
	}
	var bank *Bank
	var province *Province
	if bankCode != "" {
		bank = bankIndexCode[bankCode]
	} else if bankName != "" {
		norm := validate.NormalizeSearchSimple(bankName)
		bank = bankIndexNamN[norm]
	}

	if provinceCode != "" {
		province = proviceIndexCode[provinceCode]
	} else if provinceName != "" {
		norm := validate.NormalizeSearchSimple(provinceName)
		province = proviceIndexNamn[norm]
	}

	if bank == nil || province == nil {
		return res
	}
	index := bank.MaNganHang + "_" + province.MaTinh
	return branchIndexBankCodeProvinceCode[index]
}

func GetBankBranches(bankQuery *BankQuery, getAll bool) []*Branch {
	var res []*Branch
	bankCode, bankName := bankQuery.Code, bankQuery.Name

	if getAll {
		return Branches
	}

	if bankCode != "" {
		if bank, exist := bankIndexCode[bankCode]; exist {
			if provinces, ok := provinceIndexBankCode[bank.MaNganHang]; ok {
				res = append(res, GetBankBranchesFromProvinces(provinces)...)
			}
		}
	} else if bankName != "" {
		norm := validate.NormalizeSearchSimple(bankName)
		bank, ok := bankIndexNamN[norm]
		if !ok {
			return res
		}
		provinces, ok := provinceIndexBankCode[bank.MaNganHang]
		if !ok {
			return res
		}

		res = append(res, GetBankBranchesFromProvinces(provinces)...)
	}

	return res
}

func GetBankBranchesFromProvinces(provinces []*Province) (res []*Branch) {
	if len(provinces) == 0 {
		return res
	}
	for _, province := range provinces {
		key := fmt.Sprintf("%s_%s", province.MaNganHang, province.MaTinh)
		if data, ok := branchIndexBankCodeProvinceCode[key]; ok {
			res = append(res, data...)
		}
	}

	return res
}
