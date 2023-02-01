package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	files     []fs.FileInfo
	childrens []Children
	errFiles  error
	comm      string
	str       string
	arg       string
	pwd       string
)

type Children struct {
	Pid       int
	TimeStart time.Time
}

func main() {
	Shell()
}

func Shell() {
	pwd = "/Users/moratherest/Projects/wb_l2"
	files, _ = ioutil.ReadDir(pwd)

	for {
		fmt.Print(pwd + "> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		str = scanner.Text()
		commands := strings.Split(str, " | ")

		for _, val := range commands {
			arr := make([]string, 2)
			arr = strings.Split(val, " ")
			comm = arr[0]
			if len(arr) > 1 {
				arg = arr[1]
			}

			switch comm {
			case "cd":

				//возвращение назад
				if arg == "-" {
					split := strings.Split(pwd, "/")
					pwd = ""
					for i := 1; i < len(split)-1; i++ {
						pwd += "/" + split[i]
					}
					arg = pwd
				} else { // идем вперед
					for _, v := range files {
						if v.Name() == arg {
							arg = pwd + "/" + arg // добавляем к директории переданный путь
							pwd = arg
							break
						}
					}
				}

				files, errFiles = ioutil.ReadDir(arg)
				if errFiles != nil {
					fmt.Println("wrong dir")
					continue
				}
				for _, v := range files {
					fmt.Println(v.Name())
				}
				break
			case "pwd":
				fmt.Println(pwd)
				break
			case "echo":
				fmt.Println(arg)
				break
			case "kill":
				pid, errPid := strconv.Atoi(arg)
				if errPid != nil {
					fmt.Println("wrong argument")
					continue
				}

				pgid, errPgid := syscall.Getpgid(pid) //получаем pgid который присваивается процессу при его создании
				if errPgid != nil {
					fmt.Println("failed to kill chldrn:", errPgid)
					continue
				}

				for i, v := range childrens { // удаляем процесс из слайса
					if v.Pid == pid {
						childrens = append(childrens[:i], childrens[i+1:]...)
						break
					}
				}

				if errKill := syscall.Kill(-pgid, syscall.SIGKILL); errKill != nil { // убиваем процесс
					fmt.Println("failed to kill fork with pid: ", arg, errKill)
					continue
				}
				fmt.Println("killed pid:", pid)
				break
			case "ps":
				for _, v := range childrens {
					//находим разницу в датах и выводим
					day := time.Now().Day() - v.TimeStart.Day()
					hour := time.Now().Hour() - v.TimeStart.Hour()
					minute := time.Now().Minute() - v.TimeStart.Minute()

					fmt.Printf("Time: %v:%v:%v	Pid:%v\n", day, hour, minute, v.Pid)
				}
			case "fork":
				childPID, _ := syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr{
					Sys: &syscall.SysProcAttr{
						Setpgid: true,
					},
				})
				child := Children{
					Pid:       childPID,
					TimeStart: time.Now(),
				}

				childrens = append(childrens, child)
				fmt.Println("fork run with id:", childPID)
				break
			case "quit":
				for _, v := range childrens {
					syscall.Kill(v.Pid, syscall.SIGKILL)
				}
				os.Exit(0)
			}
		}
	}
}
