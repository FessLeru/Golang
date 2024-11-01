package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)


var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}


func applyOperator(a, b float64, op rune) float64 {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b == 0 {
			panic("деление на ноль")
		}
		return a / b
	}
	return 0
}


func Calc(expression string) (float64, error) {
	var values []float64
	var ops []rune


	expression = strings.ReplaceAll(expression, " ", "")
	n := len(expression)

	for i := 0; i < n; i++ {
		char := rune(expression[i])

		if (char >= '0' && char <= '9') || char == '.' {
			start := i
			for i < n && ((expression[i] >= '0' && expression[i] <= '9') || expression[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, errors.New("ошибка в записи выражения")
			}
			values = append(values, num)
			i--
		} else if char == '(' {
			ops = append(ops, char)
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, errors.New("ошибка в записи выражения: недостаточно значений для операции")
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				values = append(values, applyOperator(val1, val2, op))
			}
			if len(ops) == 0 {
				return 0, errors.New("ошибка в записи выражения: отсутствует открывающая скобка")
			}
			ops = ops[:len(ops)-1]
		} else if precedence[char] > 0 {
			for len(ops) > 0 && precedence[ops[len(ops)-1]] >= precedence[char] {
				if len(values) < 2 {
					return 0, errors.New("ошибка в записи выражения: недостаточно значений для операции")
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				values = append(values, applyOperator(val1, val2, op))
			}
			ops = append(ops, char)
		} else if char != ' ' {
			return 0, errors.New("ошибка в записи выражения: некорректный символ")
		}
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, errors.New("ошибка в записи выражения: недостаточно значений для операции")
		}
		val2 := values[len(values)-1]
		values = values[:len(values)-1]

		val1 := values[len(values)-1]
		values = values[:len(values)-1]

		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		values = append(values, applyOperator(val1, val2, op))
	}

	if len(values) != 1 {
		return 0, errors.New("ошибка в записи выражения: неверное количество значений")
	}

	return values[0], nil
}

func main() {
	expression := "3 + 5 * (2 - 8) *"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}
}
