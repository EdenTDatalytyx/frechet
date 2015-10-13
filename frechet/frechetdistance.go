package frechet

import (
	"github.com/artpar/frechet/deque"
	"math"
// "fmt"
)

type FrechetDistance interface {
	ComputeDistance(p, q [][]float64) float64
	initializeColumnUpperEnvelope(column int, P [][]float64) UpperEnvelope
	initializeRowUpperEnvelope(row int, Q [][]float64) UpperEnvelope
	distance(p, q []float64) float64
}

type AbstractFretchetDistance struct {
	n, m int
	P, Q [][]float64
	FrechetDistance
}

func (x AbstractFretchetDistance) ComputeDistance(p, q [][]float64) float64 {
	x.P = p
	x.Q = q
//	fmt.Printf("P, Q = %s, %s\n", x.P, x.Q)
	x.n = len(x.P) - 1
	x.m = len(x.Q) - 1
	dist := x.compute()
	x.P = nil
	x.Q = nil
	return dist
}

func (x AbstractFretchetDistance) compute() float64 {
	column_queues := make([]*deque.Deque, x.n)
	column_envelopes := make([]UpperEnvelope, x.n)
	for i, _ := range column_queues {
		column_queues[i] = deque.New()
		// fmt.Printf("compute(%d) - P, Q == %s, %s\n", i, x.P, x.Q)
		column_envelopes[i] = x.initializeColumnUpperEnvelope(i, x.P)
	}
	row_queues := make([]*deque.Deque, x.m)
	row_envelopes := make([]UpperEnvelope, x.m)
	for i, _ := range row_queues {
		temp := deque.New()
		// fmt.Printf("rowQueues(%d) P,Q == %s, %s\n", i)
		row_queues[i] = temp
		row_envelopes[i] = x.initializeRowUpperEnvelope(i, x.Q)
	}

	L_opt := make([][]float64, x.n)
	for i, _ := range L_opt {
		L_opt[i] = make([]float64, x.m)
	}
	L_opt[0][0] = x.distance(x.P[0], x.Q[0]);
	for j := 1; j < x.m; j++ {
		L_opt[0][j] = math.MaxFloat64;
	}

	B_opt := make([][]float64, x.n)
	for i, _ := range B_opt {
		B_opt[i] = make([]float64, x.m)
	}
	B_opt[0][0] = L_opt[0][0]
	for j := 1; j < x.n; j++ {
		B_opt[j][0] = math.MaxFloat64;
	}
	// fmt.Printf("Min1 - %f\n", L_opt[x.n - 1][x.m - 1]);
	// fmt.Printf("Min2 - %f\n", B_opt[x.n - 1][x.m - 1]);


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
				upperenv.Add(i + 1, x.Q[j], x.Q[j + 1], x.P[i + 1])
				h := queue.Left().(int)
				min := upperenv.FindMinimum(B_opt[h][j])
				if h < i {
					min = upperenv.FindMinimum(L_opt[i][j], B_opt[h][j])
				}

				first := queue.PopLeft()
				for queue.Size() > 1 && B_opt[queue.Left().(int)][j] <= min {
					h := queue.Left().(int)
					if !(h <= i) {
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
				queue := column_queues[i]
				upperenv := column_envelopes[i]
				for queue.Size() > 0 && L_opt[i][queue.Right().(int)] > L_opt[i][j] {
					queue.PopRight()
				}
				queue.PushRight(j)
				if queue.Size() == 1 {
					upperenv.Clear()
				}
				upperenv.Add(j + 1, x.P[i], x.P[i + 1], x.Q[j + 1])
				h := queue.Left().(int)
				min := upperenv.FindMinimum(L_opt[i][h])
				if h < i {
					min = upperenv.FindMinimum(B_opt[i][j], L_opt[i][h])
				}

				first := queue.PopLeft()
				for queue.Size() > 1 && L_opt[i][queue.Left().(int)] <= min {
					h := queue.Left().(int)
					if !(h <= j) {
						panic("h !<= j")
					}
					upperenv.RemoveUpto(h)
					min = upperenv.FindMinimum(L_opt[i][h])
					if h < j {
						min = upperenv.FindMinimum(B_opt[i][j], L_opt[i][h])
					}
					first = queue.PopLeft()
				}
				queue.PushLeft(first)

				B_opt[i][j + 1] = min;
				upperenv.TruncateLast();
			}
		}
		// fmt.Printf("Min1 - %f\n", L_opt[x.n - 1][x.m - 1]);
		// fmt.Printf("Min2 - %f\n", B_opt[x.n - 1][x.m - 1]);
	}
	distance := x.distance(x.P[x.n], x.Q[x.m])
	// fmt.Printf("Distance - %f\n", distance);
	// fmt.Printf("Min1 - %f\n", L_opt[x.n - 1][x.m - 1]);
	// fmt.Printf("Min2 - %f\n", B_opt[x.n - 1][x.m - 1]);

	return math.Max(distance, math.Min(L_opt[x.n - 1][x.m - 1], B_opt[x.n - 1][x.m - 1]))
}