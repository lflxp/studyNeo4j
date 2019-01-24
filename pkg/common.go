package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	driver   neo4j.Driver
	session  neo4j.Session
	result   neo4j.Result
	username string
	password string
	greeting interface{}
	uri      string
	err      error
)

func init() {
	username = "root"
	password = "system"
	uri = "bolt://localhost:7687"
}

func initDriver() error {
	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	// driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return err
	}
	return err
}

func initSession(mode neo4j.AccessMode) error {
	err = initDriver()
	if err != nil {
		return err
	}
	session, err = driver.Session(mode)
	if err != nil {
		return err
	}
	return nil
}

func ReadTran(cql string, arg map[string]interface{}) (string, error) {
	err = initSession(neo4j.AccessModeWrite)
	if err != nil {
		return "", err
	}
	defer driver.Close()
	defer session.Close()
	greeting, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run(cql, arg)
		if err != nil {
			return nil, err
		}

		node := []interface{}{}
		relation := []interface{}{}

		for {
			if result.Next() {
				fmt.Println(result.Record().Keys())
				// fmt.Println(result.Record().Keys(), result.Record().Values())
				// person := result.Record().GetByIndex(0).(neo4j.Node)
				// // fmt.Println(person.Id(), person.Labels(), person.Props())
				// group := result.Record().GetByIndex(1).(neo4j.Node)
				// // fmt.Println(group.Id(), group.Labels(), group.Props())
				// host := result.Record().GetByIndex(2).(neo4j.Node)
				// // fmt.Println(host.Id(), host.Labels(), host.Props())
				// relation := result.Record().GetByIndex(3).(neo4j.Relationship)
				// relation1 := result.Record().GetByIndex(4).(neo4j.Relationship)
				// fmt.Println(person.Props(), group.Props(), host.Props(), relation.Props(), relation.Type(), relation1.Props(), relation1.Type())

				// data, err := json.Marshal(host.Props())
				// if err != nil {
				// 	panic(err)
				// }
				// fmt.Println(string(data))

				for n, x := range result.Record().Keys() {
					tmp_rs := map[string]interface{}{}
					rs := result.Record().GetByIndex(n)
					switch v := rs.(type) {
					case neo4j.Node:
						fmt.Println("node")
						tmp_rs["name"] = x
						tmp_rs["props"] = rs.(neo4j.Node).Props()
						tmp_rs["id"] = rs.(neo4j.Node).Id()
						tmp_rs["labels"] = rs.(neo4j.Node).Labels()
						tmp_rs["type"] = "node"
						ss, _ := json.Marshal(tmp_rs)
						fmt.Println(string(ss))
						node = append(node, tmp_rs)
					case neo4j.Relationship:
						fmt.Println("relationship")
						tmp_rs["name"] = x
						tmp_rs["Id"] = rs.(neo4j.Relationship).Id()
						tmp_rs["props"] = rs.(neo4j.Relationship).Props()
						tmp_rs["type"] = rs.(neo4j.Relationship).Type()
						tmp_rs["startid"] = rs.(neo4j.Relationship).StartId()
						tmp_rs["endid"] = rs.(neo4j.Relationship).EndId()
						ss, _ := json.Marshal(tmp_rs)
						fmt.Println(string(ss))
						relation = append(relation, tmp_rs)
					case string:
						fmt.Println("string")
						tmp_rs["name"] = x
						tmp_rs["value"] = rs
						tmp_rs["type"] = "string"
						ss, _ := json.Marshal(tmp_rs)
						fmt.Println(string(ss))
						node = append(node, tmp_rs)
					default:
						fmt.Println("unknow type", v)
					}
				}
				// fmt.Println(result.Record().Keys(), result.Record().Values())
				// fmt.Println(result.Record().GetByIndex(0))
				// return result.Record().GetByIndex(0), nil
				// return "ok", nil
			} else {
				break
			}
		}

		node_string, _ := json.Marshal(node)

		fmt.Println(string(node_string))
		relation_string, _ := json.Marshal(relation)
		fmt.Println(string(relation_string))

		return "ok", nil
		// return nil, result.Err()
	})

	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}
