package main

import (
	"fmt"
	"sushi/sushi"
	"sync"
	"time"
)

const (
	foodMax = 100
)

type worker interface {
	Start()
}

func main() {
	dining := make(chan *sushi.Sushi, 20)
	food := make(chan struct{}, foodMax)
	done := make(chan struct{})
	exit := make(chan int)

	for i := 0; i < foodMax; i++ {
		food <- struct{}{}
	}

	cooks := []worker{
		sushi.NewCook(1, 0, 0, dining, food, done),
		sushi.NewCook(2, 0, 0, dining, food, done),
		sushi.NewCook(3, 0, 0, dining, food, done),
		sushi.NewCook(4, 0, 0, dining, food, done),
		sushi.NewCook(5, 0, 0, dining, food, done),
	}
	go run(cooks, exit, 1)

	customers := []worker{
		sushi.NewCustomer(1, 0, 0, dining, done),
		sushi.NewCustomer(2, 0, 0, dining, done),
		sushi.NewCustomer(3, 0, 0, dining, done),
		sushi.NewCustomer(4, 0, 0, dining, done),
	}

	go run(customers, exit, 2)

	go func() {
		s := <-exit
		switch s {
		case 1:
			fmt.Println("**厨师下班了**")
		case 2:
			fmt.Println("**顾客走光了**")
		}
		close(done)
	}()

	<-done
	time.Sleep(time.Second)
	fmt.Println("营业结束")
}

func run(workers []worker, exit chan int, singal int) {
	wg := &sync.WaitGroup{}
	wg.Add(len(workers))
	for _, c := range workers {
		go func(c worker) {
			c.Start()
			wg.Done()
		}(c)
	}
	wg.Wait()
	exit <- singal
}
