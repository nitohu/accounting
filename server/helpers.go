package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GenerateSessionKey generates a unique key for each session
func GenerateSessionKey() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charset += "0123456789?.,-_*!:;#+"

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	key := make([]byte, 32)

	for i := range key {
		key[i] = charset[seed.Intn(len(charset))]
	}

	return string(key)
}

// HumanReadable takes a long number and converts
// it to a human readable string
func HumanReadable(num float64, digits int) string {
	signs := []string{"k", "m", "b", "t", "q", "Q", "s", "S"}
	number := fmt.Sprintf("%f", num)
	number = strings.Split(number, ".")[0]
	negative := false

	if number[0] == '-' {
		negative = true
		number = number[1:]
	}

	for i := 1; i < len(signs)+1; i++ {
		if len(number) > (i*3) && len(number) <= (i*3+3) {
			// Calculate the indexes for cutting the number at the correct places
			firstCut := len(number) - (i * 3) - 1 + digits
			secondCut := firstCut + 1

			decimals := ""
			var roundingDevice *string
			if digits > 0 {
				// If digits is greater than zero the decimals variable should be used for rounding
				roundingDevice = &decimals
				fc := firstCut + 1
				decimals = number[(fc - digits):fc]
			} else {
				// If no digits are appended to the rounded number, the number itself should be used for rounding
				roundingDevice = &number
			}

			//
			// Rounding
			//
			// Get the part of the number which will be removed
			del, _ := strconv.Atoi(number[firstCut:])
			// Get digit for rounding
			exp, _ := strconv.ParseFloat(number[firstCut:secondCut], 64)
			// Calculate next smaller value of del
			smaller := exp * math.Pow(10, float64(i*(3-digits)))
			// Calculate next bigger value of del
			bigger := (exp + 1) * math.Pow(10, float64(i*(3-digits)))

			// Calculate the differences of both values to del
			delSmallerDiff := float64(del) - smaller
			delBiggerDiff := bigger - float64(del)

			firstCut = len(number) - i*3
			number = number[:firstCut]

			// If the difference of the bigger number with del is smaller
			// that means the number will be rounded up
			if delBiggerDiff <= delSmallerDiff {
				rounded, _ := strconv.Atoi(*roundingDevice)
				rounded++

				if rounded >= 1000 {
					rounded = rounded / 1000
					i++
				}
				*roundingDevice = fmt.Sprintf("%d", rounded)
			}

			if decimals != "" {
				number += "." + decimals
			}

			// Add the sign
			number += signs[i-1]
			break
		}
	}

	if negative {
		number = "-" + number
	}

	return number
}
