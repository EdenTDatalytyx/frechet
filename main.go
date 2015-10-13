package main
import (
	fretchet "github.com/artpar/frechet/frechet"
// "fmt"
	"fmt"
)


func main() {
	var frechet fretchet.FrechetDistance;
	var curveA, curveB [][]float64;
	var dist float64;

	// two curves in 3D
	curveA = [][]float64{[]float64{0, -497.75619757841895}, []float64{559.4405594405595, -592.7969328517269}, []float64{1153.8461538461538, -746.6108377369661}, []float64{1748.2517482517483, -909.9922089511371}, []float64{2342.6573426573427, -1000}, []float64{2937.062937062937, -949.8966475584696}, []float64{3531.4685314685316, -719.3743228266028}, []float64{4090.909090909091, -439.86293936085747}, []float64{4685.314685314685, -227.8895214941}, []float64{5279.72027972028, -109.157671952481}, []float64{5874.125874125874, -42.826627522605804}, []float64{6468.531468531469, -8.021665145351108}, []float64{7062.937062937063, 77.1438213019494}, []float64{7622.377622377623, 264.29833021335367}, []float64{8216.783216783217, 533.4807776278296}, []float64{8811.188811188811, 866.6588667982098}, []float64{9405.594405594406, 1000}, []float64{10000, 790.0316097458776}};
	curveB = [][]float64{[]float64{0, -695.1611497647639}, []float64{544.8717948717949, -790.4777475298525}, []float64{1089.7435897435898, -892.1572902869895}, []float64{1634.6153846153845, -958.3213539384153}, []float64{2179.4871794871797, -1000}, []float64{2724.358974358974, -960.1329406719366}, []float64{3269.230769230769, -826.0042769490638}, []float64{3782.051282051282, -602.0035693135744}, []float64{4326.923076923077, -413.717890336045}, []float64{4871.794871794872, -299.96868420084854}, []float64{5416.666666666667, -253.97993813506912}, []float64{5961.538461538462, -217.38757274473835}, []float64{6506.410256410257, -175.92754425757016}, []float64{7051.282051282051, -70.35247519388975}, []float64{7564.102564102564, 104.99846910576912}, []float64{8108.974358974359, 345.94199721600626}, []float64{8653.846153846154, 684.4321262067106}, []float64{9198.71794871795, 983.762187449438}, []float64{9743.589743589744, 999.9999999999998}, []float64{10000, 699.4095383377385}};

	// L-1 in 3 dimensions
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.L1(2));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 1 : %f\n\n\n\n\n", dist)
	// two curves in 4D

	// L-infinity in 4 dimensions
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.LInfinity(2));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 2 : %f\n\n\n\n\n", dist)

	// two curves in 2D
	// 1.1-approximation of Euclidean (in 2 dimensions) (NB: any value above sqrt(2) uses sqrt(2) as approximation value)
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.EpsApproximation2D(1.009));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 3 : %f\n\n\n\n\n", dist)

	// 6-regular polygon (in 2 dimensions)
	// implementation supports only symmetric polyhedra, so parameter must be even!
	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.KRegular2D(8));
	dist = frechet.ComputeDistance(curveA, curveB);
	fmt.Printf("Distance 4 : %f\n", dist)
}

//
//func main() {
//	var frechet fretchet.FrechetDistance;
//	var curveA, curveB [][]float64;
//	var dist float64;
//
//	// two curves in 3D
//	curveA = [][]float64{[]float64{0, 0, 0}, []float64{1, 1, 0}, []float64{0, 1, 2}, []float64{2, 1, 2}};
//	curveB = [][]float64{[]float64{0, 0, 0}, []float64{2, 1, 2}};
//
//	// L-1 in 3 dimensions
//	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.L1(3));
//	dist = frechet.ComputeDistance(curveA, curveB);
//	fmt.Printf("Distance 1 : %f\n\n\n\n\n", dist)
//	// two curves in 4D
//	curveA = [][]float64{[]float64{0, 0, 0, -3}, []float64{1, 1, 6, 5}, []float64{0, 8, 2, -2}};
//	curveB = [][]float64{[]float64{0, 0, 0, 1}, []float64{2, 1, 2, 7}};
//
//	// L-infinity in 4 dimensions
//	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.LInfinity(4));
//	dist = frechet.ComputeDistance(curveA, curveB);
//	fmt.Printf("Distance 2 : %f\n\n\n\n\n", dist)
//
//	// two curves in 2D
//	curveA = [][]float64{[]float64{0, 0}, []float64{1, 6}, []float64{0, 8}};
//	curveB = [][]float64{[]float64{1, 0}, []float64{2, 7}, []float64{-1, 5}};
//
//	// 1.1-approximation of Euclidean (in 2 dimensions) (NB: any value above sqrt(2) uses sqrt(2) as approximation value)
//	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.EpsApproximation2D(1.1));
//	dist = frechet.ComputeDistance(curveA, curveB);
//	fmt.Printf("Distance 3 : %f\n\n\n\n\n", dist)
//
//	// 6-regular polygon (in 2 dimensions)
//	// implementation supports only symmetric polyhedra, so parameter must be even!
//	frechet = fretchet.NewPolyhedralFretchetDistance(fretchet.KRegular2D(6));
//	dist = frechet.ComputeDistance(curveA, curveB);
//	fmt.Printf("Distance 4 : %f\n", dist)
//}
//
