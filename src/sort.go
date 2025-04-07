package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	test()

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
}

func test() {

	// This is a test function to ensure the package compiles.
	// You can add test cases here if needed.
	log.Println("This is a test function in the sort package.")
}
