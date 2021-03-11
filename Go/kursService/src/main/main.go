package main

import (
	"core"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: kursService portNumber compressionType kursFolder")
		os.Exit(1)
	}

	portNumber, err := strconv.Atoi(os.Args[1])
	checkErr(err)

	var compressionType int
	compressionType, err = core.GetCompressionType(os.Args[2])
	checkErr(err)

	err = core.ServerStart(portNumber, compressionType, os.Args[3])
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
