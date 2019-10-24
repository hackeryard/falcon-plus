package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_completion_timestamp_seconds",
		Help: "The timestamp of the last completion of a DB backup, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_success_timestamp_seconds",
		Help: "The timestamp of the last successful completion of a DB backup.",
	})
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_duration_seconds",
		Help: "The duration of the last DB backup in seconds.",
	})
	records = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_records_processed",
		Help: "The number of records processed in the last DB backup.",
	})
)

func performBackup() (int, error) {
	// Perform the backup and return the number of backed up records and any
	// applicable error.
	// ...
	return 42, nil
}

func main() {
	// We use a registry here to benefit from the consistency checks that
	// happen during registration.
	registry := prometheus.NewRegistry()
	registry.MustRegister(completionTime, duration, records, successTime)
	// Note that successTime is not registered.

	// you can add Gatherer or Collector to a pusher
	pusher := push.New("http://10.10.26.24:9091", "vm_monitor").Gatherer(registry)

	start := time.Now()
	// metrics 1
	records.Set(float64(n))
	// Note that time.Since only uses a monotonic clock in Go1.9+.
	//metrics 2
	duration.Set(time.Since(start).Seconds())
	// metrics 3
	completionTime.SetToCurrentTime()

	successTime.SetToCurrentTime()

	// error 触发条件：任何pusher中添加的东西发送失败时
	if err := pusher.Push(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}
