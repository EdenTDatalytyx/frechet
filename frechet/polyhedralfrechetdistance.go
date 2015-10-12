package frechet

type PolyhedralFrechetDistance struct {
	AbstractFretchetDistance
	distfunc PolyhedralDistanceFunction
}

func NewPolyhedralFretchetDistance(disfunc PolyhedralDistanceFunction) PolyhedralFrechetDistance {
	x := PolyhedralFrechetDistance{distfunc:disfunc}
	x.FrechetDistance = x
	return x
}



func (this PolyhedralFrechetDistance) distance(p, q []float64) float64 {
	return this.distfunc.DistanceFromTwo(p, q);
}


func (this PolyhedralFrechetDistance) initializeRowUpperEnvelope(row int, Q [][]float64) UpperEnvelope {
	// fmt.Printf("initializeRowUpperEnvelope(%d) - Q == %s\n", row, Q)
	return NewPolyhedralUpperEnvelope(this.distfunc, Q[row], Q[row + 1]);
}

func (this PolyhedralFrechetDistance) initializeColumnUpperEnvelope(column int, P [][]float64) UpperEnvelope {
	// fmt.Printf("initializeColumnUpperEnvelope(%d) - P == %s\n", column, P)
	return NewPolyhedralUpperEnvelope(this.distfunc, P[column], P[column + 1]);
}