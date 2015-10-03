package frechet

type UpperEnvelope interface {
	Add(i int, P1, P2, Q []float64)
	RemoveUpto(i int);
	Clear();
	FindMinimum(constants ...float64) float64;
	TruncateLast();
}