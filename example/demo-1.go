package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	cfg "github.com/davidwalter0/envconfig"
)

// Specification configuration structure receives envconfig
type Specification struct {
	Debug        bool   `envconfig:"Debug" usage:"enable debug mode"`
	Port         int    `short:"p" default:"8080" usage:"primary ip port"`
	User         string `usage:"user for ..."`
	UserName     string `envconfig:"USER_NAME"`
	Users        []string
	Rate         float64
	RateOfTravel float64
	Timeout      time.Duration
	// ColorCodes map[string]int
}

func main() {
	var s Specification

	err := cfg.Process("myapp", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(s)

	format := "Debug: %v\nPort: %d\nUser: %s\nUserName: %s\nRate: %f\nTimeout: %s\n"
	_, err = fmt.Printf(format, s.Debug, s.Port, s.User, s.UserName, s.Rate, s.Timeout)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Users:")
	for _, u := range s.Users {
		fmt.Printf("  %s\n", u)
	}

	// fmt.Println("Color codes:")
	// for k, v := range s.ColorCodes {
	// 	fmt.Printf("  %s: %d\n", k, v)
	// }

	flag.Usage()
}
