package sushi

import (
	"fmt"
	"math/rand"
	time "time"
)

const (
	ER1 = "吃饱了"
	ER2 = "没有寿司了"
	ER3 = "营业结束"
)

var (
	ereasones = []string{ER1, ER2, ER3}
)

type Customer struct {
	id          int
	max         int
	spend       int //单位 ms
	diningTable <-chan *Sushi
	done        chan struct{}
	count       int
}

func NewCustomer(id int, max int, spend int, diningTable <-chan *Sushi, done chan struct{}) *Customer {
	rand.Seed(time.Now().UnixNano())
	if max == 0 {
		max = rand.Intn(30) + 5
	}
	if spend == 0 {
		spend = rand.Intn(20) + 10
	}
	return &Customer{id: id, max: max, count: 0, spend: spend, diningTable: diningTable, done: done}
}

func (c *Customer) Start() {
	fmt.Printf("顾客 %02d:%+v 开始就餐\n", c.id, c)
	for {
		select {
		case s, ok := <-c.diningTable:
			if !ok {
				c.exit(2)
				return
			}
			if ok = c.eat(s); !ok {
				c.exit(1)
				return
			}
		case <-c.done:
			c.exit(3)
			return
		}
	}
}

func (c *Customer) eat(s *Sushi) bool {
	if c.count >= c.max {
		return false
	}
	time.Sleep(time.Millisecond * time.Duration(c.spend))
	fmt.Printf("顾客 %02d 吃掉了一个寿司 %+v\n", c.id, *s)
	c.count++
	return true
}

func (c *Customer) exit(r int) {
	fmt.Printf("顾客 %02d 离席，原因：%s,一共吃了 %d/%d 个寿司\n", c.id, ereasones[r-1], c.count, c.max)
}
