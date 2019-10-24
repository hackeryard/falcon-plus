package rpc

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

func main() {
	gatewayURL := "http://10.10.26.24:9091/"

	throughputGuage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "throughput",
		Help: "Throughput in Mbps",
	})
	throughputGuage.Set(800)

	if err := push.Collectors("throughput_job", push.HostnameGroupingKey(), gatewayURL, throughputGuage); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}
