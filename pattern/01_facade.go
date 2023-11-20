//Применимость:
//Паттерн facade подходит для тех случаев, когда существует сложная система,
//которую нужно представить в упрощенном виде, или когда необходимо сделать внешний коммуникационный слой
//поверх существующей системы.
//Его цель - скрыть внутреннюю сложность за единым классом, который внешне выглядит просто.
//
//Плюсы:
//Инкапсулируя логику и взаимодействие подсистем за фасадом,
//можно изолировать их от остальной части системы и избежать нежелательных зависимостей.
//Обеспечивая единую точку доступа,
//можно скрыть сложность и детали подсистем и сделать их более понятными и удобными для использования.
//
//Минусы:
//Добавляя дополнительный уровень абстракции и взаимодействия,
//можно увеличить задержки и потребление ресурсов системы.
//Это может повлиять на скорость, эффективность и масштабируемость системы,
//особенно если фасад должен обрабатывать большое количество данных.
//Еще одним недостатком паттерна фасада является ограничение гибкости и настраиваемости.
//Скрывая детали и опции подсистем, можно ограничить контроль над функциональностью и поведением системы.
//
//Примеры использования:
//Паттерн facade чаще всего встречается при использовании сторонних библиотек.
//В таком случае внутри фасада происходит взаимодействие с функционалом библиотеки, а в основном коде - с фасадом.
//Это помогает упростить код, избавиться от ненужных зависимостей
//и определить только необходимые методы взаимодействия с библиотекой, не используя сложный и ненужный функционал.

package main

import "fmt"

// MiscDoer - класс внутри фасада
type MiscDoer struct{}

// DoSomething - метод для имитации деятельности
func (*MiscDoer) DoSomething() {
	fmt.Println("MiscDoer did something!")
}

// DoSomethingWithParameter - метод для имитации деятельности с параметрам
func (*MiscDoer) DoSomethingWithParameter(parameter int) {
	fmt.Printf("MiscDoer did something with parameter %d!\n", parameter)
}

// DoSomethingElse - альтернативный метод для имитации деятельности
func (*MiscDoer) DoSomethingElse() {
	fmt.Println("MiscDoer did something else!")
}

// MoreComplexDoer - другой класс внутри фасада
type MoreComplexDoer struct{}

// DoSomethingComplex - метод для имитации деятельности с несколькими параметрами
func (*MoreComplexDoer) DoSomethingComplex(parameter int, anotherParameter string) {
	fmt.Printf("MoreComplexDoer did something complex with parameters '%s' and %d!\n", anotherParameter, parameter)
}

// ExtraDoer - еще один класс внутри фасада
type ExtraDoer struct{}

// DoSomethingExtra - метод для имитации деятельности, не используемый в фасаде
func (*ExtraDoer) DoSomethingExtra(parameter int, anotherParameter int) string {
	data := fmt.Sprintf("ExtraDoer did something extra with parameters %d and %d, why?", parameter, anotherParameter)
	fmt.Println(data)
	return data
}

// SimpleFacade - класс фасада
type SimpleFacade struct {
	doer        *MiscDoer
	complexDoer *MoreComplexDoer
	extraDoer   *ExtraDoer
}

// NewSimpleFacade - конструктор фасада
func NewSimpleFacade() *SimpleFacade {
	return &SimpleFacade{
		doer:        &MiscDoer{},
		complexDoer: &MoreComplexDoer{},
		extraDoer:   &ExtraDoer{},
	}
}

// Start - метод фасада с несколькими методами объектов внутри фасада
func (c *SimpleFacade) Start() {
	c.doer.DoSomething()
	c.complexDoer.DoSomethingComplex(1, "aboba")
	c.doer.DoSomethingWithParameter(2)
	c.doer.DoSomethingElse()
}

func main() {
	facade := NewSimpleFacade() // создание фасада
	facade.Start()              // вызов метода старт для вызова методов объектов внутри фасада
}
