package frechet
import (
	"math"
	"github.com/artpar/frechet/vectorutil"
)

type PolyhedralDistanceFunction struct {
	facets         [][]float64;
	facetSqrLength []float64;
}

func NewPolyhedralDistanceFunction(facets [][]float64) PolyhedralDistanceFunction {
	p := PolyhedralDistanceFunction{facets:facets}
	p.facetSqrLength = make([]float64, len(facets))
	for i := 0; i < len(p.facetSqrLength); i++ {
		p.facetSqrLength[i] = vectorutil.SquaredLength(p.facets[i]);
	}
	return p
}

func ( p PolyhedralDistanceFunction) Complexity() int {
	return len(p.facets)
}
func (p PolyhedralDistanceFunction) Facet(i int) []float64 {
	return p.facets[i]
}
func (this PolyhedralDistanceFunction) DistanceFromTow(p, q []float64) float64 {
	return this.Distance(vectorutil.Subtract(q, p))
}

func (p PolyhedralDistanceFunction) Distance(d []float64) float64 {
	max := -1 * math.MaxFloat64
	for i := 0; i < len(p.facets); i++ {
		fd := p.getFacetDistance(d, i);
		max = math.Max(max, fd);
	}
	return max
}

func (this PolyhedralDistanceFunction) getFacetDistanceFromTwo(p, q []float64, facet int) float64 {
	return this.getFacetDistance(vectorutil.Subtract(q, p), facet);
}

func (p PolyhedralDistanceFunction) getFacetDistance(d []float64, facet int) float64 {
	return vectorutil.DotProduct(p.facets[facet], d) / p.facetSqrLength[facet];
}

func (p PolyhedralDistanceFunction) getFacetSlope(p1, p2 []float64, facet int) float64 {
	return p.getFacetDistance(vectorutil.Subtract(p2, p1), facet);
}

func Custom(facetNormals, facetPoints [][]float64, normalize bool) PolyhedralDistanceFunction {
	if len(facetNormals) != len(facetPoints) {
		panic("facet normal and facet points should have same length")
	}
	if len(facetNormals[0]) != len(facetPoints[0]) {
		panic("must have same dimentions")
	}
	facetdescripts := make([][]float64, len(facetNormals));

	for i := 0; i < len(facetdescripts); i++ {
		N := facetNormals[i];
		if (normalize) {
			N = vectorutil.Normalise(N);
		}
		facetdescripts[i] = vectorutil.Scale(vectorutil.DotProduct(facetNormals[i], facetPoints[i]), facetNormals[i]);
	}

	return NewPolyhedralDistanceFunction(facetdescripts);
}
func KRegular2D(k int) PolyhedralDistanceFunction {
	if !(k >= 4 && k % 2 == 0) {
		panic("k must be even number greater then or equal to 4")
	}

	facetdescripts := make([][]float64, k);
	for i, _ := range facetdescripts {
		facetdescripts[i] = make([]float64, 2)
	}
	alpha := float64(2) * math.Pi / float64(k);

	for i := 0; i < k; i++ {
		facetdescripts[i][0] = 0.5 * (math.Cos(float64(i) * alpha) + math.Cos(float64(i + 1) * alpha));
		facetdescripts[i][1] = 0.5 * (math.Sin(float64(i) * alpha) + math.Sin(float64(i + 1) * alpha));
	}

	return NewPolyhedralDistanceFunction(facetdescripts);
}

func EpsApproximation2D(eps float64) PolyhedralDistanceFunction {
	if !(eps > 1) {
		panic("eps is not greater then 1")
	}

	var k int;
	if (eps >= math.Sqrt(2)) {
		k = 4;
	} else {
		k = int(math.Ceil(math.Pi * 2.0 / math.Acos(1.0 / eps)));
		if (k % 2 == 1) {
			k++;
		}
	}

	return KRegular2D(k);
}

func LInfinity(dimensions int) PolyhedralDistanceFunction {
	if !(dimensions >= 2) {
		panic("dimensions must be greater then or equal to 2")
	}

	facetdescripts := make([][]float64, 2 * dimensions);
	for i, _ := range facetdescripts {
		facetdescripts[i] = make([]float64, dimensions)
	}
	for i := 0; i < dimensions; i++ {
		for j := 0; j < dimensions; j++ {
			if i == j {
				facetdescripts[2 * i][j] = 1
				facetdescripts[2 * i + 1][j] = -1
			} else {
				facetdescripts[2 * i][j] = 0
				facetdescripts[2 * i + 1][j] = 0
			}
		}
	}

	return NewPolyhedralDistanceFunction(facetdescripts);
}

func Round(x float64) float64 {
	v, frac := math.Modf(x)
	if x > 0.0 {
		if frac > 0.5 || (frac == 0.5 && uint64(v) % 2 != 0) {
			v += 1.0
		}
	} else {
		if frac < -0.5 || (frac == -0.5 && uint64(v) % 2 != 0) {
			v -= 1.0
		}
	}

	return v
}

func L1(dimensions int) PolyhedralDistanceFunction {
	if !(dimensions >= 2) {
		panic("dimennsions should be more then or equal to 2")
	}

	k := int(Round(math.Pow(2, float64(dimensions))));
	facetdescripts := make([][]float64, k);
	for i, _ := range facetdescripts {
		facetdescripts[i] = make([]float64, dimensions)
	}

	val := float64(1.0 / dimensions);

	totalblock := k;
	for d := 0; d < dimensions; d++ {
		halfblock := totalblock / 2;

		for f := 0; f < k; f++ {
			if (f % totalblock < halfblock) {
				facetdescripts[f][d] = val;
			} else {
				facetdescripts[f][d] = -val;
			}
		}

		totalblock = halfblock;
	}

	return NewPolyhedralDistanceFunction(facetdescripts);
}


