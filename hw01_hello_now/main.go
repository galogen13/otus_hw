package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal(err)
	}
	locTime := t.Local()
	fmt.Println("current time:", locTime.Round(time.Minute))
	fmt.Println("exact time:", locTime.Round(time.Nanosecond))
}
