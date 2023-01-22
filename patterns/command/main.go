package main

import "fmt"

//  Команда — это паттерн проектирования, который превращает запросы в объекты, позволяя передавать их как аргументы при вызове методов, ставить запросы в очередь, логировать их, а также поддерживать отмену операций.
// Плюсы :Убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их непосредственно выполняют.
// Позволяет реализовать простую отмену и повтор операций.
// Позволяет реализовать отложенный запуск операций.
// Позволяет собирать сложные команды из простых.
// Реализует принцип открытости/закрытости.
// Минусы:Усложняет код программы из-за введения множества дополнительных классов.

// отправитель
// Кнопка от пультика
// с полем типа интерфейс Command

type Button struct {
	command Command
}

// метод нажатия на кнопку пультика press()
func (b *Button) press() {
	b.command.execute()
}

// Интерфейс команды
// //  команда задаёт общий интерфейс для конкретных
// // классов команд и содержит базовое поведение выполнение операции

type Command interface {
	execute()
}

// Интерфейс получателя с методами  on и off

type Device interface {
	on()
	off()
}

// Конкретная команда включения

type OnCommand struct {
	device Device
}

// выполнение этой команды
func (c *OnCommand) execute() {
	c.device.on()
}

// Конкретная команда выключения

type OffCommand struct {
	device Device
}

// выполнение выключения
func (c *OffCommand) execute() {
	c.device.off()
}

// Конкретный получатель

type Tv struct {
	isRunning bool
}

// конкретный получатель имплиментирует интерфейс Device
func (t *Tv) on() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}
func (t *Tv) off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}
func main() {

	tv := &Tv{}
	onCommand := &OnCommand{
		device: tv,
	}
	offCommand := &OffCommand{
		device: tv,
	}
	onButton := &Button{
		command: onCommand,
	}
	onButton.press()

	offButton := &Button{
		command: offCommand,
	}
	offButton.press()
}
