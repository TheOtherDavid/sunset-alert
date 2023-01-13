package main

import (
	"fmt"
	"time"

	"github.com/TheOtherDavid/sunset-alert"
)

func main() {
	for {
		fmt.Printf("Executing Sunset Alert at %s\n", time.Now())
		alert.SunsetAlert()
	}
}
