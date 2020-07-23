package db

import (
	"github.com/jackc/pgx"
)

var NextPoolCon func() *pgx.Conn //pointer to a pgx conneciton struct

type connCache chan *pgx.Conn

//
func (c connCache) populate() {
	for {
		c <- connectOrDie()
	}
}

func StartConnCache(size int) func() *pgx.Conn {
	var c connCache = make(chan *pgx.Conn, size)
	go c.populate()
	return func() *pgx.Conn { return <-c }
}
