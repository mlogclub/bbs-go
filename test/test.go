package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5}
	for _, v := range arr[0:len(arr)] {
		fmt.Println(v)
	}
}
