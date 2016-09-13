package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/codegangsta/cli"
)

// IrohaCity は 路線情報を持つ.
type IrohaCity struct {
	Lines []Line `json:"lines"`
}

// Line は 路線名と駅情報を持つ.
type Line struct {
	Name     string    `json:"name"`
	Stations []Station `json:"stations"`
}

// Station は駅名、上り方面における次の駅までの所要時間、ガレージの有無、停泊する電車の数を持つ.
type Station struct {
	Name      string `json:"name"`
	Duration  int    `json:"up"`
	HasGarage bool   `json:"hasGarage"`
	Trains    int    `json:"trains"`
}

func getCityInformationFromJSON() IrohaCity {
	raw, err := ioutil.ReadFile("./iroha-city.json")
	if err != nil {
		log.Println(err)
	}
	irohaCity := IrohaCity{}
	json.Unmarshal(raw, &irohaCity)
	return irohaCity
}

func createTimeTable() map[string][]string {
	// 始発電車の時刻から終電までの時刻を求める
	// キーに時間、値は分の配列とする
	timeTable := map[string][]string{}
	var trainsByHour []string

	firstTrain, _ := time.Parse("15:04", "06:00")
	hour := firstTrain.Format("15")
	minutes := firstTrain.Format("04")
	trainsByHour = append(trainsByHour, minutes)
	timeTable[hour] = trainsByHour

	secondTrain, _ := time.Parse("15:04", "06:12")
	hour = secondTrain.Format("15")
	minutes = secondTrain.Format("04")
	trainsByHour = append(trainsByHour, minutes)
	timeTable[hour] = trainsByHour

	return timeTable
}

func printTimeTable(timeTable map[string][]string) {
	for hour, trains := range timeTable {
		fmt.Print(hour, ": ")
		for _, minutes := range trains {
			fmt.Print(minutes, " ")
		}
		fmt.Println("")
	}
}

func doMain(c *cli.Context) {
	// line := c.Args()[0]
	// station := c.Args()[1]
	// direction := c.Args()[2]
	// if len(c.Args()) == 4 {
	// 	hour, err := strconv.Atoi(c.Args()[3])
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println("hour: ", hour)
	// }
	// IrohaCity := getCityInformationFromJSON()
	timeTable := createTimeTable()
	printTimeTable(timeTable)
}
