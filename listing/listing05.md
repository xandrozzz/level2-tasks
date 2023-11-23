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

Программа выведет

    error

Строка error выводится потому, что переменная err не равна nil.
Тип customError реализует интерфейс error, реализуя его единственный метод Error().
Тем не менее, из функции test на выходе получается интерфейс с типом *customError,
хотя обычная функция, возвращающая ошибку, вернула бы интерфейс с типом nil.
Поэтому значение переменной err и nil - неодинаковы.
    