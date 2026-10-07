package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadEquation(path string) (float64, float64, error) {
	nums, err := numbers(path)
	if err != nil {
		return 0, 0, err
	}
	if len(nums) < 2 {
		return 0, 0, fmt.Errorf("input: %s: нужно 2 числа (eps, x0), найдено %d", path, len(nums))
	}
	return nums[0], nums[1], nil
}

func ReadSystem(path string) (int, float64, []float64, error) {
	nums, err := numbers(path)
	if err != nil {
		return 0, 0, nil, err
	}
	if len(nums) == 0 {
		return 0, 0, nil, fmt.Errorf("input: %s: пустой файл", path)
	}
	n := int(nums[0])
	if float64(n) != nums[0] || n <= 0 {
		return 0, 0, nil, fmt.Errorf("input: %s: первое число должно быть целой размерностью, найдено %v", path, nums[0])
	}
	need := 2 + n
	if len(nums) < need {
		return 0, 0, nil, fmt.Errorf("input: %s: нужно %d чисел (n, eps, x0), найдено %d", path, need, len(nums))
	}
	return n, nums[1], nums[2:need], nil
}

func numbers(path string) ([]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}
	defer f.Close()
	var nums []float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		for _, tok := range strings.Fields(sc.Text()) {
			v, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return nil, fmt.Errorf("input: %s: '%s': %w", path, tok, err)
			}
			nums = append(nums, v)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("input: %s: %w", path, err)
	}
	return nums, nil
}
