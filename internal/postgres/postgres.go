package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
)

// Config represents the configuration for connecting to the Postgres database.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// Client represents the Postgres client.
type Client struct {
	Conn *pgx.Conn
}

type ClientPool struct {
	mu      *sync.Mutex
	conns   []*Client
	maxConn int
	channel chan interface{}
}

// Create a Postgres client
func NewClient(config Config) (*Client, error) {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: conn,
	}, nil
}

// Close closes the connection to the Postgres database.
func (c *Client) Close() error {
	return c.Conn.Close(context.Background())
}

func NewClientPool(config Config, maxConn int) (*ClientPool, error) {
	var mu = sync.Mutex{}
	cPool := &ClientPool{
		mu:      &mu,
		conns:   make([]*Client, 0, maxConn),
		maxConn: maxConn,
		channel: make(chan interface{}, maxConn),
	}

	for range maxConn {
		client, err := NewClient(config)
		if err != nil {
			return nil, err
		}

		cPool.conns = append(cPool.conns, client)
		cPool.channel <- nil
	}

	return cPool, nil
}

func (cPool *ClientPool) Acquire() *Client {
	<-cPool.channel

	cPool.mu.Lock()
	c := cPool.conns[0]
	cPool.conns = cPool.conns[1:]
	cPool.mu.Unlock()

	return c
}

func (cPool *ClientPool) Release(c *Client) {
	cPool.mu.Lock()
	cPool.conns = append(cPool.conns, c)
	cPool.channel <- nil
	cPool.mu.Unlock()
}

func (cPool *ClientPool) Close() error {
	for i := range cPool.conns {
		err := cPool.conns[i].Conn.Close(context.Background())
		if err != nil {
			fmt.Printf("Unable to close the connection")
			return err
		}
	}

	return nil
}
