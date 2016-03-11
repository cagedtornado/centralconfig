package main

import (
	"fmt"

	"github.com/danesparza/centralconfig/datastores"
)

func main() {

	fmt.Println("Initializing BoltDB database...")

	db := datastores.BoltDB{
		Database: "testing.db"}

	//	Initialize the store
	db.InitStore(true)

	//	Try storing some config items:
	ct1 := datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem1",
		Value:       "Value1"}
	db.Set(ct1)

	ct2 := datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem2",
		Value:       "Value2"}
	db.Set(ct2)

	//	Get a config item:
	query := datastores.ConfigItem{
		Application: "Formbuilder",
		Name:        "TestItem2"}
	ct3, _ := db.Get(query)
	fmt.Printf("Found config item: %s value: %s \n", ct3.Name, ct3.Value)

	//	Get all config items:
	cis, _ := db.GetAll("Formbuilder")
	for _, configItem := range cis {
		fmt.Printf("%s - %s:%s\n", configItem.Application, configItem.Name, configItem.Value)
	}
}
