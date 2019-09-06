package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Base int

const (
	operand Base = iota + 1
	number
)

func main() {
	fmt.Println(operand, number)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		str, err := eval(text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(str)
		fmt.Print("> ")
	}
}

func eval(text string) (string, error) {
	var currentType Base
	var prevType Base
	var result float64
	var str string
	skip := false
	a := strings.Fields(text)
	if getIndexOfArray(a[0], []string{"*", "/"}) > -1 {
		err := errors.New(`The first element is invalid, it must be either number or (*) or (/)`)
		return str, err
	}
	if getIndexOfArray(a[len(a)-1], []string{"+", "-", "*", "/"}) > -1 {
		err := errors.New(`The last element is invalid, it must be a number`)
		return str, err
	}
	for i, value := range a {
		if skip {
			skip = false
			continue
		}
		c, err := strconv.ParseFloat(value, 10)
		if err != nil {
			currentType = operand
			if getIndexOfArray(value, []string{"+", "-", "*", "/"}) < 0 {
				err := fmt.Errorf(`character at index %d does not match`, i)
				return str, err
			}
		} else {
			currentType = number
		}

		if prevType != 0 && prevType == currentType {
			var mType string
			if currentType == operand {
				mType = "number"
			} else {
				mType = "operand"
			}
			err := fmt.Errorf(`character at index %d must be %s`, i, mType)
			return str, err
		}
		if currentType == operand {
			d, err := strconv.ParseFloat(a[i+1], 10)
			if err != nil {
				err := fmt.Errorf(`character at index %d must be %s`, i + 1, "number")
				return str, err
			}
			result, err = calc(result, d, value)
			skip = true
			currentType = number
		} else {
			result = c
		}
		prevType = currentType
	}
	str = strings.Join(a, " ") + " = " + fmt.Sprintf("%g", result)
	return str, nil
}

func getIndexOfArray(text string, slice []string) int {
	for i, v := range slice {
		if v == text {
			return i
		}
	}
	return -1
}

func calc(a float64, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "/":
		if b == 0 {
			err := errors.New("can not divide zero")
			return 0, err
		}
		return a / b, nil
	case "*":
		return a * b, nil
	default:
		return 0, nil
	}
}
