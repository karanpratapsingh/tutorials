package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var size = 8
var total = 1_000_000
var subject = "subject"
var data = []byte("benchmark")

func main() {
	var external []opts.BarData
	var embedded []opts.BarData

	for i := 0; i < size; i++ {
		external = append(external, opts.BarData{Value: Benchmark()})
		embedded = append(embedded, opts.BarData{Value: BenchmarkEmbedded()})
	}

	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)

	bar.
		SetXAxis([]int{1, 2, 3, 4, 5, 6, 7, 8}).
		AddSeries("External", external).
		AddSeries("Embedded", embedded)

	file, err := os.Create("benchmark/results.html")

	if err != nil {
		panic(err)
	}

	bar.Render(file)
	fmt.Println("Results have been saved to results.html")
}

func Benchmark() (duration int64) {
	wait := make(chan bool)
	nc := GetClient(nats.DefaultURL)

	start := time.Now()

	processed := 0

	nc.Subscribe(subject, func(msg *nats.Msg) {
		processed += 1

		if processed >= total {
			duration = time.Since(start).Milliseconds()
			wait <- false
		}
	})

	for i := 0; i < total; i++ {
		nc.Publish(subject, data)
	}

	<-wait
	return duration
}

func BenchmarkEmbedded() (duration int64) {
	opts := &server.Options{
		Port: 4223,
	}
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	nc := GetClient(ns.ClientURL())

	start := time.Now()

	processed := 0

	nc.Subscribe(subject, func(msg *nats.Msg) {
		processed += 1

		if processed >= total {
			duration = time.Since(start).Milliseconds()
			ns.Shutdown()
		}
	})

	for i := 0; i < total; i++ {
		nc.Publish(subject, data)
	}

	ns.WaitForShutdown()
	return duration
}

func GetClient(url string) *nats.Conn {
	nc, err := nats.Connect(url)

	if err != nil {
		panic(err)
	}
	return nc
}
