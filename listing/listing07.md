Что выведет программа? Объяснить вывод программы.

    package main
    
    import (
        "fmt"
        "math/rand"
        "time"
    )
    
    func asChan(vs ...int) <-chan int {
        c := make(chan int)
        go func() {
            for _, v := range vs {
                c <- v
                time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
            }
            close(c)
        }()
        return c
    }
    
    func merge(a, b <-chan int) <-chan int {
        c := make(chan int)
        go func() {
            for {
                select {
                case v := <-a:
                    c <- v
                case v := <-b:
                    c <- v
                }
            }
        }()
        return c
    }
    
    func main() {
        a := asChan(1, 3, 5, 7)
        b := asChan(2, 4 ,6, 8)
        c := merge(a, b)
        for v := range c {
            fmt.Println(v)
        }
    }

Ответ:

Программа выведет

    Цифры от 1 до 8 в случайном порядке, а потом бесконечные нули

В данном случае вывод данных не останавливается также по причине взаимодействия range
с каналом, который никогда не закрывается, в данном случае - каналом c.
Даже когда отправка данных из функции asChan завершается, и выражение select внутри for в функции merge
блокирует поток, возвращаемый функцией merge канал не закрывается.
Из-за этого range бесконечно итерируется по пустому каналу, выводя нули.