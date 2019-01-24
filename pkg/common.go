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
				// fmt.Println(result.Record().Keys())
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
				node_map := map[int64]map[string]interface{}{}
				for n, x := range result.Record().Keys() {
					var tmp_rs map[string]interface{}
					rs := result.Record().GetByIndex(n)
					switch v := rs.(type) {
					case neo4j.Node:
						// fmt.Println("node")
						if _, ok := node_map[rs.(neo4j.Node).Id()]; !ok {
							tmp_rs = map[string]interface{}{}
							tmp_rs["group"] = x
							tmp_rs["props"] = rs.(neo4j.Node).Props()
							tmp_rs["id"] = rs.(neo4j.Node).Id()
							tmp_rs["labels"] = rs.(neo4j.Node).Labels()
							tmp_rs["type"] = "node"
							// ss, _ := json.Marshal(tmp_rs)
							// fmt.Println(string(ss))
							// node = append(node, tmp_rs)
							node_map[rs.(neo4j.Node).Id()] = tmp_rs
						}

					case neo4j.Relationship:
						tmp_rs = map[string]interface{}{}
						// fmt.Println("relationship")
						tmp_rs["name"] = x
						// tmp_rs["id"] = rs.(neo4j.Relationship).Id()
						tmp_rs["props"] = rs.(neo4j.Relationship).Props()
						tmp_rs["relation"] = rs.(neo4j.Relationship).Type()
						tmp_rs["source"] = rs.(neo4j.Relationship).StartId()
						tmp_rs["target"] = rs.(neo4j.Relationship).EndId()
						tmp_rs["value"] = 1
						// ss, _ := json.Marshal(tmp_rs)
						// fmt.Println(string(ss))
						relation = append(relation, tmp_rs)
					case string:
						tmp_rs = map[string]interface{}{}
						// fmt.Println("string")
						tmp_rs["group"] = x
						tmp_rs["value"] = rs
						tmp_rs["type"] = "string"
						// ss, _ := json.Marshal(tmp_rs)
						// fmt.Println(string(ss))
						node = append(node, tmp_rs)
					// case int64:
					// 	tmp_rs["length"] = rs
					default:
						fmt.Println("unknow type", v)
					}
				}
				// fmt.Println(result.Record().Keys(), result.Record().Values())
				// fmt.Println(result.Record().GetByIndex(0))
				// return result.Record().GetByIndex(0), nil
				// return "ok", nil

				for _, v := range node_map {
					node = append(node, v)
				}
			} else {
				break
			}
		}

		countNode := map[int64]int{}
		for n1, x := range node {
			countNode[x.(map[string]interface{})["id"].(int64)] = n1
		}
		node_string, _ := json.Marshal(node)

		fmt.Println(string(node_string))
		fmt.Println("##################################################")
		for n2, y := range relation {
			y.(map[string]interface{})["source"] = countNode[y.(map[string]interface{})["source"].(int64)]
			y.(map[string]interface{})["target"] = countNode[y.(map[string]interface{})["target"].(int64)]
			relation[n2] = y
		}
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
