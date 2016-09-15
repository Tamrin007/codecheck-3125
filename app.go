package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
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

func createTimeTable(line string, station string, direction string) map[string][]string {
	IrohaCity := getCityInformationFromJSON()
	delay := 0
	for _, Line := range IrohaCity.Lines {
		if Line.Name == line {
			for _, Station := range Line.Stations {
				if Station.Name == station {
					break
				}
				delay += Station.Duration
			}
		}
	}
	// 始発電車の時刻から終電までの時刻を求める
	// キーに時間、値は分の配列とする
	timeTable := map[string][]string{}

	firstTrain, _ := time.Parse("15:04", "06:00")
	firstTrain = firstTrain.Add(time.Duration(delay) * time.Minute)
	limit, _ := time.Parse("15:04", "23:00")
	limit = limit.Add(time.Duration(delay) * time.Minute)
	for train := firstTrain; train.Before(limit); train = train.Add(6 * time.Minute) {
		hour := train.Format("15")
		minutes := train.Format("04")
		timeTable[hour] = append(timeTable[hour], minutes)
	}

	return timeTable
}

func printAllTimeTable(timeTable map[string][]string) {
	var keys []string
	for hour := range timeTable {
		keys = append(keys, hour)
	}
	sort.Strings(keys)
	for _, hour := range keys {
		output := fmt.Sprint(hour, ":")
		for _, minutes := range timeTable[hour] {
			output += fmt.Sprint(" ", minutes)
		}
		fmt.Println(output)
	}
}

func printHourlyTimeTable(timeTable map[string][]string, hour string) {
	var trains string
	for _, minutes := range timeTable[hour] {
		trains += fmt.Sprint(" ", minutes)
	}
	if trains == "" {
		fmt.Println("No train")
		return
	}
	fmt.Print(hour, ":")
	fmt.Println(trains)
}

func doMain(c *cli.Context) {
	line := c.Args()[0]
	station := c.Args()[1]
	direction := c.Args()[2]
	timeTable := createTimeTable(line, station, direction)
	if len(c.Args()) == 4 {
		hour, err := strconv.Atoi(c.Args()[3])
		if err != nil {
			fmt.Println(err)
		}
		zeroPaddedHour := fmt.Sprintf("%02d", hour)
		printHourlyTimeTable(timeTable, zeroPaddedHour)
	} else {
		printAllTimeTable(timeTable)
	}
}
