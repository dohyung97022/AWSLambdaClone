package main

import (
	"fmt"
	"src/lambda"
	"src/mongodb"
	"src/server"
)

func main() {
	if err := mongodb.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := mongodb.Disconnect(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	lambda.SetController()
	if err := server.Listen(); err != nil {
		fmt.Println(err)
		return
	}
}
