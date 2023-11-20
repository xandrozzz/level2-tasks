//Применимость:
//Паттерн visitor подходит для случаев, когда нужно провести операцию над группой однотипных объектов.
//В таком случае, логика операций будет находиться не в основных классах, а в классе visitor.
//
//Плюсы:
//Преимущество паттерна visitor заключается в том, что при изменении логики работы
//необходимо вносить изменения только в сами классы visitor, а не делать это во всех целевых классах.
//Еще одним преимуществом является то, интерфейс visitor подходит для всех целевых объектов,
//реализуя всю схожую логику в одном месте.
//
//Минусы:
//При изменении структуры одного класса приходится изменять все классы visitor, нацеленные на него.
//
//Примеры использования:
//Паттерн visitor применяется в случаях, когда нужно проводить однотипные операции над коллекцией объектов,
//реализующих один интерфейс, но для каждого класса требуется своя логика.
//Также одним из оснований применения паттерна visitor является сбор определенных данных об объектах
//в процессе выполнения операций над ними.
//

package main

import (
	"fmt"
	"strconv"
)

// Person - интерфейс, принимающий Visitor
type Person interface {
	GetSummary() string
	Accept(Visitor)
}

// Engineer - конкретный класс, реализующий интерфейс Person
type Engineer struct {
	Name                string
	Age                 int
	YearsWorked         int
	MachinesConstructed int
}

// Accept - метод для приема Visitor, реализуемый классом Engineer
func (e *Engineer) Accept(v Visitor) {
	v.VisitEngineer(e)
}

// GetSummary - метод для получения данных о человеке, реализуемый классом Engineer
func (e *Engineer) GetSummary() string {
	return "Engineer " + e.Name + ", Age: " + strconv.Itoa(e.Age)
}

// Artist - конкретный класс, реализующий интерфейс Person
type Artist struct {
	Name           string
	Age            int
	YearsWorked    int
	PaintingsDrawn int
}

// Accept - метод для приема Visitor, реализуемый классом Artist
func (a *Artist) Accept(v Visitor) {
	v.VisitArtist(a)
}

// GetSummary - метод для получения данных о человеке, реализуемый классом Artist
func (a *Artist) GetSummary() string {
	return "Artist" + a.Name + ", Age: " + strconv.Itoa(a.Age)
}

// Manager - конкретный класс, реализующий интерфейс Person
type Manager struct {
	Name           string
	Age            int
	YearsWorked    int
	EmployeesHired int
	EmployeesFired int
}

// Accept - метод для приема Visitor, реализуемый классом Manager
func (m *Manager) Accept(v Visitor) {
	v.VisitManager(m)
}

// GetSummary - метод для получения данных о человеке, реализуемый классом Manager
func (m *Manager) GetSummary() string {
	return "Manager" + m.Name + ", Age: " + strconv.Itoa(m.Age)
}

// Visitor - интерфейс, представляющий паттерн visitor
type Visitor interface {
	VisitEngineer(*Engineer) // метод для приема классом Engineer
	VisitArtist(*Artist)     // метод для приема классом Artist
	VisitManager(*Manager)   // метод для приема классом Manager
}

// SumCalculator - конкретный класс, реализующий интерфейс Visitor
type SumCalculator struct {
	AdditionalThingsDone int
}

// VisitEngineer - реализация метода VisitEngineer для класса SumCalculator
func (c *SumCalculator) VisitEngineer(e *Engineer) {
	fmt.Println(e.Name, "has constructed", e.MachinesConstructed+c.AdditionalThingsDone, "machines in total")
}

// VisitArtist - реализация метода VisitArtist для класса SumCalculator
func (c *SumCalculator) VisitArtist(a *Artist) {
	fmt.Println(a.Name, "has drawn", a.PaintingsDrawn+c.AdditionalThingsDone, "paintings in total")
}

// VisitManager - реализация метода VisitManager для класса SumCalculator
func (c *SumCalculator) VisitManager(m *Manager) {
	fmt.Println(m.Name, "has hired", m.EmployeesHired-m.EmployeesFired+c.AdditionalThingsDone, "employees in total")
}

// StatisticsCalculator - другой конкретный класс, реализующий интерфейс visitor
type StatisticsCalculator struct {
	AdditionalThingsDoneThisYear int
}

// VisitEngineer - реализация метода VisitEngineer для класса StatisticsCalculator
func (c *StatisticsCalculator) VisitEngineer(e *Engineer) {
	fmt.Println(e.Name, "constructs",
		float32(e.MachinesConstructed+c.AdditionalThingsDoneThisYear)/float32(e.YearsWorked), "machines per year")
}

// VisitArtist - реализация метода VisitArtist для класса StatisticsCalculator
func (c *StatisticsCalculator) VisitArtist(a *Artist) {
	fmt.Println(a.Name, "draws",
		float32(a.PaintingsDrawn+c.AdditionalThingsDoneThisYear)/float32(a.YearsWorked), "paintings per year")
}

// VisitManager - реализация метода VisitManager для класса StatisticsCalculator
func (c *StatisticsCalculator) VisitManager(m *Manager) {
	fmt.Println(m.Name, "hires",
		float32(m.EmployeesHired-m.EmployeesFired+c.AdditionalThingsDoneThisYear)/float32(m.YearsWorked), "employees per year")
}

func main() {

	// создание объектов классов, реализующих интерфейс Person
	engineer := &Engineer{Name: "Ivan", Age: 25, YearsWorked: 6, MachinesConstructed: 30}
	artist := &Artist{Name: "Petr", Age: 39, YearsWorked: 15, PaintingsDrawn: 20}
	manager := &Manager{Name: "Sasha", Age: 43, YearsWorked: 21, EmployeesHired: 62, EmployeesFired: 24}

	// объявление объекта класса, реализующего интерфейс Visitor
	sumCalculator := &SumCalculator{AdditionalThingsDone: 15}

	// вызов метода Accept у объектов интерфейса Person
	engineer.Accept(sumCalculator)
	artist.Accept(sumCalculator)
	manager.Accept(sumCalculator)

	fmt.Println()

	// объявление другого объекта класса, реализующего интерфейс Visitor
	statisticsCalculator := &StatisticsCalculator{AdditionalThingsDoneThisYear: 3}

	// вызов метода Accept у объектов интерфейса Person
	engineer.Accept(statisticsCalculator)
	artist.Accept(statisticsCalculator)
	manager.Accept(statisticsCalculator)

}
