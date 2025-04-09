package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"sort"
)

// Read a big-endian uint32 from a byte slice of length at least 4
func ReadBigEndianUint32(buffer []byte) uint32 {
	if len(buffer) < 4 {
		panic("buffer too short to read uint32")
	}
	return binary.BigEndian.Uint32(buffer[:])
}

// Write a big-endian uint32 to a byte slice of length at least 4
func WriteBigEndianUint32(buffer []byte, num uint32) {
	if len(buffer) < 4 {
		panic("buffer too short to write uint32")
	}
	binary.BigEndian.PutUint32(buffer, num)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	inFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v\n", err)
	}
	defer inFile.Close()

	len_buffer := make([]byte, 4)
	key_buffer := make([]byte, 10)
	var records [][]byte

	/* Store records in memory */
	for {
		// Read the length of the record
		l, err := inFile.Read(len_buffer)
		// Read until EOF
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading length: %v\n", err)
		}
		length := ReadBigEndianUint32(len_buffer[:l])
		// Read the key
		k, err := inFile.Read(key_buffer)
		if err != nil || k != 10 {
			log.Fatalf("Error reading key: %v\n", err)
		}
		// Read the value
		val_buffer := make([]byte, length-14)
		v, err := inFile.Read(val_buffer)
		if err != nil {
			log.Fatalf("Error reading value: %v\n", err)
		}
		value := val_buffer[:v]

		// Store the record
		var record []byte
		record = append(record, len_buffer[:]...)
		record = append(record, key_buffer[:]...)
		record = append(record, value...)
		records = append(records, record)
	}

	/* Sort the records in memory */
	sort.Slice(records, func(i, j int) bool {
		return bytes.Compare(records[i][4:14], records[j][4:14]) < 0
	})

	/* Write sorted records to output file */
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v\n", err)
	}
	for _, record := range records {
		n, err := outFile.Write(record)
		if err != nil || n != len(record) {
			log.Fatalf("Error writing record: %v\n", err)
		}
	}
	defer outFile.Close()
}
