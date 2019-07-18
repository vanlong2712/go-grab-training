package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		error := eval(text)
		if error != nil {
			fmt.Println("err", error)
			fmt.Print("> ")
			continue
		}
		fmt.Print("> ")
	}
}

func parse(text string) (a float64, b float64, op string, err error) {
	expr := strings.Split(text, " ")
	if len(expr) != 3 {
		err = errors.New("Invalid length")
		return
	}
	op = expr[1]
	a, err = strconv.ParseFloat(expr[0], 10)
	if err != nil {
		return
	}
	b, err = strconv.ParseFloat(expr[2], 10)
	if err != nil {
		return
	}
	return
}

func eval(text string) error {
	a, b, op, err := parse(text)
	if err != nil {
		return err
	}
	switch op {
	case "+":
		fmt.Println(a, op, b, "=", a + b)
	case "-":
		fmt.Println(a, op, b, "=", a - b)
	case "/":
		if b == 0 && op == "/" {
			err = errors.New("Invalid operand")
			return
		}
		fmt.Println(a, op, b, "=", a / b)
	case "*":
		fmt.Println(a, op, b, "=", a * b)
	default:
		err = errors.New("OP does not match with +,-,*,/")
		return err
	}
	return nil
}