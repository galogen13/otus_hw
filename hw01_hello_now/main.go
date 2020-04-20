package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("current time: %s\n", t.Round(time.Second))
	fmt.Printf("exact time: %s\n", t.Local())
}
