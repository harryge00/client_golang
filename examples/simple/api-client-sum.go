package main

import (
	"log"
	"fmt"
	"time"
	"golang.org/x/net/context"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type MetricData struct {
	TimeStamps   []time.Time
	MetricValues []float64
}

type S float64
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
	go func() {
		values, err := prometheus.
			NewQueryAPI(client).
			QueryRange(context.Background(), "sum(rate(container_cpu_usage_seconds_total{kubernetes_pod_name=\"nginx-pc-535175895-sxhet\"}[10m]))", pRange)
		if err != nil {
			panic(err)
		}
		matrix, ok := values.(model.Matrix)
		if !ok {
			panic("typer error")
		}
		log.Println(matrix)
		
	} ()
	go func() {
		values2, err := prometheus.
			NewQueryAPI(client).
			QueryRange(context.Background(), "sum(rate(container_cpu_usage_seconds_total{kubernetes_pod_name=\"nginx-pc-535175895-sxhet\"}[10m]))+sum(rate(container_cpu_usage_seconds_total{kubernetes_pod_name=\"nginx-pc-535175895-4nv51\"}[5m]))", pRange)
		if err != nil {
			panic(err)
		}
		matrix2, ok := values2.(model.Matrix)
		if !ok {
			panic("typer error")
		}
		log.Println(matrix2)
	} ()
	var abc string
	fmt.Scanln(&abc)
}