package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/moleculer/broker"
)

func add(a, b int) int {
	return a + b
}

func main() {

	var bkr *broker.ServiceBroker

	bkr = broker.New(&moleculer.Config{})
	bkr.AddService(moleculer.Service{
		Name: "math",
		Actions: []moleculer.Action{{
			Name: "add",
			Handler: func(context moleculer.Context, params moleculer.Payload) interface{} {
				return params.Get("a").Int() + params.Get("b").Int()
			}}},
	})
	bkr.Start()
	params := map[string]interface{}{"a": 5, "b": 3}
	runAndPrint(func(b *bench) {
		total, rps, avg, seconds := doForOneSecond(func() {
			<-bkr.Call("math.add", params)
			//add(5, 3)
		})
		b.save("rps", rps)
		b.save("total", total)
		b.save("avg", avg)
		b.save("seconds", seconds)
	}, 10)

	bkr.Stop()
}

type bench struct {
	values *sync.Map
}

func (b *bench) save(name string, value float64) {
	var list []float64
	temp, exists := b.values.Load(name)
	if !exists {
		list = []float64{value}
		b.values.Store(name, list)
		return
	}
	b.values.Store(name, append(temp.([]float64), value))
}

func stats(list []float64) (float64, float64, float64, float64) {
	var min, max, avg, total float64 = -1, 0, 0, 0
	for _, item := range list {
		if item < min || min == -1 {
			min = item
		}
		if item > max {
			max = item
		}
		total += item
	}
	avg = total / float64(len(list))
	return min, max, avg, total
}

func (b *bench) print() {
	b.values.Range(func(key, value interface{}) bool {
		min, max, avg, total := stats(value.([]float64))
		fmt.Println("\n **** Stats for [ ", key, " ] **** ")
		fmt.Println("\n avg: ", avg)
		fmt.Println("\n min: ", min)
		fmt.Println("\n max: ", max)
		fmt.Println("\n total: ", total)
		return true
	})
}

func runAndPrint(fn func(*bench), times int) {
	b := &bench{&sync.Map{}}
	for i := 0; i < times; i++ {
		fn(b)
	}
	b.print()
}

// doForOneSecond call the function in a loop for one second.
func doForOneSecond(do func()) (float64, float64, float64, float64) {
	var total int64
	var rps, avg, seconds float64
	start := time.Now()
	for {
		do()
		total++
		if total%1000 == 0 {
			seconds = time.Since(start).Seconds()
			if seconds > 1 {
				rps = float64(total) / seconds
				avg = seconds / float64(total)
				return float64(total), rps, avg, seconds
			}
		}
	}
	return float64(total), rps, avg, seconds
}
