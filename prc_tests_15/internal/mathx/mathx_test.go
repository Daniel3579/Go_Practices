package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSum_Table(t *testing.T) {
	cases := []struct {
		name string
		a, b int
		want int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative and positive", 10, -5, 5},
		{"both negative", -2, -3, -5},
		{"zero + zero", 0, 0, 0},
		{"zero + number", 0, 42, 42},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Sum(c.a, c.b)
			assert.Equal(t, c.want, got, "Sum(%d, %d) должна вернуть %d", c.a, c.b, c.want)
		})
	}
}

func TestSum_WithTestify(t *testing.T) {
	result := Sum(100, 200)
	require.Equal(t, 300, result)

	result = Sum(-50, 30)
	assert.Equal(t, -20, result)
}

func TestMultiply_Table(t *testing.T) {
	cases := []struct {
		name string
		a, b int
		want int
	}{
		{"positive", 5, 3, 15},
		{"by zero", 0, 999, 0},
		{"by one", 1, 42, 42},
		{"negative", -2, 3, -6},
		{"both negative", -2, -3, 6},
		{"large", 1000, 1000, 1000000},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Multiply(c.a, c.b)
			require.Equal(t, c.want, got, "Multiply(%d, %d) = %d", c.a, c.b, got)
		})
	}
}

func TestDivide_Success(t *testing.T) {
	result, err := Divide(10, 2)
	require.NoError(t, err, "должна быть успешная операция")
	assert.Equal(t, 5, result, "10 / 2 = 5")
}

func TestDivide_DivideByZero(t *testing.T) {
	result, err := Divide(10, 0)
	assert.Error(t, err, "ожидается ошибка при делении на ноль")
	assert.ErrorIs(t, err, ErrDivideByZero, "ошибка должна быть ErrDivideByZero")
	assert.Equal(t, 0, result, "результат должен быть 0 при ошибке")
}

func TestDivide_Table(t *testing.T) {
	cases := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		// Успешные случаи
		{"normal division", 20, 4, 5, false},
		{"divide by one", 42, 1, 42, false},
		{"divide zero", 0, 5, 0, false},
		{"negative dividend", -10, 2, -5, false},
		{"negative divisor", 10, -2, -5, false},
		{"both negative", -10, -2, 5, false},

		// Ошибочные случаи
		{"divide by zero", 10, 0, 0, true},
		{"zero divide zero", 0, 0, 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Divide(c.a, c.b)

			if c.wantErr {
				require.Error(t, err, "кейс %s: ожидается ошибка", c.name)
				assert.ErrorIs(t, err, ErrDivideByZero)
			} else {
				require.NoError(t, err, "кейс %s: не должно быть ошибки", c.name)
				assert.Equal(t, c.want, got, "кейс %s: неправильный результат", c.name)
			}
		})
	}
}

func TestAbs_EdgeCases(t *testing.T) {
	cases := []struct {
		name string
		x    int
		want int
	}{
		{"positive", 42, 42},
		{"negative", -42, 42},
		{"zero", 0, 0},
		{"large positive", 999999, 999999},
		{"large negative", -999999, 999999},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Abs(c.x)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestMax_Table(t *testing.T) {
	cases := []struct {
		name string
		a, b int
		want int
	}{
		{"a > b", 10, 5, 10},
		{"b > a", 5, 10, 10},
		{"a == b", 5, 5, 5},
		{"negative", -1, -5, -1},
		{"mixed", -5, 5, 5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Max(c.a, c.b)
			assert.Equal(t, c.want, got, "Max(%d, %d) = %d", c.a, c.b, got)
		})
	}
}

func TestMin_Table(t *testing.T) {
	cases := []struct {
		name string
		a, b int
		want int
	}{
		{"a < b", 5, 10, 5},
		{"b < a", 10, 5, 5},
		{"a == b", 5, 5, 5},
		{"negative", -10, -5, -10},
		{"mixed", -5, 5, -5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Min(c.a, c.b)
			assert.Equal(t, c.want, got, "Min(%d, %d) = %d", c.a, c.b, got)
		})
	}
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum(100, 200)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Multiply(50, 30)
	}
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Divide(1000, 2)
	}
}

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Abs(-12345)
	}
}

func BenchmarkMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Max(100, 50)
	}
}
