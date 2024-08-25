package utils

import "fmt"

func PrintList(list []string) {
	for _, item := range list {
		fmt.Println(item)
	}
}
