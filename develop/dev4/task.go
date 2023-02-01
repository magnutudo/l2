package main

import (
	"fmt"
	"strings"
)

func main() {
	arr := []string{"слИток", "ПЯтка", "столик", "стилок"}
	//arr2 := []string{"
	mapa := GetSets(arr)
	fmt.Println(mapa)
}

func GetSets(dictionary []string) map[string][]string {
	mapa := make(map[string][]string)
	for _, v := range dictionary {
		v = strings.ToLower(v)
		sum := GetCharSum(v)

		var key = v
		for k, _ := range mapa {
			kSum := GetCharSum(k)
			if sum == kSum {
				key = k
				break
			}
		}
		// TODO: при создании ключа не добавлять сам ключ
		mapa[key] = appendWithSort(mapa[key], v)
	}

	return mapa
}

func appendWithSort(arr []string, value string) []string {
	var newArr []string
	var index int
	var stop bool

	if len(arr) == 0 {
		newArr = append(newArr, value)
		return newArr
	}

	for i, v := range arr {
		if value == v {
			return arr
		}

		for j, _ := range value {
			if j >= len(v) || value[j] > v[j] {
				index = i + 1
				stop = true
				break
			} else if value[j] == v[j] {
				continue
			} else if value[j] < v[j] {
				index = i
				stop = true
				break
			}
		}
		if stop {
			break
		}
	}

	newArr = append(newArr, arr[:index]...)
	newArr = append(newArr, value)
	newArr = append(newArr, arr[index:]...)
	return newArr
}

func GetCharSum(str string) int {
	var sum int32
	for _, v := range str {
		sum += v
	}

	return int(sum)
}
