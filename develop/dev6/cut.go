package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	f             string
	d             string
	text          string
	numOfMismatch int
	s             bool
)

func main() {
	flag.StringVar(&f, "f", "0-0", "\"fields\" - выбрать поля (колонки)")
	flag.StringVar(&d, "d", "\t", "\"delimiter\" - использовать другой разделитель")
	flag.BoolVar(&s, "s", false, "\"separated\" - только строки с разделителем")
	flag.Parse()

	arr, errOpen := OpenFile(flag.Arg(0))
	if errOpen != nil {
		text = flag.Arg(0)
		arr = strings.Split(text, "\n")
	}

	split := Split(arr, d)

	left, right := getBorders(f)
	fields := GetFields(split, left, right)
	Print(fields)
}

func OpenFile(path string) ([]string, error) {
	var arr []string

	file, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}

	return arr, nil
}

// Получаем границу выделения
func getBorders(borders string) (int, int) {
	b := strings.Split(borders, "-")
	if len(b) > 2 {
		code, _ := fmt.Fprintln(os.Stderr, fmt.Errorf("wrong arguments"))
		os.Exit(code)
	}
	left, err := strconv.Atoi(b[0])
	if err != nil {
		code, _ := fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
	right, err := strconv.Atoi(b[1])
	if err != nil {
		code, _ := fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
	return left, right
}

// Split разбиваем массив строк и получаем массив с массивом строк
func Split(arr []string, delimiter string) [][]string {

	var arrSplit = make([][]string, len(arr))
	for i, v := range arr {
		a := strings.Split(v, delimiter)
		if len(a) == 1 && s {
			numOfMismatch++
		}
		arrSplit[i] = a
	}

	return arrSplit
}

// GetFields получаем поля
func GetFields(arr [][]string, leftBorder, rightBorder int) [][]string {
	var fields = make([][]string, len(arr)-numOfMismatch) // создаем массив без учета строк где нет разделителя
	var index int
	for _, v := range arr {
		if len(v) < rightBorder {
			if s {
				continue
			} else { // если выводим не вошедшие в условие
				fields[index] = append(fields[index], v[0:]...)
				index++
				continue
			}
		}
		if len(v) > 1 {
			for k := leftBorder; k <= rightBorder; k++ {
				fields[index] = append(fields[index], v[k])
			}
			index++
		}

	}

	return fields
}

func Print(arr [][]string) {
	for _, v := range arr {
		for _, val := range v {
			fmt.Print(val + " ")
		}
		fmt.Println()
	}
}
