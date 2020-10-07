package main

import (
	c "github.com/Ichliebedich0505/T_traceroute/trace_cli"
	"log"
	"os"
)

func main() {
	if err := c.Run(os.Args); err != nil {
		log.Printf("%v", err)
	}
}
