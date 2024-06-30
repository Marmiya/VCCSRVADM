package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// GetCPULoad queries the Prometheus API for CPU utilization of the node exporter at the specified IP address.
func GetCPULoad(ip string) (float64, error) {
	// Create a new Prometheus API client
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("http://%s:9090", ip),
	})
	if err != nil {
		return 0, fmt.Errorf("error creating client: %v", err)
	}

	// Create a new API v1 client
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query the CPU utilization
	// query := `100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
	result, warnings, err := v1api.Query(ctx, "up", time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		return 0, fmt.Errorf("error querying Prometheus: %v", err)
	}
	if len(warnings) > 0 {
		log.Printf("warnings: %v", warnings)
	}

	// Parse the result
	// if result.Type() == model.ValVector {
	// 	vector := result.(model.Vector)
	// 	if len(vector) > 0 {
	// 		return float64(vector[0].Value), nil
	// 	}
	// }

	fmt.Printf("Result:\n%v\n", result)

	return 0, fmt.Errorf("no data returned")
}