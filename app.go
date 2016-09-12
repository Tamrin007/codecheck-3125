package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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

func doMain(c *cli.Context) {
	for i := 0; i < len(c.Args()); i++ {
		fmt.Println(c.Args()[i])
	}
	IrohaCity := getCityInformationFromJSON()
	for _, line := range IrohaCity.Lines {
		fmt.Println(line.Name)
		for _, station := range line.Stations {
			fmt.Print(station.Name + " ")
		}
		fmt.Println("")
	}
}
