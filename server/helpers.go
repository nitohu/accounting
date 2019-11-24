package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// HumanReadable takes a long number and converts
// it to a human readable string
func HumanReadable(num float64) string {
	signs := []string{"k", "m", "b", "t", "q", "Q", "s", "S"}
	number := fmt.Sprintf("%f", num)
	number = strings.Split(number, ".")[0]

	for i := 1; i < len(signs)+1; i++ {
		if len(number) > (i*3) && len(number) < (i*3+3) {
			// Calculate the indexes for cutting the number at the correct places
			firstCut := len(number) - (i * 3) - 1
			secondCut := firstCut + 1

			// Get the part of the number which will be removed
			del, _ := strconv.Atoi(number[firstCut:])

			exp, _ := strconv.ParseFloat(number[firstCut:secondCut], 64)
			// Calculate next smaller value of del
			smaller := exp * math.Pow(10, float64(i*3))
			// Calculate next bigger value of del
			bigger := (exp + 1) * math.Pow(10, float64(i*3))

			// Calculate the differences of both values to del
			delSmallerDiff := float64(del) - smaller
			delBiggerDiff := bigger - float64(del)

			firstCut = len(number) - i*3
			number = number[:firstCut]

			// If the difference of the bigger number with del is smaller
			// that means the number will be rounded up
			if delBiggerDiff <= delSmallerDiff {
				rounded, _ := strconv.Atoi(number)
				rounded++

				if rounded >= 1000 {
					rounded = rounded / 1000
					i++
				}
				number = fmt.Sprintf("%d", rounded)
			}

			// Add the sign
			number += signs[i-1]
			break
		}
	}

	return number
}
