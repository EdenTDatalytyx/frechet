package vectorutil
import (
	"math"
	// "fmt"
)


func SquaredLength(A []float64) float64 {
	var sum float64;
	for i := 0; i < len(A); i++ {
		sum += A[i] * A[i];
	}
	// fmt.Printf("Squared Length of %s: %f\n", A, sum)
	return sum
}

func Add(A, B []float64) []float64 {
	if len(A) != len(B) {
		panic("Length are not equal")
	}
	R := make([]float64, len(A))
	for i, _ := range R {
		R[i] = A[i] + B[i]
	}
	return R
}


func Subtract(A, B []float64) []float64 {
	if len(A) != len(B) {
		panic("Length are not equal")
	}
	R := make([]float64, len(A))
	for i, _ := range R {
		R[i] = A[i] - B[i]
	}
	return R
}

func DotProduct(A, B []float64) float64 {
	if len(A) != len(B) {
		panic("Length are not equal")
	}
	var R float64
	for i, _ := range A {
		R += A[i] * B[i]
	}
	return R
}

func Distance(A, B []float64) float64 {
	// fmt.Printf("Check distance A,B == %s,%s", A, B)
	return Length(Subtract(B, A))
}

func Length(A []float64) float64 {
	return math.Sqrt(SquaredLength(A))
}

func Scale(S float64, A []float64) []float64 {
	R := make([]float64, len(A))
	for i, _ := range R {
		R[i] = A[i] * S
	}
	return R
}

func Normalise(A []float64) []float64 {
	return Scale(1 / Length(A), A)
}
