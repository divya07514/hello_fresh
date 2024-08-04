package main

import (
	"divya07514-hftest-golang/db"
	"divya07514-hftest-golang/file_reader"
	"divya07514-hftest-golang/logger"
	"divya07514-hftest-golang/service"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"
)

func main() {
	// Define command-line flags
	var inputFile string
	var postCode string
	var fromTime string
	var toTime string
	var recipeList string
	var files = make(map[string]bool)
	flag.StringVar(&inputFile, "input", "", "Path to the JSON file")
	flag.StringVar(&postCode, "postcode", "", "Custom Postcode")
	flag.StringVar(&fromTime, "fromTime", "", "Custom From Time")
	flag.StringVar(&toTime, "toTime", "", "Custom To Time")
	flag.StringVar(&recipeList, "recipe", "", "Recipe Names to search for. Provide this as comma separate with no spaces in between. e.g Potato,Steak,Tomato")
	flag.Parse()
	files[""] = true
	fmt.Println(inputFile)
	dbConn, _ := db.GetMySqlDbConnection("root:root@tcp(192.168.5.2:3306)/hello_fresh")
	database := db.MySqlDB{Connection: dbConn}
	_ = database.CreateTable()
	checkInputFile(inputFile, files, &database)
	postCode, fromTime, toTime, list := setupDefaults(postCode, fromTime, toTime, recipeList)
	generateReport(postCode, fromTime, toTime, list, &database)
}

func generateReport(postCode string, fromTime string, toTime string, list []string, db db.DataSource) {
	stats := &service.GlobalStats{Db: db}
	report, _ := stats.Report(postCode, fromTime, toTime, list)
	jsonData, _ := json.MarshalIndent(report, "", "    ")
	if jsonData != nil {
		fmt.Println(string(jsonData))
	}
}

func checkInputFile(inputFile string, files map[string]bool, db db.DataSource) {
	if _, exists := files[inputFile]; !exists {
		reader := &file_reader.JsonFileReader{
			Wg: sync.WaitGroup{},
			Db: db,
		}
		err := reader.ProcessFile(inputFile)
		if err != nil {
			logger.ErrorLogger.Error().Msg("Error reading file")
		}
		files[inputFile] = true
	}
}

func setupDefaults(postCode string, fromTime string, toTime string, recipeList string) (string, string, string, []string) {
	if postCode == "" {
		postCode = "10120"
	}
	if fromTime == "" {
		fromTime = "10AM"
	}
	if toTime == "" {
		toTime = "3PM"
	}
	var list []string
	if recipeList == "" {
		list = append(list, "Potato")
		list = append(list, "Veggie")
		list = append(list, "Mushroom")
	} else {
		list = strings.Split(recipeList, ",")
	}
	return postCode, fromTime, toTime, list
}
