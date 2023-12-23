package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/procfs"
)

var (
	ProcFS, _ = procfs.NewDefaultFS()
	memFree   = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mem_free",
			Help: "Free memory",
		},
	)
	processesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "processesCount",
			Help: "Total count of running processes",
		},
	)
)

func init() {
	prometheus.MustRegister(memFree)
	prometheus.MustRegister(processesCount)
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go collectMemInfo()
	go collectProcStat()

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":8080"}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stopChan

	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped.")
}

func collectMemInfo() {
	for {
		memStat, err := ProcFS.Meminfo()
		if err != nil {
			log.Printf("Failed to retrieve memory stats: %s\n", err)
			continue
		}
		memFree.Set(float64(*memStat.MemFree))

		time.Sleep(10 * time.Second)
	}
}

func collectProcStat() {
	for {
		processes, err := ProcFS.AllProcs()
		if err != nil {
			log.Printf("Failed to retrieve processes: %s\n", err)
			continue
		}

		processesCount.Set(float64(len(processes)))

		time.Sleep(10 * time.Second)
	}
}
