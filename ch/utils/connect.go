package utils

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
)

func Connect(addrs []string, database string, user string, password string) (clickhouse.Conn, error) {
	connect, err := clickhouse.Open(&clickhouse.Options{
		Addr: addrs,
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		//Debug: true,
	})
	if err == nil {
		err = connect.Ping(context.Background())
	}
	return connect, err
}
