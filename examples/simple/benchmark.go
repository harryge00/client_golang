package main

import (
	"log"
	"os"
	"strconv"
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
var count = 0
func query(client prometheus.Client) time.Duration {
	log.Printf("query %d", count)
	count++
	pRange := prometheus.Range{
		Start: time.Now().Add(- 60 * time.Second),
		End:   time.Now(),
		Step:  10 * time.Second,
	}
	var wg sync.WaitGroup

	queries := []string {
		"sum(rate(container_cpu_usage_seconds_total{id=\"/\"}[10m]))",
		"sum(rate(container_network_receive_errors_total{id=\"/\"}[10m])) ",
		"sum(container_spec_memory_limit_bytes{id=\"/\"}) ",
		"sum(rate(container_network_transmit_errors_total{id=\"/\"}[10m])) ",
		"sum(rate(container_network_receive_bytes_total{id=\"/\"}[10m]))",
		"sum(rate(container_network_transmit_bytes_total{id=\"/\"}[10m])) ",
		"sum(container_fs_usage_bytes{id=\"/\",device!~\"/dev/mapper.*\"}) ",
		"sum(container_fs_limit_bytes{id=\"/\",device!~\"/dev/mapper.*\"}) ",
		"sum(machine_cpu_cores{}) ",
		"sum(container_memory_usage_bytes{id=\"/\"}) ",
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
				log.Println(err)
				panic(err)
			}
			matrix, ok := values.(model.Matrix)
			if !ok {
				panic("typer error")
			}
			// log.Println(i, time.Since(beginTime), q)
			var res MetricData 
			for _, m := range matrix {
				// log.Println("metric:", m.Metric, len(m.Values))
				res.TimeStamps = make([]time.Time, len(m.Values))
				res.MetricValues = make([]float64, len(m.Values))
				for i, val := range m.Values {
					res.TimeStamps[i] = val.Timestamp.Time()
					res.MetricValues[i] = float64(val.Value)
				}
			}
			log.Println(i, time.Since(beginTime), q, res.MetricValues)
			wg.Done()
		} (i, queries[i])
	}
	wg.Wait()
	return time.Since(totalBeginTime)
}

func main() {
	config := prometheus.Config{
		Address: "http://localhost:9090",
	}
	client, err := prometheus.New(config)
	if err != nil {
		panic(err)
	}
	num, _ := strconv.Atoi(os.Getenv("Num"))
	arr := make([]time.Duration, num)
	var sum time.Duration
	for i := 0; i < num; i++ {
		res := query(client)
		arr[i] = res
		log.Println(i, "cost", res)
		sum += res
	}
	log.Printf("arr:%v, sum:%v, avg:%v", arr, sum, sum/time.Duration(num))
}
/*
[58.979754ms 49.768949ms 1.745104735s 237.006725ms 40.112603ms 335.634895ms 559.074211ms 45.498542ms 625.465484ms 336.705404ms]
[843.204348ms 611.274505ms 621.44978ms 418.945613ms 34.44283ms 50.751781ms 36.919145ms 42.633877ms 792.171588ms 83.543277ms]
[173.942615ms 49.289646ms 50.678209ms 43.04082ms 62.90851ms 59.77217ms 826.748453ms 72.621498ms 644.75162ms 1.537237969s]
[435.865712ms 40.499205ms 29.519696ms 624.762537ms 614.590206ms 619.703962ms 93.427435ms 51.342981ms 58.494447ms 50.941069ms 661.676994ms 163.64658ms 38.825446ms 35.330977ms 28.726485ms 351.333495ms 512.02552ms 39.651302ms 55.007725ms 621.663388ms 1.846205872s 534.780508ms 36.620826ms 651.809455ms 258.409261ms 40.981923ms 623.620909ms 919.261299ms 470.492406ms 53.436962ms 39.725158ms 38.316671ms 855.946341ms 52.122704ms 27.153101ms 616.540288ms 240.591621ms 64.74024ms 614.333944ms 613.494755ms 726.456359ms 25.217852ms 25.99735ms 44.483449ms 36.185502ms 31.170533ms 32.689478ms 615.162133ms 215.66575ms 28.482002ms 49.878074ms 326.408605ms 1.223806493s 611.779704ms 681.20799ms 59.358791ms 37.67763ms 49.548092ms 53.550198ms 31.975805ms 318.851703ms 454.10314ms 50.082328ms 29.367082ms 35.473314ms 349.804837ms 537.441907ms 383.257328ms 920.322957ms 616.208886ms 117.742677ms 43.898056ms 45.39063ms 49.755626ms 358.598823ms 460.093812ms 2.918169854s 148.871023ms 46.558622ms 103.62331ms 1.236446066s 923.848951ms 623.84398ms 74.396735ms 136.120249ms 392.110908ms 448.937562ms 97.814322ms 58.519045ms 624.381665ms 211.299424ms 52.747977ms 350.288806ms 554.577078ms 1.197904851s 56.375622ms 740.205303ms 206.169105ms 622.943077ms 272.429589ms]

*/