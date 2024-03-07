package main

import (
	"fmt"
	"src/lambdaClone"
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
	if err := lambdaClone.SetS3CClient(); err != nil {
		fmt.Println(err)
		return
	}
	lambdaClone.SetController()
	if err := server.Listen(); err != nil {
		fmt.Println(err)
		return
	}
}
