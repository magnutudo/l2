Что выведет программа? Объяснить вывод программы.

package main

type customError struct {
msg string
}

func (e *customError) Error() string {
return e.msg
}

func test() *customError {
{
// do something
}
return nil
}

func main() {
var err error
err = test()
if err != nil {
println("error")
return
}
println("ok")
}

Ответ:

Т.к. в функции test() мы возвращаем структуру которая наследуется от interface'а error возвращаемое значение nil
оборачивается в структуру customError и возвращаемое значение равняется customError(nil) что не является нулевым значением