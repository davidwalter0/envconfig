package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	// cfg "github.com/davidwalter0/envconfig"
	cfg ".."
)

// Specification configuration structure receives envconfig
type Specification struct {
	Debug        bool          `envconfig:"Debug" usage:"enable debug mode"`
	Port         int           `short:"p" default:"8080" usage:"primary ip port"`
	User         string        `usage:"user for ..."`
	UserName     string        `envconfig:"USER_NAME"`
	Users        []string      ``
	UserArray    []string      `default:"asdfx,asdfy,asdfz,asdf0,asdf1"`
	Rate         float64       ``
	RateOfTravel float32       `default:"3.14"`
	Timeout      time.Duration `default:"720h1m3s"`
	Timeout2     time.Duration `default:"720h1m3s"`
	Int8         int8          `short:"i8" default:"127" usage:"int8 test"`
	Nint8        int8          `short:"n8" default:"-128" usage:"nint8 test"`
	Uint8        uint8         `short:"u8" default:"255" usage:"uint8 test"`
	Int16        int16         `short:"i16" default:"32767" usage:"int16 test"`
	Nint16       int16         `short:"n16" default:"-32768" usage:"nint16 test"`
	Uint16       uint16        `short:"u16" default:"65535" usage:"uint16 test"`
	Int32        int32         `short:"i32" default:"1048576" usage:"int32 test"`
	Nint32       int32         `short:"n32" default:"-1232" usage:"nint32 test"`
	Uint32       uint32        `short:"u32" default:"255" usage:"uint32 test"`
	// ColorCodes map[string]int ``
}

func main() {
	var s Specification

	err := cfg.Process("myapp", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(s)

	format := "Debug: %v\nPort: %d\nUser: %s\nUserName: %s\nRate: %f\nTimeout: %s\nTimeout2: %s\nRateOfTravel: %f\nInt8: %d\nNInt8: %d\nUInt8: %d\n"
	_, err = fmt.Printf(format, s.Debug, s.Port, s.User, s.UserName, s.Rate, s.Timeout, s.Timeout2, s.RateOfTravel, s.Int8, s.Nint8, s.Uint8)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Int16: %d\nNInt16: %d\nUInt16: %d\n", s.Int16, s.Nint16, s.Uint16)
	fmt.Printf("Int32: %d\nNInt32: %d\nUInt32: %d\n", s.Int32, s.Nint32, s.Uint32)

	fmt.Println("Users:")
	for _, u := range s.Users {
		fmt.Printf("  %s\n", u)
	}

	fmt.Println("UserArray:")
	for _, u := range s.UserArray {
		fmt.Printf("  %s\n", u)
	}

	// fmt.Println("Color codes:")
	// for k, v := range s.ColorCodes {
	// 	fmt.Printf("  %s: %d\n", k, v)
	// }

	flag.Usage()
}
