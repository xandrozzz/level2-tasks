//Применимость:
//Паттерн state стоит использовать для объектов, ведущих себя по-разному в зависимости от состояния,
//количество состояний очень велико и код для различных состояний части меняется.
//Также паттерн состояния рационально использовать для замены больших условных выражений внутри классов,
//опирающихся на поля объекта для определения его поведения.
//Кроме того, состояние можно применять для замены дублирующегося кода между похожими состояниями,
//и переходами между ними.
//Такое применение рационально в машине состояний, основанной на условиях.
//
//Плюсы:
//Возможна реализация кода для различных состояний в отдельных классах.
//Возможно добавление новых состояний без изменений в коде существующих классов или контекста.
//Упрощение кода контекста засчет избегания большого количества условий машины состояний.
//
//Минусы:
//Применение паттерна будет нерационально в случае, если машина состояний использует мало различных состояний,
//или состояние меняется нечасто.
//
//Примеры использования:
//Паттерн state можно использовать для онлайн-стриминга данных.
//В таком случае состояние будет меняться во время взаимодействия с конечным пользователем,
//а также во время загрузки данных.
//Кроме того, паттерн можно использовать для организации воркфлоу,
//отправляя необходимые уведомления на конкретные адреса в правильное время.
//Такой подход обеспечивает легкость изменения структуры потока.
//

package main

import (
	"fmt"
	"log"
)

// VendingMachine - класс контекст для State
type VendingMachine struct {
	hasItem       State
	itemRequested State
	hasMoney      State
	noItem        State

	currentState State

	itemCount int
	itemPrice int
}

// newVendingMachine - конструктор класса VendingMachine
func newVendingMachine(itemCount, itemPrice int) *VendingMachine {
	// создание изначального объекта VendingMachine
	v := &VendingMachine{
		itemCount: itemCount,
		itemPrice: itemPrice,
	}
	// привязка созданного объекта к состояниям
	hasItemState := &HasItemState{
		vendingMachine: v,
	}

	itemRequestedState := &ItemRequestedState{
		vendingMachine: v,
	}

	hasMoneyState := &HasMoneyState{
		vendingMachine: v,
	}

	noItemState := &NoItemState{
		vendingMachine: v,
	}

	v.setState(hasItemState) // установка изначального состояния
	// привязка состояний к объекту
	v.hasItem = hasItemState
	v.itemRequested = itemRequestedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState

	return v
}

// методы для смены состояния
func (v *VendingMachine) requestItem() error {
	return v.currentState.requestItem()
}

func (v *VendingMachine) addItem(count int) error {
	return v.currentState.addItem(count)
}

func (v *VendingMachine) insertMoney(money int) error {
	return v.currentState.insertMoney(money)
}

func (v *VendingMachine) dispenseItem() error {
	return v.currentState.dispenseItem()
}

// метод для прямой установки состояния
func (v *VendingMachine) setState(s State) {
	v.currentState = s
}

// метод для изменения поля itemCount
func (v *VendingMachine) incrementItemCount(count int) {
	fmt.Printf("Adding %d items\n", count)
	v.itemCount = v.itemCount + count // увеличение значения поля itemCount
}

// State - интерфейс, реализующий паттерн состояния
type State interface {
	addItem(int) error
	requestItem() error
	insertMoney(money int) error
	dispenseItem() error
}

// NoItemState - конкретный класс состояния
type NoItemState struct {
	vendingMachine *VendingMachine
}

// реализация методов интерфейса State конкретным классом NoItemState
func (i *NoItemState) requestItem() error {
	return fmt.Errorf("item out of stock") // ошибка из-за некорректной смены состояния
}

func (i *NoItemState) addItem(count int) error {
	i.vendingMachine.incrementItemCount(count) // увеличение значения itemCount с помощью метода incrementItemCount

	i.vendingMachine.setState(i.vendingMachine.hasItem) // смена состояния на hasItem

	return nil
}

func (i *NoItemState) insertMoney(money int) error {
	return fmt.Errorf("item out of stock") // ошибка из-за некорректной смены состояния
}

func (i *NoItemState) dispenseItem() error {
	return fmt.Errorf("item out of stock") // ошибка из-за некорректной смены состояния
}

// HasItemState - конкретный класс состояния
type HasItemState struct {
	vendingMachine *VendingMachine
}

// реализация методов интерфейса State конкретным классом HasItemState
func (i *HasItemState) requestItem() error {
	// проверка поля itemCount
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem) // смена состояния на noItem
		return fmt.Errorf("no item present")               // ошибка из-за некорректной смены состояния
	}
	fmt.Printf("Item requested\n")
	i.vendingMachine.setState(i.vendingMachine.itemRequested) // смена состояния на itemRequested
	return nil
}

func (i *HasItemState) addItem(count int) error {
	fmt.Printf("%d items added\n", count)

	i.vendingMachine.incrementItemCount(count) // увеличение значения itemCount с помощью метода incrementItemCount
	return nil
}

func (i *HasItemState) insertMoney(money int) error {
	return fmt.Errorf("please select item first") // ошибка из-за некорректной смены состояния
}

func (i *HasItemState) dispenseItem() error {
	return fmt.Errorf("please select item first") // ошибка из-за некорректной смены состояния
}

// ItemRequestedState - конкретный класс состояния
type ItemRequestedState struct {
	vendingMachine *VendingMachine
}

// реализация методов интерфейса State конкретным классом ItemRequestedState
func (i *ItemRequestedState) requestItem() error {
	return fmt.Errorf("item already requested") // ошибка из-за некорректной смены состояния
}

func (i *ItemRequestedState) addItem(count int) error {
	return fmt.Errorf("item dispense in progress") // ошибка из-за некорректной смены состояния
}

func (i *ItemRequestedState) insertMoney(money int) error {
	// проверка поля itemPrice
	if money < i.vendingMachine.itemPrice {
		return fmt.Errorf("inserted money is not enough. please insert %d", i.vendingMachine.itemPrice) // ошибка из-за некорректной смены состояния
	}
	fmt.Println("Money entered is sufficient")
	i.vendingMachine.setState(i.vendingMachine.hasMoney) // смена состояния на hasMoney
	return nil
}

func (i *ItemRequestedState) dispenseItem() error {
	return fmt.Errorf("please insert money first") // ошибка из-за некорректной смены состояния
}

// HasMoneyState - конкретный класс состояния
type HasMoneyState struct {
	vendingMachine *VendingMachine
}

// реализация методов интерфейса State конкретным классом HasMoneyState
func (i *HasMoneyState) requestItem() error {
	return fmt.Errorf("item dispense in progress") // ошибка из-за некорректной смены состояния
}

func (i *HasMoneyState) addItem(count int) error {
	return fmt.Errorf("item dispense in progress") // ошибка из-за некорректной смены состояния
}

func (i *HasMoneyState) insertMoney(money int) error {
	return fmt.Errorf("item out of stock") // ошибка из-за некорректной смены состояния
}

func (i *HasMoneyState) dispenseItem() error {
	fmt.Println("Dispensing item")

	i.vendingMachine.itemCount = i.vendingMachine.itemCount - 1 // уменьшение значения поля itemCount

	// проверка поля itemCount
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem) // смена значения на noItem
	} else {
		i.vendingMachine.setState(i.vendingMachine.hasItem) // смена значения на hasItem
	}

	return nil
}

func main() {

	vendingMachine := newVendingMachine(1, 10) // создание объекта VendingMachine через конструктор

	err := vendingMachine.requestItem() // изменение состояния через метод requestItem

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10) // изменение состояния через метод insertMoney

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem() // изменение состояния через метод dispenseItem

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.addItem(2) // изменение состояния через метод addItem

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	err = vendingMachine.requestItem() // изменение состояния через метод requestItem

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.insertMoney(10) // изменение состояния через метод insertMoney

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = vendingMachine.dispenseItem() // изменение состояния через метод dispenseItem

	// проверка на ошибку при некорректной смене состояния
	if err != nil {
		log.Fatalf(err.Error())
	}

}
