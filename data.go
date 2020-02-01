package ml

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"linear"
	"math/rand"
	"os"
	"strconv"
)

func LoadMatrixCSV(filename string) (linear.Matrix, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	recordLength := -1
	values := []float64{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if recordLength == -1 {
			recordLength = len(record)
			continue // first record is header
		} else if len(record) != recordLength {
			return nil, fmt.Errorf("expected record of length %d but found length %d", recordLength, len(record))
		}
		for _, v := range record {
			x, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			values = append(values, x)
		}
	}
	if len(values) == 0 {
		return nil, errors.New("no values found in csv")
	}
	ins := recordLength
	outs := len(values) / recordLength
	A := linear.NewArrayMatrix(ins, outs)
	for o := 0; o < outs; o++ {
		for i := 0; i < ins; i++ {
			A.Set(i, o, values[o*ins+i])
		}
	}
	return A, nil
}

func randomInt(lo, hi int) int {
	return rand.Intn(hi-lo) + lo
}

func GenerateRandomOrder(d int) []int {
	order := make([]int, d)
	for i := 0; i < d; i++ {
		order[i] = i
	}
	for i := 0; i < d; i++ {
		j := randomInt(i, d)
		order[i], order[j] = order[j], order[i]
	}
	return order
}

func swapRows(A linear.Matrix, a, b, ins int) {
	for i := 0; i < ins; i++ {
		tmp := A.Get(i, a)
		A.Set(i, a, A.Get(i, b))
		A.Set(i, b, tmp)
	}
}

func ShuffleRows(A linear.Matrix, order []int) {
	ins, outs := A.Shape()
	location := make([]int, outs)
	for o := 0; o < outs; o++ {
		location[o] = o
	}
	for dst, src := range order {
		srcLoc := location[src]
		swapRows(A, dst, srcLoc, ins)
		location[dst], location[src] = srcLoc, dst
	}
}

func SplitRows(A linear.Matrix, hi int) (top, bottom linear.Matrix) {
	ins, outs := A.Shape()
	return linear.Slice(A, 0, ins, 0, hi), linear.Slice(A, 0, ins, hi, outs)
}
