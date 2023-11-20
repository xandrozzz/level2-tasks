//Применимость:
//Паттерн chain of responsibility используется в тех случаях, когда предполагается,
//что программа будет обрабатывать различные типы запросов различными способами,
//но точные типы запросов и их последовательность заранее неизвестны.
//Также этот паттерн можно использовать, когда необходимо запустить
//несколько обработчиков в определенном порядке или порядок должен меняться динамически.
//
//Плюсы:
//Возможно контролировать порядок вызова обработчиков.
//Возможно отделить классы, которые вызывают операции, от классов, которые выполняют, распределяя ответственность.
//Возможно легко добавить новые обработчики без изменения большей части существующего кода.
//
//Минусы:
//Некоторые запросы могут остаться необработанными.
//Вероятность получить рекурсивный вызов.
//
//Примеры использования:
//Паттерн chain of responsibility может использоваться при логировании,
//цепочка логгеров в таком случае будет использоваться для записи в разные выходы
//на основании уровня или категории сообщения.
//Также паттерн может применяться при валидации данных,
//в таком случае разные валидаторы будут проверять данные по разным правилам.
//Еще один вариант применения паттерна - авторизация.
//Возможно создать цепочку хэндлеров авторизации, проверяющих доступ к чтению или изменению различных частей данных.
//

package main

import "fmt"

// RestaurantSpace - интерфейс для хэндлеров цепочки
type RestaurantSpace interface {
	Handle(*HungryPerson)
	SetNext(RestaurantSpace)
}

// OrderTerminal - конкретный класс хэндлера цепочки
type OrderTerminal struct {
	next RestaurantSpace
}

// Handle - реализация метода Handle для класса OrderTerminal
func (r *OrderTerminal) Handle(p *HungryPerson) {
	// проверка поля foodOrdered
	if p.foodOrdered {
		fmt.Println("Food is already ordered")
		r.next.Handle(p) // вызов метода Handle у следующего хэндлера
		return
	}
	fmt.Println("Order terminal was used to order food")
	p.foodOrdered = true // установка поля foodOrdered на true
	r.next.Handle(p)     // вызов метода Handle у следующего хэндлера
}

// SetNext - реализация метода SetNext для класса OrderTerminal
func (r *OrderTerminal) SetNext(next RestaurantSpace) {
	r.next = next
}

// Reception - конкретный класс хэндлера цепочки
type Reception struct {
	next RestaurantSpace
}

// Handle - реализация метода Handle для класса Reception
func (d *Reception) Handle(p *HungryPerson) {
	// проверка поля foodOrdered
	if p.foodReceived {
		fmt.Println("Food was already picked up ")
		d.next.Handle(p) // вызов метода Handle у следующего хэндлера
		return
	}
	fmt.Println("Food was picked up from reception")
	p.foodReceived = true // установка поля foodReceived на true
	d.next.Handle(p)      // вызов метода Handle у следующего хэндлера
}

// SetNext - реализация метода SetNext для класса Reception
func (d *Reception) SetNext(next RestaurantSpace) {
	d.next = next
}

// Table - конкретный класс хэндлера цепочки
type Table struct {
	next RestaurantSpace
}

// Handle - реализация метода Handle для класса Table
func (m *Table) Handle(p *HungryPerson) {
	// проверка поля foodEaten
	if p.foodEaten {
		fmt.Println("Food was already eaten")
		m.next.Handle(p) // вызов метода Handle у следующего хэндлера
		return
	}
	fmt.Println("Food was eaten at the table")
	p.foodEaten = true // установка поля foodReceived на true
	m.next.Handle(p)   // вызов метода Handle у следующего хэндлера
}

// SetNext - реализация метода SetNext для класса Table
func (m *Table) SetNext(next RestaurantSpace) {
	m.next = next
}

// TrashBin - конкретный класс хэндлера цепочки
type TrashBin struct {
	next RestaurantSpace
}

// Handle - реализация метода Handle для класса TrashBin
func (c *TrashBin) Handle(p *HungryPerson) {
	// проверка поля trashDisposed
	if p.trashDisposed {
		fmt.Println("The trash was already disposed of")
	}
	fmt.Println("The trash was thrown into the trash bin")
}

// SetNext - реализация метода SetNext для класса Table
func (c *TrashBin) SetNext(next RestaurantSpace) {
	c.next = next
}

// HungryPerson - целевой класс для цепочки
type HungryPerson struct {
	name string
	// поля, проверяющие прохождение элементов цепочки
	foodOrdered   bool
	foodReceived  bool
	foodEaten     bool
	trashDisposed bool
}

func main() {

	trashBin := &TrashBin{} // создание хэндлера цепочки TrashBin

	table := &Table{}       // создание хэндлера цепочки Table
	table.SetNext(trashBin) // установка следующего хэндлера цепочки

	reception := &Reception{} // создание хэндлера цепочки Reception
	reception.SetNext(table)  // установка следующего хэндлера цепочки

	orderTerminal := &OrderTerminal{} // создание хэндлера цепочки OrderTerminal
	orderTerminal.SetNext(reception)  // установка следующего хэндлера цепочки

	hungryPerson := &HungryPerson{name: "Misha"} // создание объекта для цепочки
	reception.Handle(hungryPerson)               // вызов метода Handle для обработки объекта цепочкой

}
