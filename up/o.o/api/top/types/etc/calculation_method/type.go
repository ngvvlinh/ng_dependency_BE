package calculation_method

// +enum
// +enum:zero=null
type CalculationMethodType int

type NullCalculationMethodType struct {
	Enum  CalculationMethodType
	Valid bool
}

// Context:
// - Có một mảng các công thức tính toán (theo rule gì đó định nghĩa sẵn)
// - Có 2 cách để tính ra kết quả:
//      + cumulative (Lũy kế): công thức nào thõa mãn thì tính ra kết quả, sau đó cộng dồn các kết quả này lại.
//      + first_satisfy (công thức đầu tiên thõa mãn): trả về kết quả của công thức thõa mãn đầu tiên

const (
	// +enum=unknown
	Unknown CalculationMethodType = 0

	// +enum=cumulative
	Cumulation CalculationMethodType = 1

	// +enum=first_satisfy
	FirstSatisfy CalculationMethodType = 2
)
