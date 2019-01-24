package main

import (
	"fmt"

	"github.com/lflxp/studyNeo4j/pkg"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func helloWorld(uri, username, password string) (string, error) {
	var (
		err      error
		driver   neo4j.Driver
		session  neo4j.Session
		result   neo4j.Result
		greeting interface{}
	)

	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return "", err
	}
	defer session.Close()

	greeting, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func main() {
	// rs, err := helloWorld("bolt://localhost:7687", "root", "system")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(rs)
	// rs, err := pkg.ReadTran("match (b:Person{name:'张三'})-[x:BELONG]->(a)<-[r:分配]-(c) return a,c,b,r,x", nil)
	// rs, err := pkg.ReadTran("match (a)-[r:BELONG]->(b) return a,b,r", nil)
	// rs, err := pkg.ReadTran("match (a)-[r:ACTED_IN]->(b) return a.name,b,r", nil)
	// rs, err := pkg.ReadTran("match (a:Person) return a.name", nil)
	// rs, err := pkg.ReadTran("create (a:Person{name:'123'})", nil)
	rs, err := pkg.ReadTran("match (a:Person{name:'Tom Hanks'})-[r:ACTED_IN]->(b) with a,b,r match (b)<-[x:DIRECTED]-(m) return a,r,b,x,m", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(rs)
}
