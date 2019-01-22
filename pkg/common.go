package pkg

import (
	"github.com/jmcvetta/neoism"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	driver  neo4j.Driver
	session neo4j.Session
	result  neo4j.Result
	err     error
)

func NewConn(path string) (*neoism.Database, error) {
	db, err := neoism.Connect(path)
	return db, err
}
