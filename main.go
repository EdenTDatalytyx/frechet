package main
import (
	fretchet "github.com/artpar/frechet/frechet"
	"fmt"
)


func main() {
	var frechet fretchet.FrechetDistance;
	var curveA, curveB [][]float64;
	var dist float64;

	// two curves in 3D
	curveA = [][]float64{[]float64{0, 0, 0}, []float64{1, 1, 0}, []float64{0, 1, 2}, []float64{2, 1, 2}};
	curveB = [][]float64{[]float64{0, 0, 0}, []float64{2, 1, 2}};

	// L-1 in 3 dimensions
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.L1(3));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 1 : %f\n\n\n\n\n", dist)
	// two curves in 4D
	curveA = [][]float64{[]float64{0, 0, 0, -3}, []float64{1, 1, 6, 5}, []float64{0, 8, 2, -2}};
	curveB = [][]float64{[]float64{0, 0, 0, 1}, []float64{2, 1, 2, 7}};

	// L-infinity in 4 dimensions
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.LInfinity(4));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 2 : %f\n\n\n\n\n", dist)

	// two curves in 2D
	curveA = [][]float64{[]float64{0, 0}, []float64{1, 6}, []float64{0, 8}};
	curveB = [][]float64{[]float64{1, 0}, []float64{2, 7}, []float64{-1, 5}};

	// 1.1-approximation of Euclidean (in 2 dimensions) (NB: any value above sqrt(2) uses sqrt(2) as approximation value)
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.EpsApproximation2D(1.1));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 3 : %f\n\n\n\n\n", dist)

	// 6-regular polygon (in 2 dimensions)
	// implementation supports only symmetric polyhedra, so parameter must be even!
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.KRegular2D(6));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 4 : %f\n", dist)
}