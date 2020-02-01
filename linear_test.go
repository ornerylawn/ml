package ml

import (
	"linear"
	"math"
	"testing"
)

func ExpectFloat(expect, got float64, t *testing.T) {
	if math.Abs(got-expect) > 1e-9 {
		t.Errorf("expected %f but got %f", expect, got)
	}
}

func ExpectInt(expect, got int, t *testing.T) {
	if got != expect {
		t.Errorf("expected %d but got %d", expect, got)
	}
}

func TestLinearRegression(t *testing.T) {
	X := linear.NewArrayMatrix(2, 2)
	X.Set(0, 0, 1)
	X.Set(1, 0, 0)
	X.Set(0, 1, 1)
	X.Set(1, 1, 2)

	y := linear.NewArrayVector(2)
	y.Set(0, 6)
	y.Set(1, 0)

	theta_hat := LinearRegression(X, y)

	ExpectInt(2, theta_hat.Dimension(), t)
	ExpectFloat(6, theta_hat.Get(0), t)
	ExpectFloat(-3, theta_hat.Get(1), t)
}

func TestLinearRegressionNonSquare(t *testing.T) {
	X := linear.NewArrayMatrix(2, 3)
	X.Set(0, 0, 1)
	X.Set(1, 0, 0)
	X.Set(0, 1, 1)
	X.Set(1, 1, 2)
	X.Set(0, 2, -2)
	X.Set(1, 2, 1)

	y := linear.NewArrayVector(3)
	y.Set(0, 6)
	y.Set(1, 0)
	y.Set(2, -15)

	theta_hat := LinearRegression(X, y)

	ExpectInt(2, theta_hat.Dimension(), t)
	ExpectFloat(6, theta_hat.Get(0), t)
	ExpectFloat(-3, theta_hat.Get(1), t)
}

func TestLinearRegressionBoston(t *testing.T) {
	boston, err := LoadMatrixCSV("boston_housing.csv")
	if err != nil {
		t.Fatal("can't read boston dataset:", err)
	}
	ins, outs := boston.Shape()
	ExpectInt(ins, 14, t)
	ExpectInt(outs, 506, t)

	// TODO: do linear regression, match to tensorflow
}
