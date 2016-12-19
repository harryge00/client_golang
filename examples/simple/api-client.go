package main

import (
	"log"
	"time"
	"golang.org/x/net/context"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type MetricData struct {
	TimeStamps   []time.Time
	MetricValues []float64
}

func main() {
	config := prometheus.Config{
		Address: "http://localhost:9090",
	}
	client, err := prometheus.New(config)
	if err != nil {
		panic(err)
	}
	pRange := prometheus.Range{
		Start: time.Now().Add(- 24 * time.Hour),
		End:   time.Now(),
		Step:  300 * time.Second,
	}
	// q := prometheus.NewQueryAPI(client)
	values, err := prometheus.
		NewQueryAPI(client).
		QueryRange(context.TODO(), "sum(container_memory_usage_bytes{kubernetes_namespace=\"allen\",kubernetes_pod_name=\"allen-test-3232540786-vfsyg\"})", pRange)
	if err != nil {
		panic(err)
	}
	matrix, ok := values.(model.Matrix)
	if !ok {
		panic("typer error")
	}
	var res MetricData 
	for _, m := range matrix {
		log.Println("metric:", m.Metric, len(m.Values))
		res.TimeStamps = make([]time.Time, len(m.Values))
		res.MetricValues = make([]float64, len(m.Values))
		for i, val := range m.Values {
			res.TimeStamps[i] = val.Timestamp.Time()
			res.MetricValues[i] = float64(val.Value)
		}
	}
	log.Println(res)
	values2, err := prometheus.
		NewQueryAPI(client).
		QueryRange(context.TODO(), "sum(container_spec_cpu_quota{kubernetes_pod_name=\"monitoring-server-v1.0.0-8d0vs\"})", pRange)
	if err != nil {
		panic(err)
	}
	matrix2, ok := values2.(model.Matrix)
	if !ok {
		panic("typer error")
	}
	var res2 MetricData 
	for _, m := range matrix2 {
		log.Println("metric2:", m.Metric, len(m.Values))
		res2.TimeStamps = make([]time.Time, len(m.Values))
		res2.MetricValues = make([]float64, len(m.Values))
		for i, val := range m.Values {
			res2.TimeStamps[i] = val.Timestamp.Time()
			res2.MetricValues[i] = float64(val.Value)
		}
	}

	log.Println(res2)
}