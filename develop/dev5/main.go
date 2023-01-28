package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Количество строк
	Ccount := flag.Bool("c", false, "number of rows")
	// -n номер строки

	NFlag := flag.Bool("n", false, "number of row")
	NLine := os.Args[2]
	Aflag := flag.Bool("A", false, "print after")
	Bflag := flag.Bool("B", true, "print before")
	ANnumber := os.Args[3]
	if file := os.Args[1]; file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal("Unable to open file")
		}
		defer f.Close()
		if *Ccount {
			lines := bufio.NewScanner(f)
			LineCount := 0
			for lines.Scan() {
				LineCount++
			}
			fmt.Println(LineCount)
		}
		if *NFlag {
			lines := bufio.NewScanner(f)
			fmt.Println(NumOfRow(lines, NLine))
		}
		lines := []string{}
		counter := 0
		j := 0
		if *Aflag {
			mapOfstr := make(map[int]string)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				lines = append(lines, line)
				if strings.Contains(line, NLine) {
					mapOfstr[j] = line
					counter++
				}
				j++
			}
			a, err := strconv.Atoi(ANnumber)
			if err != nil {
				return
			}
			if a > 0 {
				for value, _ := range mapOfstr {
					fmt.Println("Row : ", value+1)
					if value <= j-a-1 {
						for k := value; k <= a+value-0; k++ {
							fmt.Println("another row ----> ", lines[k])

						}
						fmt.Println("--------------------")
					} else {
						for l := value; l < j; l++ {
							fmt.Println("another row ----> ", lines[l])
						}
					}
				}
			}
		}
		if *Bflag {
			mapOfstr := make(map[int]string)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				lines = append(lines, line)
				if strings.Contains(line, NLine) {
					mapOfstr[j] = line
					counter++
				}
				j++
			}
			b, err := strconv.Atoi(ANnumber)
			if err != nil {
				return
			}
			if b > 0 {
				for value, _ := range mapOfstr {
					if value > b {
						for i := value - b; i <= value-0; i++ {
							fmt.Println("another row -----> ", lines[i])
						}
						fmt.Println("--------------------")
					} else {
						for i := 0; i <= b-1; i++ {
							fmt.Println("another step -----> ", lines[i])
						}
						fmt.Println("--------------------")
					}
				}
			}
		}
	}
}
func NumOfRow(lines *bufio.Scanner, Nl string) (int, error) {
	line := 1
	for lines.Scan() {
		if strings.Contains(lines.Text(), Nl) {
			return line, nil
		}
		line++
	}
	if err := lines.Err(); err != nil {
		log.Fatal(err)
	}
	return line, nil
}
