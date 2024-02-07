package main

import "log"

func main() {
	err := setup()
	if err != nil {
		log.Fatal()
	}
}
