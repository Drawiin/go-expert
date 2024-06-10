package calculator

func Calculate(a, b int, op string) int {
	if op == "+" {
		return a + b
	}
	if op == "-" {
		return a - b
	}
	if op == "*" {
		return a * b
	}
	if op == "/" && b != 0 && a != 0 {
		return a / b
	}
	if op == "/" && a == 0{
		return 0
	}
	return 0

}
