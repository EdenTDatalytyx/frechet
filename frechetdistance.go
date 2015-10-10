package frechet

import (
	"github.com/artpar/frechet/deque"
	"math"
)

type FrechetDistance interface {
	ComputeDistance(p, q [][]float64) float64
	initializeColumnUpperEnvelope(column int) UpperEnvelope
	initializeRowUpperEnvelope(row int) UpperEnvelope
	distance(p, q []float64) float64
}

type abstractFrechetDistance struct {
	n, m     int
	p, q     [][]float64
	instance FrechetDistance
}

func (x abstractFrechetDistance) ComputeDistance(p, q [][]float64) float64 {
	x.p = p
	x.q = q
	x.n = len(x.p) - 1
	x.m = len(x.q) - 1
	dist := x.compute()
	x.p = nil
	x.q = nil
	return dist
}

func (x abstractFrechetDistance) compute() float64 {
	column_queues := make([]*deque.Deque, x.n)
	column_envelopes := make([]UpperEnvelope, x.n)
	for i, _ := range column_queues {
		column_queues[i] = deque.New()
		column_envelopes[i] = x.instance.initializeColumnUpperEnvelope(i)
	}
	row_queues := make([]*deque.Deque, x.m)
	row_envelopes := make([]UpperEnvelope, x.m)
	for i, _ := range column_queues {
		temp := deque.New()
		row_queues[i] = temp
		row_envelopes[i] = x.instance.initializeRowUpperEnvelope(i)
	}

	L_opt := make([][]float64, x.n)
	for i, _ := range L_opt {
		L_opt[i] = make([]float64, x.m)
	}
	L_opt[0][0] = x.instance.distance(x.p[0], x.q[0]);
	for j := 1; j < x.m; j++ {
		L_opt[0][j] = math.MaxInt64;
	}

	B_opt := make([][]float64, x.n)
	for i, _ := range B_opt {
		L_opt[i] = make([]float64, x.m)
	}
	B_opt[0][0] = L_opt[0][0]
	for j := 1; j < x.m; j++ {
		B_opt[0][j] = math.MaxInt64;
	}

	for i := 0; i < x.n; i++ {
		for j := 0; j < x.m; j++ {
			if i < x.n - 1 {
				queue := row_queues[j]
				upperenv := row_envelopes[j]
				for queue.Size() > 0 && B_opt[queue.Right().(int)][j] > B_opt[i][j] {
					queue.PopRight()
				}
				queue.PushRight(i)
				if queue.Size() == 1 {
					upperenv.Clear()
				}
				upperenv.Add(i + 1, x.q[j], x.q[j + 1], x.p[i + 1])
				h := queue.Left().(int)
				min := upperenv.FindMinimum(B_opt[h][j])
				if h < i {
					min = upperenv.FindMinimum(L_opt[i][j], B_opt[h][j])
				}

				first := queue.PopLeft()
				for queue.Size() > 1 && B_opt[queue.Left().(int)][j] <= min {
					h := queue.Left().(int)
					if h <= i {
						panic("h !<= i")
					}
					upperenv.RemoveUpto(h)
					min = upperenv.FindMinimum(B_opt[h][j])
					if h < i {
						min = upperenv.FindMinimum(L_opt[i][j], B_opt[h][j])
					}
					first = queue.PopLeft()
				}
				queue.PushLeft(first)

				L_opt[i + 1][j] = min;
				upperenv.TruncateLast();
			}

			if j < x.m - 1 {
				queue := column_queues[j]
				upperenv := column_envelopes[j]
				for queue.Size() > 0 && B_opt[queue.Right().(int)][j] > B_opt[i][j] {
					queue.PopRight()
				}
				queue.PushRight(i)
				if queue.Size() == 1 {
					upperenv.Clear()
				}
				upperenv.Add(i + 1, x.q[j], x.q[j + 1], x.p[i + 1])
				h := queue.Left().(int)
				min := upperenv.FindMinimum(L_opt[h][j])
				if h < i {
					min = upperenv.FindMinimum(B_opt[i][j], L_opt[h][j])
				}

				first := queue.PopLeft()
				for queue.Size() > 1 && L_opt[queue.Left().(int)][j] <= min {
					h := queue.Left().(int)
					if h <= i {
						panic("h !<= i")
					}
					upperenv.RemoveUpto(h)
					min = upperenv.FindMinimum(B_opt[h][j])
					if h < i {
						min = upperenv.FindMinimum(B_opt[i][j], L_opt[h][j])
					}
					first = queue.PopLeft()
				}
				queue.PushLeft(first)

				B_opt[i + 1][j] = min;
				upperenv.TruncateLast();
			}
		}
	}
	return math.Max(x.instance.distance(x.p[x.n], x.q[x.m]), math.Min(L_opt[x.n - 1][x.m - 1], B_opt[x.n - 1][x.m - 1]))
}