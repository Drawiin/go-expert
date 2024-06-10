package calculator

import "testing"

func TestBatch(t *testing.T) {
	tests := []struct {
		a, b, expected int
		op             string
	}{
		{1, 2, 3, "+"},
		{2, 1, 1, "-"},
		{2, 3, 6, "*"},
		{6, 3, 2, "/"},
		{7, 0, 0, "/"},
	}
	for _, test := range tests {
		if output := Calculate(test.a, test.b, test.op); output != test.expected {
			t.Errorf("Test failed for test: %v, expected: %v, got: %v", test, test.expected, output)
		}
	}

}

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Calculate(1, 2, "+")
	}
}
