package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
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
	// 始発電車の時刻から終電までの時刻を求める
	// キーに時間、値は分の配列とする
	IrohaCity := getCityInformationFromJSON()
	timeTable := map[string][]string{}
	var (
		firstTrain     time.Time
		firstTrainAtA7 time.Time
		interval       int
	)

	if line == "A" {
		firstTrain, _ = time.Parse("15:04", "05:55")
		firstTrainAtA7, _ = time.Parse("15:04", "06:10")
		interval = 5
	}
	if line == "B" {
		firstTrain, _ = time.Parse("15:04", "06:00")
		interval = 6
	}

	limit, _ := time.Parse("15:04", "23:00")
	delay := 0
	delayAtA7 := 0
	for _, Line := range IrohaCity.Lines {
		if Line.Name == line {
			for _, Station := range Line.Stations {
				if Station.Name == "A7" {
					delayAtA7 = delay
				}
				if Station.Name == station {
					break
				}
				delay += Station.Duration
			}
		}
	}

	firstTrain = firstTrain.Add(time.Duration(delay) * time.Minute)
	firstTrainAtA7 = firstTrainAtA7.Add(time.Duration(delay-delayAtA7) * time.Minute)
	limit = limit.Add(time.Duration(delay) * time.Minute)
	for i, train := 0, firstTrain; train.Before(limit); train = train.Add(time.Duration(interval) * time.Minute) {
		// 1 本おきに A7 行と A13 行 かつ 始発は A7 行
		re := regexp.MustCompile("A")
		stationNum, _ := strconv.Atoi(re.ReplaceAllString(station, ""))
		if line == "A" && stationNum >= 7 && i == 0 {
			timeTable["06"] = append(timeTable[firstTrainAtA7.Format("15")], firstTrainAtA7.Format("04"))
		}
		if line == "A" && stationNum >= 7 && i%2 == 0 {
			i++
			continue
		}
		hour := train.Format("15")
		minutes := train.Format("04")
		timeTable[hour] = append(timeTable[hour], minutes)
		i++
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
