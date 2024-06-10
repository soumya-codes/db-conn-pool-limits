package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/soumya-codes/postgres-conn-limits/internal/postgres"
)

const (
	query = "SELECT pg_sleep(0.01);"
)

func main() {
	config := postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "admin_password",
		Database: "test_conn_pool_db",
	}

	// Benchmark the time taken to execute queries with and without connection pooling
	BenchmarkNonPooledConnection(config, 10)
	BenchmarkPooledConnection(config, 10, 10)
}

func BenchmarkNonPooledConnection(config postgres.Config, noOfQueries int) {

	startTime := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < noOfQueries; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := postgres.NewClient(config)
			if err != nil {
				log.Panic("Error connecting to Postgres:", err)
			}

			defer client.Close()

			_, err = client.Conn.Exec(context.Background(), query)
			if err != nil {
				log.Panic("Error executing query:", err)
			}
		}()

		wg.Wait()
	}

	fmt.Println("Time taken to execute queries without connection pooling:", time.Since(startTime))
}

func GetConnectionPool(config postgres.Config, maxConn int) (*postgres.ClientPool, error) {
	pool, err := postgres.NewClientPool(config, maxConn)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return nil, err
	}

	return pool, nil
}

func BenchmarkPooledConnection(config postgres.Config, noOfQueries int, maxConn int) {
	pool, err := GetConnectionPool(config, maxConn)
	if err != nil {
		log.Panic("Error creating connection pool:", err)
	}

	startTime := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < noOfQueries; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := pool.Acquire()

			_, err = client.Conn.Exec(context.Background(), query)
			if err != nil {
				log.Panic("Error executing query:", err)
			}

			pool.Release(client)
		}()

		wg.Wait()
	}

	fmt.Println("Time taken to execute queries with connection pooling:", time.Since(startTime))
}
