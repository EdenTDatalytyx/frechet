package frechet

import (

	"github.com/artpar/frechet/deque"
	"sort"
	"github.com/artpar/frechet/vectorutil"
	"reflect"
	"math"
)

type UpperEnvelope interface {
	Add(i int, P1, P2, Q []float64)
	RemoveUpto(i int);
	Clear();
	FindMinimum(constants ...float64) float64;
	TruncateLast();
}

type facetList struct {
	deque.Deque
	facet int
	slope float64
}

func newFacetList(facet int, slope float64) facetList {
	f := facetList{facet:facet, slope: slope}
	return f
}

type PolyhedralUpperEnvelope struct {
	distfunc     PolyhedralDistanceFunction
	p1, p2       []float64
	sortedfacets []facetList
}

type BySlope []facetList

func (a BySlope) Len() int { return len(a) }
func (a BySlope) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySlope) Less(i, j int) bool { return a[i].slope < a[j].slope }

func NewPolyhedralUpperEnvelope(distfunc PolyhedralDistanceFunction, p1, p2 []float64) PolyhedralUpperEnvelope {
	e := PolyhedralUpperEnvelope{distfunc:distfunc, p1:p1, p2:p2}
	e.sortedfacets = make([]facetList, 0)
	for i := 0; i < distfunc.Complexity(); i++ {
		e.sortedfacets = append(e.sortedfacets, newFacetList(i, distfunc.getFacetSlope(p1, p2, i)))
	}

	sort.Sort(BySlope(e.sortedfacets))

	return e
}

type FacetListElement struct {
	index  int
	height float64
	slope  float64
}



func (this PolyhedralUpperEnvelope) Add(i int, P1, P2, Q []float64) {
	if !(reflect.DeepEqual(this.p1, P1)) {
		panic("p1 is not equal to P1")
	}
	if !(reflect.DeepEqual(this.p2, P2)) {
		panic("p2 is not equal to P2")
	}

	for f := 0; f < this.distfunc.Complexity(); f++ {

		fl := this.sortedfacets[f];

		fle := FacetListElement{index:i};
		fle.slope = fl.slope;
		fle.height = this.distfunc.getFacetDistance(vectorutil.Subtract(this.p1, Q), fl.facet);

		for ; fl.Size() > 0 && fl.Left().(FacetListElement).height <= fle.height; {
			fl.PopLeft();
		}
		fl.PushLeft(fle);
	}
}


func (this PolyhedralUpperEnvelope) RemoveUpto(i int) {
	for _, fl := range this.sortedfacets {

		for ; fl.Size() > 0 && fl.Right().(FacetListElement).index <= i; {
			fl.PopRight();
		}
		if !(fl.Size() > 0) {
			panic("how is this size 0 ?")
		}
	}
}

func (this PolyhedralUpperEnvelope) Clear() {
	for _, fl := range this.sortedfacets {
		fl.Empty()
	}
}

func (this PolyhedralUpperEnvelope) findMinimumFullProcedure() float64 {
	upperenvelope := make([]FacetListElement, this.distfunc.Complexity())

	first := this.sortedfacets[0].Right().(FacetListElement)
	upperenvelope[0] = first;
	intersectPreviousAt := make([]float64, this.distfunc.Complexity())
	intersectPreviousAt[0] = -1 * math.MaxFloat64;
	n := 1

	for i := 1; i < this.distfunc.Complexity(); i++ {
		fl := this.sortedfacets[i]
		fle := fl.PopRight().(FacetListElement);
		upperenvelope[n] = fle;
		intersectPreviousAt[n] = this.findIntersection(upperenvelope[n - 1], fle);
		n++;

		if (math.IsInf(intersectPreviousAt[n - 1], 0) || math.IsNaN(intersectPreviousAt[n - 1])) {
			// NB: can only occur with the first previous one (slopes are unique in UE)
			if (upperenvelope[n - 1].height > upperenvelope[n - 2].height) {
				// toss previous
				upperenvelope[n - 2] = upperenvelope[n - 1];
				if (n > 2) {
					intersectPreviousAt[n - 2] = this.findIntersection(upperenvelope[n - 3], fle);
				} else {
					// nothing to intersect with
				}
				n--;
			} else {
				// toss this one
				n--;
			}
		}

		for ; n > 1 && intersectPreviousAt[n - 1] < intersectPreviousAt[n - 2]; {
			// [n-2] not on upper envelope

			if !(n > 2) {
				panic("[n-2] not on upper envolope")
			}

			upperenvelope[n - 2] = upperenvelope[n - 1];
			intersectPreviousAt[n - 2] = this.findIntersection(upperenvelope[n - 3], fle);

			n--;
		}
	}

	// upperenvelope contains only lines that actually occur on it

	// minimum given by first point with positve slop
	min_index := -1;
	for i := 0; i < n; i++ {
		if (upperenvelope[i].slope > 0) {
			min_index = i;
			break;
		}
	}

	var min float64;
	if (intersectPreviousAt[min_index] < 0) {
		// minimum before interval [0,1]
		// find (first) intersection with 0 by going forward

		for ; min_index < n - 1 && intersectPreviousAt[min_index + 1] < 0; {
			min_index++;
		}
		min = upperenvelope[min_index].height;

	} else if (intersectPreviousAt[min_index] > 1) {
		// minimum after interval [0,1]
		// find first intersection with 1 by going backward

		for ; min_index > 0 && intersectPreviousAt[min_index] > 1; {
			min_index--;
		}
		min = upperenvelope[min_index].slope + upperenvelope[min_index].height;

	} else {
		// minimum point in interval [0,1]
		min = upperenvelope[min_index].slope * intersectPreviousAt[min_index] + upperenvelope[min_index].height;
	}

	return min;
}


func (this PolyhedralUpperEnvelope) findMinimumTrimmedProcedure() float64 {
	upperenvelope := make([]FacetListElement, this.distfunc.Complexity())

	first := this.sortedfacets[0].Right().(FacetListElement);
	upperenvelope[0] = first;
	intersectPreviousAt := make([]float64, this.distfunc.Complexity())
	intersectPreviousAt[0] = 0;
	n := 1;

	for i := 1; i < this.distfunc.Complexity(); i++ {
		fl := this.sortedfacets[i];
		fle := fl.Right().(FacetListElement)

		upperenvelope[n] = fle;
		intersectPreviousAt[n] = this.findIntersection(upperenvelope[n - 1], fle);
		n++;

		if (math.IsInf(intersectPreviousAt[n - 1], 0) || math.IsNaN(intersectPreviousAt[n - 1])) {
			// NB: can only occur with the first previous one (slopes are unique in UE)
			if (upperenvelope[n - 1].height > upperenvelope[n - 2].height) {
				// toss previous
				upperenvelope[n - 2] = upperenvelope[n - 1];
				if (n > 2) {
					intersectPreviousAt[n - 2] = this.findIntersection(upperenvelope[n - 3], fle);
				} else {
					// nothing to intersect
				}
				n--;
			} else {
				// toss this one
				n--;
			}
		}

		for ; (n > 1 && intersectPreviousAt[n - 1] < intersectPreviousAt[n - 2]); {
			// [n-2] not on upper envelope
			upperenvelope[n - 2] = upperenvelope[n - 1];
			if (n == 2) {
				// intersectPreviousAt[n-2] = 0; // doesn't change
			} else {
				intersectPreviousAt[n - 2] = this.findIntersection(upperenvelope[n - 3], fle);
			}
			n--;
		}

		if (intersectPreviousAt[n - 1] > 1) {
			if !(n > 1) {
				panic("was expecting n > 1")
			}
			n--;
		} else if (n > 1 && upperenvelope[n - 2].slope > 0) {
			n--;
		}
	}

	// upper envelope now contains the decreasing part starting at zero
	// up to the first increasing function before 1 on the upper envelope (if any)

	var min float64;
	if (upperenvelope[n - 1].slope > 0) {
		if (n > 1) {
			// get height at intersection with previous
			min = upperenvelope[n - 1].slope * intersectPreviousAt[n - 1] + upperenvelope[n - 1].height;
		} else {
			// get height at 0
			min = upperenvelope[n - 1].height;
		}
	} else {
		// get height at 1
		min = upperenvelope[n - 1].slope + upperenvelope[n - 1].height;
	}

	return min;
}


func (this PolyhedralUpperEnvelope) findIntersection(fle1, fle2 FacetListElement) float64 {
	// find x of intersection point (assume slopes are different)
	// y = fle1.slope * x + fle1.height
	// y = fle2.slope * x + fle1.height
	//
	// (fle1.slope - fle2.slope) * x = fle2.height - fle1.height
	// x = (fle2.height - fle1.height) / (fle1.slope - fle2.slope)

	xint := (fle2.height - fle1.height) / (fle1.slope - fle2.slope);
	return xint;
}