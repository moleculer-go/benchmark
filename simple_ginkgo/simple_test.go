package moleculer_test

import (
	"time"

	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/moleculer/broker"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Benchmarking", func() {
	var bkr *broker.ServiceBroker
	BeforeSuite(func() {
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
	})

	AfterSuite(func() {
		bkr.Stop()
	})
	nLoops := 15
	Measure("Simple call to simple action.", func(bench Benchmarker) {
		var total float64
		elapsed := bench.Time("time for loop", func() {
			total = doForOneSecond(func() {
				//bench.Time("math.add local call", func() {
				<-bkr.Call("math.add", map[string]interface{}{"a": 5, "b": 3})
				//sum(5, 3)
				//})
			})
		})
		bench.RecordValue("rps ? ", total/elapsed.Seconds())
	}, nLoops)
})

func sum(a, b int) int {
	return a + b
}

// doForOneSecond call the function in a loop for one second.
func doForOneSecond(do func()) float64 {
	var total int64 = 0
	start := time.Now()
	for {
		do()
		total++
		if (total%1000 == 0) && (time.Since(start) >= time.Second) {
			return float64(total)
		}
	}
}
