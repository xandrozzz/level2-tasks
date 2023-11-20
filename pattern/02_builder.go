//Применимость:
//Паттерн builder подходит для случаев, когда нужно создавать объекты с множеством параметров.
//В таком случае обычный конструктор получится громоздким и неудобным
//из-за количества параметров и операций для создания объекта.
//Builder позволяет легко конфигурировать создаваемые объекты,
//а также создавать инкапсулированные объекты внутри основных.
//
//Плюсы:
//Параметры конструктора сокращаются и предоставляются в виде хорошо читаемых вызовов методов.
//Builder помогает минимизировать количество параметров в конструкторе, без необходимости передавать в конструктор
//пустые значения для необязательных параметров.
//Неизменяемые объекты могут быть созданы без использования сложной логики в процессе построения.
//
//Минусы:
//Количество строк кода для создания объекта значительно увеличивается, но код становится более читаемым.
//Требуется создание конкретного класса builder для каждого типа создаваемого объекта.
//
//Примеры использования:
//Паттерн builder чаще всего применяется при необходимости создания неизменяемых объектов, в случаях,
//когда обычный конструктор такого объекта получится слишком громоздким.
//Также паттерн builder применяется в случаях, когда возможно переиспользовать частично созданный объект для
//создания похожих на него объектов с немного различными свойствами.

package main

import "fmt"

// ActualBuilding - класс для создания через builder
type ActualBuilding struct {
	address          string
	postalCode       int
	residents        int
	builtAtYear      int
	isCommercial     bool
	ownerCompanyName string
}

// BuildingBuilder - основной класс builder для создания здания
type BuildingBuilder struct {
	building *ActualBuilding
}

// BuildingAddressBuilder - класс для создания адреса здания
type BuildingAddressBuilder struct {
	BuildingBuilder
}

// BuildingCompanyBuilder - класс для создания компании здания
type BuildingCompanyBuilder struct {
	BuildingBuilder
}

// NewBuildingBuilder - конструктор для класса BuildingBuilder
func NewBuildingBuilder() *BuildingBuilder {
	return &BuildingBuilder{building: &ActualBuilding{}}
}

// Located - метод, переходящий к builder'у адреса здания
func (b *BuildingBuilder) Located() *BuildingAddressBuilder {
	return &BuildingAddressBuilder{*b}
}

// Built - метод, переходящий к builder'у компании здания
func (b *BuildingBuilder) Built() *BuildingCompanyBuilder {
	return &BuildingCompanyBuilder{*b}
}

// At - метод, устанавливающий адрес здания
func (a *BuildingAddressBuilder) At(address string) *BuildingAddressBuilder {
	a.building.address = address
	return a
}

// WithPostalCode - метод, устанавливающий почтовый индекс здания
func (a *BuildingAddressBuilder) WithPostalCode(postalCode int) *BuildingAddressBuilder {
	a.building.postalCode = postalCode
	return a
}

// By - метод, устанавливающий название компании здания
func (c *BuildingCompanyBuilder) By(ownerCompanyName string) *BuildingCompanyBuilder {
	c.building.ownerCompanyName = ownerCompanyName
	return c
}

// AtYear - метод, устанавливающий год постройки здания
func (c *BuildingCompanyBuilder) AtYear(builtAtYear int) *BuildingCompanyBuilder {
	c.building.builtAtYear = builtAtYear
	return c
}

// IsCommercial - метод, определяющий, коммерческое здание или нет
func (b *BuildingBuilder) IsCommercial(isCommercial bool) *BuildingBuilder {
	b.building.isCommercial = isCommercial
	return b
}

// WithNumberOfResidents - метод, устанавливающий количество жильцов
func (b *BuildingBuilder) WithNumberOfResidents(residents int) *BuildingBuilder {
	b.building.residents = residents
	return b
}

// Build - метод, возвращающий созданный объект ActualBuilding
func (b *BuildingBuilder) Build() *ActualBuilding {
	return b.building
}

func main() {
	builder := NewBuildingBuilder() // объявление объекта builder
	// установка параметров создаваемого объекта через методы builder'a
	builder.Located().
		At("Moscow").
		WithPostalCode(424151).
		Built().
		By("PIK").
		AtYear(2020).
		WithNumberOfResidents(500).
		IsCommercial(true)

	// получение созданного объекта
	building := builder.Build()

	fmt.Println(building) // вывод результата

}
