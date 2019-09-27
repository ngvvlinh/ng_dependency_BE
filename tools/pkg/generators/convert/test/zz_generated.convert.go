// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package test

import (
	scheme "etop.vn/backend/pkg/common/scheme"
)

/*
Custom conversions:
    ConvertAB       // in use
    ConvertC01      // in use
    ConvertC10      // in use

Ignored functions: (none)
*/

func init() {
	registerConversionFunctions(scheme.Global)
}

func registerConversionFunctions(s *scheme.Scheme) {
	s.Register((*B)(nil), (*A)(nil), func(arg, out interface{}) error {
		Convert_B_A(arg.(*B), out.(*A))
		return nil
	})
	s.Register(([]*B)(nil), (*[]*A)(nil), func(arg, out interface{}) error {
		out0 := Convert_Bs_As(arg.([]*B))
		*out.(*[]*A) = out0
		return nil
	})
	s.Register((*A)(nil), (*B)(nil), func(arg, out interface{}) error {
		Convert_A_B(arg.(*A), out.(*B))
		return nil
	})
	s.Register(([]*A)(nil), (*[]*B)(nil), func(arg, out interface{}) error {
		out0 := Convert_As_Bs(arg.([]*A))
		*out.(*[]*B) = out0
		return nil
	})
	s.Register((*C1)(nil), (*C0)(nil), func(arg, out interface{}) error {
		Convert_C1_C0(arg.(*C1), out.(*C0))
		return nil
	})
	s.Register(([]*C1)(nil), (*[]*C0)(nil), func(arg, out interface{}) error {
		out0 := Convert_C1s_C0s(arg.([]*C1))
		*out.(*[]*C0) = out0
		return nil
	})
	s.Register((*C0)(nil), (*C1)(nil), func(arg, out interface{}) error {
		Convert_C0_C1(arg.(*C0), out.(*C1))
		return nil
	})
	s.Register(([]*C0)(nil), (*[]*C1)(nil), func(arg, out interface{}) error {
		out0 := Convert_C0s_C1s(arg.([]*C0))
		*out.(*[]*C1) = out0
		return nil
	})
}

//-- convert etop.vn/backend/tools/pkg/generators/convert/test.A --//

func Convert_B_A(arg *B, out *A) *A {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &A{}
	}
	convert_B_A(arg, out)
	return out
}

func convert_B_A(arg *B, out *A) {
	out.Value = 0                   // types do not match
	out.Int = int64(arg.Int)        // simple conversion
	out.String = string(arg.String) // simple conversion
	out.Strings = arg.Strings       // simple assign
	out.C = Convert_C1_C0(arg.C, nil)
	out.Cs = Convert_C1s_C0s(arg.Cs)
}

func Convert_Bs_As(args []*B) (outs []*A) {
	tmps := make([]A, len(args))
	outs = make([]*A, len(args))
	for i := range tmps {
		outs[i] = Convert_B_A(args[i], &tmps[i])
	}
	return outs
}

func Convert_A_B(arg *A, out *B) *B {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &B{}
	}
	ConvertAB(arg, out)
	return out
}

func convert_A_B(arg *A, out *B) {
	out.Value = ""             // types do not match
	out.Int = int32(arg.Int)   // simple conversion
	out.String = S(arg.String) // simple conversion
	out.Strings = arg.Strings  // simple assign
	out.C = Convert_C0_C1(arg.C, nil)
	out.Cs = Convert_C0s_C1s(arg.Cs)
}

func Convert_As_Bs(args []*A) (outs []*B) {
	tmps := make([]B, len(args))
	outs = make([]*B, len(args))
	for i := range tmps {
		outs[i] = Convert_A_B(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/backend/tools/pkg/generators/convert/test.C0 --//

func Convert_C1_C0(arg *C1, out *C0) *C0 {
	return ConvertC10(arg, out)
}

func convert_C1_C0(arg *C1, out *C0) {
	out.Value = 0 // types do not match
}

func Convert_C1s_C0s(args []*C1) (outs []*C0) {
	tmps := make([]C0, len(args))
	outs = make([]*C0, len(args))
	for i := range tmps {
		outs[i] = Convert_C1_C0(args[i], &tmps[i])
	}
	return outs
}

func Convert_C0_C1(arg *C0, out *C1) *C1 {
	return ConvertC01(arg)
}

func convert_C0_C1(arg *C0, out *C1) {
	out.Value = "" // types do not match
}

func Convert_C0s_C1s(args []*C0) (outs []*C1) {
	tmps := make([]C1, len(args))
	outs = make([]*C1, len(args))
	for i := range tmps {
		outs[i] = Convert_C0_C1(args[i], &tmps[i])
	}
	return outs
}
