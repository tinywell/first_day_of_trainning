package sushi

import (
	"fmt"
	"math/rand"
	"time"
)

type Sushi struct {
	id   string
	cook int
}

const (
	MR1 = "原材料不足"
	MR2 = "厨师任务完成"
	MR3 = "营业结束"
)

var (
	reasones = []string{MR1, MR2, MR3}
)

type Cook struct {
	id          int
	max         int
	spend       int //单位 ms
	diningTable chan<- *Sushi
	food        chan struct{}
	done        chan struct{}
	count       int
}

func NewCook(id int, max int, spend int, diningTable chan<- *Sushi, food chan struct{}, done chan struct{}) *Cook {
	rand.Seed(time.Now().UnixNano())
	if max == 0 {
		max = rand.Intn(90) + 10
	}
	if spend == 0 {
		spend = rand.Intn(45) + 5
	}
	return &Cook{id: id, max: max, spend: spend, diningTable: diningTable, food: food, count: 0, done: done}
}

func (c *Cook) Start() {
	fmt.Printf("厨师 %02d:%+v 开始工作\n", c.id, c)
	for {
		select {
		case _, ok := <-c.food:
			if !ok {
				c.exit(1)
				return
			}
			if ss, ok := c.makeSushi(); ok {
				c.diningTable <- ss
			} else {
				c.exit(2)
				return
			}
		case <-c.done:
			c.exit(3)
			return
		default:
			c.exit(1)
			return
		}
	}
}

func (c *Cook) makeSushi() (*Sushi, bool) {
	if c.count > c.max {
		return nil, false
	}
	time.Sleep(time.Millisecond * time.Duration(c.spend))
	s := &Sushi{id: fmt.Sprintf("%02d%04d", c.id, c.count), cook: c.id}
	fmt.Printf("厨师 %02d 制作了一个寿司 %+v\n", c.id, *s)
	c.count++
	return s, true
}

func (c *Cook) exit(r int) {
	fmt.Printf("厨师 %d 结束制作，原因：%s，共制作了 %d/%d 个寿司\n", c.id, reasones[r-1], c.count, c.max)
}
