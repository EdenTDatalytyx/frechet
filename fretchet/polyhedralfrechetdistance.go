package frechet

type PolyhedralFrechetDistance struct {
	AbstractFrechetDistance
	distfunc PolyhedralDistanceFunction
}

func NewPolyhedralFrechetDistance(disfunc PolyhedralDistanceFunction) PolyhedralFrechetDistance {
	return PolyhedralFrechetDistance{distfunc:disfunc}
}



func (this PolyhedralFrechetDistance) distance(p, q []float64) float64 {
	return this.distfunc.DistanceFromTow(p, q);
}


func (this PolyhedralFrechetDistance) initializeRowUpperEnvelope(row int) UpperEnvelope {
	return NewPolyhedralUpperEnvelope(this.distfunc, this.AbstractFrechetDistance.Q[row], this.Q[row + 1]);
}

func (this PolyhedralFrechetDistance) initializeColumnUpperEnvelope(column int) UpperEnvelope {
	return NewPolyhedralUpperEnvelope(this.distfunc, this.Q[column], this.P[column + 1]);
}