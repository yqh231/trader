package main

import "fmt"

func main() {
	const (
		a = iota + 1
	)

	const (
		b = iota + 1
		c
	)
	switch a {
	case 2:
		fmt.Print(1)
	}

	fmt.Println(a, b, c)
}