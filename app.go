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

func calcDelay(IrohaCity IrohaCity, line, direction, station string) (delay, delayFromA7 int) {
	for _, Line := range IrohaCity.Lines {
		if Line.Name == line && direction == "U" {
			for _, Station := range Line.Stations {
				if Station.Name == "A7" {
					delayFromA7 = delay
				}
				if Station.Name == station {
					break
				}
				delay += Station.Duration
			}
		}
		if Line.Name == line && direction == "D" {
			for i := len(Line.Stations) - 1; i >= 0; i-- {
				delay += Line.Stations[i].Duration
				if Line.Stations[i].Name == "A7" {
					delayFromA7 = delay
				}
				if Line.Stations[i].Name == station {
					break
				}
			}
		}
	}
	return delay, delayFromA7
}

func setInterval(line, direction string) (interval int) {
	if line == "A" && direction == "U" {
		interval = 5
	}
	if line == "A" && direction == "D" {
		interval = 5
	}
	if line == "B" && direction == "U" {
		interval = 6
	}
	if line == "B" && direction == "D" {
		interval = 6
	}
	return interval
}

func setLimit(line, direction string, delay int) (limit time.Time) {
	if line == "A" && direction == "U" {
		limit, _ = time.Parse("15:04", "23:00")
	}
	if line == "A" && direction == "D" {
		limit, _ = time.Parse("15:04", "23:00")
	}
	if line == "B" && direction == "U" {
		limit, _ = time.Parse("15:04", "23:00")
	}
	if line == "B" && direction == "D" {
		// TODO: B5 U 最終電車の発車時刻にしたい
		// 上りテーブルも下りテーブルも作れば参照できる？
		limit, _ = time.Parse("15:04", "23:06")
	}
	limit = limit.Add(time.Duration(delay) * time.Minute)
	return limit
}

func setFirstTrain(line, direction string, delay int) (firstTrain time.Time) {
	if line == "A" && direction == "U" {
		firstTrain, _ = time.Parse("15:04", "05:55")
	}
	if line == "A" && direction == "D" {
		firstTrain, _ = time.Parse("15:04", "05:52")
	}
	if line == "B" && direction == "U" {
		firstTrain, _ = time.Parse("15:04", "06:00")
	}
	if line == "B" && direction == "D" {
		firstTrain, _ = time.Parse("15:04", "06:11")
	}
	firstTrain = firstTrain.Add(time.Duration(delay) * time.Minute)
	return firstTrain
}

func createTimeTable(line string, station string, direction string) map[string][]string {
	// 始発電車の時刻から終電までの時刻を求める
	// キーに時間、値は分の配列とする
	IrohaCity := getCityInformationFromJSON()
	timeTable := map[string][]string{}
	var (
		firstTrainFromA7 time.Time
		limitAtA7        time.Time
	)

	delay, delayFromA7 := calcDelay(IrohaCity, line, direction, station)
	firstTrain := setFirstTrain(line, direction, delay)
	interval := setInterval(line, direction)
	limit := setLimit(line, direction, delay)

	if line == "A" && direction == "U" {
		firstTrainFromA7, _ = time.Parse("15:04", "06:10")
	}
	if line == "A" && direction == "D" {
		firstTrainFromA7, _ = time.Parse("15:04", "06:06")
		// TODO: A8 D 最終電車の発車時刻にしたい
		limitAtA7, _ = time.Parse("15:04", "23:07")
	}
	firstTrainFromA7 = firstTrainFromA7.Add(time.Duration(delay-delayFromA7) * time.Minute)

	re := regexp.MustCompile("A")
	stationNum, _ := strconv.Atoi(re.ReplaceAllString(station, ""))
	for i, train := 0, firstTrain; train.Before(limit); train = train.Add(time.Duration(interval) * time.Minute) {
		// A7 発 A13 行きが 06:10 に 1 本
		if line == "A" && direction == "U" && stationNum >= 7 && i == 0 {
			timeTable["06"] = append(timeTable[firstTrainFromA7.Format("15")], firstTrainFromA7.Format("04"))
		}
		// 1 本おきに A7 行と A13 行 かつ 始発は A7 行
		if line == "A" && direction == "U" && stationNum >= 7 && i%2 == 0 {
			i++
			continue
		}
		// A7 発 A1 行きが 06:06 に 1 本（これより後は A13 からの 5 分毎に発車しているとみなせる）
		// TODO: A13 発と A7 発は分けて考え、 A7 以下は結合させる
		// ↑にすれば最終発車時刻の決め打ちをやめることができそう
		if line == "A" && direction == "D" && stationNum <= 7 && i == 0 {
			timeTable["06"] = append(timeTable[firstTrainFromA7.Format("15")], firstTrainFromA7.Format("04"))
		}
		// A13 - A8 は 10 分に 1 本なので間引く
		if line == "A" && direction == "D" && stationNum > 7 && i%2 == 1 {
			i++
			continue
		}
		hour := train.Format("15")
		minutes := train.Format("04")
		// A7 の最終より後の電車は無視
		if line == "A" && direction == "D" && stationNum <= 7 {
			if train.After(limitAtA7.Add(time.Duration(delay-delayFromA7) * time.Minute)) {
				break
			}
		}
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
