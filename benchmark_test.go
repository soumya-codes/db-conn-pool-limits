package main

import (
	"testing"

	"github.com/soumya-codes/postgres-conn-limits/internal/postgres"
)

// Table driven tests for Benchmarking Query Performance with and without connection pooling
func TestAdd(t *testing.T) {
	// Define test cases
	tests := []struct {
		name                 string
		noOfQueries          int
		maxConnectionsInPool int
	}{
		{"Number Of Queries: 10, Max Connection In Connection Pool: 1", 10, 1},
		{"Number Of Queries: 10, Max Connection In Connection Pool: 10", 10, 10},
		{"Number Of Queries: 100, Max Connection In Connection Pool: 10", 100, 10},
		{"Number Of Queries: 1000, Max Connection In Connection Pool: 10", 1000, 10},
		{"Number Of Queries: 1000, Max Connection In Connection Pool: 50", 1000, 50},
		{"Number Of Queries: 1000, Max Connection In Connection Pool: 100", 1000, 100},
	}

	config := postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "admin_password",
		Database: "test_conn_pool_db",
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the benchmark tests
			BenchmarkPooledConnection(config, tt.noOfQueries, tt.maxConnectionsInPool)
			BenchmarkNonPooledConnection(config, tt.noOfQueries)
		})
	}
}
