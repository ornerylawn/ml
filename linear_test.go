package ml

import (
	"math"
	"testing"

	"github.com/ornerylawn/linear"
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

	X := linear.Slice(boston, 0, ins-1, 0, outs)
	y := linear.VectorFromColumn(linear.Slice(boston, ins-1, ins, 0, outs))

	theta_hat := LinearRegression(X, y)

	// tensorflow
	ExpectFloat(-9.28965170e-02, theta_hat.Get(0), t)
	ExpectFloat(4.87149552e-02, theta_hat.Get(1), t)
	ExpectFloat(-4.05997958e-03, theta_hat.Get(2), t)
	ExpectFloat(2.85399882e+00, theta_hat.Get(3), t)
	ExpectFloat(-2.86843637e+00, theta_hat.Get(4), t)
	ExpectFloat(5.92814778e+00, theta_hat.Get(5), t)
	ExpectFloat(-7.26933458e-03, theta_hat.Get(6), t)
	ExpectFloat(-9.68514157e-01, theta_hat.Get(7), t)
	ExpectFloat(1.71151128e-01, theta_hat.Get(8), t)
	ExpectFloat(-9.39621540e-03, theta_hat.Get(9), t)
	ExpectFloat(-3.92190926e-01, theta_hat.Get(10), t)
	ExpectFloat(1.49056102e-02, theta_hat.Get(11), t)
	ExpectFloat(-4.16304471e-01, theta_hat.Get(12), t)
}
