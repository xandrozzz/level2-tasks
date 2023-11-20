//Применимость:
//Паттерн factory method применяется в случаях, когда заранее неизвестно, с какими типами и зависимостями объектов
//должен работать основной код.
//Кроме того, паттерн применяется для отделения слоя создания объектов от их конкретных классов.
//Также фабричный метод полезен для экономии ресурсов системы
//засчет переиспользования созданных объектов вместо создания новых.
//
//Плюсы:
//Избегание связи между классом - создателем объектов и созданными объектами.
//Вероятность получить рекурсивный вызов.
//Перенос кода для создания объектов в одно место и упрощая его поддержку.
//Удобство добавления новых классов без больших изменений кода.
//
//Минусы:
//Значительное усложнение кода из-за введения множества новых подклассов для паттерна.
//
//Примеры использования:
//Паттерн factory method может использоваться для связи между независимыми классами.
//Например, при использовании приложением нескольких разных баз данных
//введение фабричного метода значительно упростит переключение между ними.
//Кроме того, возможно применение фабричного метода для создания расширений к основному коду.
//В таком случае, основная часть приложения не должна может быть соединена с любым расширением,
//вне зависимости от изменений в основном коде.
//

package main

import (
	"errors"
	"fmt"
)

// iPaper - интерфейс для разных конкретных классов, создаваемых фабричным методом
type iPaper interface {
	setHeight(int)
	getHeight() int
	setWidth(int)
	getWidth() int
}

// paper - класс, встраиваемый в конкретные классы фабричного метода
type paper struct {
	width  int
	height int
}

// реализация методов интерфейса iPaper классом paper
func (p *paper) setHeight(newHeight int) {
	p.height = newHeight
}

func (p *paper) getHeight() int {
	return p.height
}

func (p *paper) setWidth(newWidth int) {
	p.width = newWidth
}

func (p *paper) getWidth() int {
	return p.width
}

// a4Paper - конкретный класс для фабричного метода
type a4Paper struct {
	paper
}

// реализация методов интерфейса iPaper классом a4Paper
func (p *a4Paper) setHeight(newHeight int) {
	p.height = newHeight
}

func (p *a4Paper) getHeight() int {
	return p.height
}

func (p *a4Paper) setWidth(newWidth int) {
	p.width = newWidth
}

func (p *a4Paper) getWidth() int {
	return p.width
}

// newA4Paper - конструктор класса a4Paper
func newA4Paper() iPaper {
	return &a4Paper{
		paper: paper{
			width:  210,
			height: 297,
		},
	}
}

// a3Paper - конкретный класс для фабричного метода
type a3Paper struct {
	paper
}

// реализация методов интерфейса iPaper классом a3Paper
func (p *a3Paper) setHeight(newHeight int) {
	p.height = newHeight
}

func (p *a3Paper) getHeight() int {
	return p.height
}

func (p *a3Paper) setWidth(newWidth int) {
	p.width = newWidth
}

func (p *a3Paper) getWidth() int {
	return p.width
}

// newA3Paper - конструктор класса a3Paper
func newA3Paper() iPaper {
	return &a3Paper{
		paper: paper{
			width:  297,
			height: 420,
		},
	}
}

// getPaper - фабричный метод, возвращающий объекты интерфейса iPaper
func getPaper(paperType string) (iPaper, error) {
	switch paperType {
	case "a3":
		return newA3Paper(), nil // вызов конструктора в случае корректного аргумента типа
	case "a4":
		return newA4Paper(), nil // вызов конструктора в случае корректного аргумента типа
	default:
		return nil, errors.New(fmt.Sprintf("wrong type given: %v", paperType)) // возврат ошибки в случае некорректного аргумента типа
	}
}

func main() {
	// создание объектов двух разных классов через фабричный метод
	a4, _ := getPaper("a4")
	a3, _ := getPaper("a3")

	// вывод результатов
	fmt.Println(a4)
	fmt.Println(a3)
}
