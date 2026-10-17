package demo

import (
	"fmt"

	"lab3/input"
	"lab3/lsm"
)

func LSM() {
	fmt.Println("=== 3.3 МНК: приближающие многочлены P1(x), P2(x), P3(x) ===")
	pts, err := input.ReadPoints("data/lsm.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("исходные точки (xi, yi):")
	for i := range pts.X {
		fmt.Printf("  (%8.4f, %12.6f)\n", pts.X[i], pts.Y[i])
	}

	n := float64(len(pts.X))
	best := ""
	bestSSE := 0.0
	for _, deg := range []int{1, 2, 3} {
		p, err := lsm.Fit(pts.X, pts.Y, deg)
		if err != nil {
			fmt.Println("ошибка:", err)
			return
		}
		sse := p.SSE(pts.X, pts.Y)
		fmt.Printf("P%d(x) = %s\n", deg, p.Form())
		fmt.Printf("сумма квадратов ошибок SSE(P%d) = %g, MSE = %g\n", deg, sse, sse/n)
		if deg == 1 || sse < bestSSE {
			best = fmt.Sprintf("P%d", deg)
			bestSSE = sse
		}
	}
	fmt.Printf("наилучший многочлен (мин. SSE): %s (SSE = %g)\n", best, bestSSE)
}
