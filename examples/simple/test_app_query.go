package main

import (
	"log"
	"time"
	"sync"
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
		Step:  10 * time.Second,
	}
	// q := prometheus.NewQueryAPI(client)
	
	var wg sync.WaitGroup
	queries := []string {
		"sum(container_spec_cpu_quota{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(rate(container_cpu_usage_seconds_total{kubernetes_namespace=\"test1\"}[10m]) * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_namespace=\"test1\", kubernetes_dp_name=\"test1001\"})",
		"sum(container_memory_working_set_bytes{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(container_memory_usage_bytes{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(container_fs_usage_bytes{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(container_fs_limit_bytes{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(rate(container_network_receive_errors_total{kubernetes_namespace=\"test1\"}[10m]) * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(rate(container_network_transmit_errors_total{kubernetes_namespace=\"test1\"}[10m]) * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(rate(container_network_receive_bytes_total{kubernetes_namespace=\"test1\"}[10m]) * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(rate(container_network_transmit_bytes_total{kubernetes_namespace=\"test1\"}[10m]) * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"test1001\", kubernetes_namespace=\"test1\"}) ",
		"sum(container_memory_working_set_bytes{kubernetes_namespace=\"test1\"} * on (io_kubernetes_pod_uid) group_left(kubernetes_dp_name) kubernetes_resource_hierarchy {kubernetes_dp_name=\"redis\", kubernetes_namespace=\"micosvc\"}) ",
	}
	wg.Add(len(queries))
	totalBeginTime := time.Now()
	for i := range queries {
		go func(i int, q string) {
			beginTime := time.Now()
			values, err := prometheus.
			NewQueryAPI(client).
			QueryRange(context.TODO(), q, pRange)
			if err != nil {
				panic(err)
			}
			matrix, ok := values.(model.Matrix)
			if !ok {
				panic("typer error")
			}
			log.Println(i, time.Since(beginTime), q)
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
			log.Println(i, time.Since(beginTime), "processed", len(res.TimeStamps), len(res.MetricValues))
			wg.Done()
		} (i, queries[i])
	}
	wg.Wait()
	log.Println("total", time.Since(totalBeginTime))
}