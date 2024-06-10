# postgres-conn-pool-limits
* Compare query execution time with and without connection pools.

## Prerequisites
* Docker
* Docker Compose
* Golang 1.22 or later

## Setup
* Clone the repository
* Run `make benchmark` to start the benchmark

#Infrastructure
* Machine: MacBook M2 Pro

## Results
* The results are as follows for different number of queries and connection pool limits:
* Sample Query: `SELECT pg_sleep(0.01);`

| Number Of <br> Queries | Number of Conn <br> In Pool | Exec Time Without <br> Connection Pool | Exec Time With <br> Connection Pool |
|:----------------------:|:---------------------------:|:--------------------------------------:|:-----------------------------------:|
|           10           |              1              |                 252 ms                 |               133 ms                |
|          100           |             10              |                1.33 sec                |               2.3 sec               |
|          1000          |             10              |                13.5 sec                |              22.8 sec               |
|          1000          |             10              |                13.5 sec                |              22.8 sec               |
|          1000          |             50              |                13.4 sec                |              22.8 sec               |
|          1000          |             100             |    Error(Connection Limit Reached)     |              22.8 sec               |


## Conclusion
* Connection pools are useful constructs to speed up the execution of queries.
* After a certain limit, the dividends of adding additional connections to the pool diminishes.
* There is no generic formula to determine the ideal number of connection pools. It depends on a multitude of factors such as:
  * Number of queries 
  * Complexity/Kind of queries 
  * Network latency
  * Underlying hardware
  * Database load, configuration and capacity
  * etc.