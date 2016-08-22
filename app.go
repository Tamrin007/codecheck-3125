package main

import (
  "fmt"
  "github.com/codegangsta/cli"
)

func doMain(c *cli.Context) {
  for i := 0; i<len(c.Args()); i++ {
    fmt.Println(c.Args()[i])
  }
}