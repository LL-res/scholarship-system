package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	var strA, strB string
	fmt.Scan(&strA, &strB)
	str := make([]uint8, 0, 2*n)
	for i := 0; i < n; i++ {
		str = append(str, strA[i], strB[i])
	}
	fmt.Println(string(str))
}
