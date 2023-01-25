Что выведет программа? Объяснить вывод программы.

package main

func main() {
ch := make(chan int)
go func() {
for i := 0; i < 10; i++ {
ch <- i
}
}()

for n := range ch {
println(n)
}
}
Ответ:

программа выведет от 0 до 9, после чего произойдет ошибка "fatal error: all goroutines are asleep - deadlock!"
данная ошибка вызвана из-за цикла на месте range ch, range читает данные из канала пока он не будет закрыт,
а так-как в горутине не произошло закрытия канала, range будет слушать канал бесконечно.
