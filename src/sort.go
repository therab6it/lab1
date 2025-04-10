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

type Record struct {
	length []byte
	key    []byte
	value  []byte
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	var err error

	inFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	len_buffer := make([]byte, 4)
	key_buffer := make([]byte, 10)
	var records []Record

	/* Store records in memory */
	for {
		// Read the length of the record
		_, err := io.ReadFull(inFile, len_buffer)
		// Read until EOF
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading input file: %v\n", err)
		}
		length := ReadBigEndianUint32(len_buffer[:])
		// Read the key
		_, err = io.ReadFull(inFile, key_buffer)
		if err != nil {
			log.Fatalf("Error reading key: %v\n", err)
		}
		// Read the value
		var val_buffer []byte
		if length != 10 {
			val_buffer = make([]byte, length-10)
			_, err = io.ReadFull(inFile, val_buffer)
			if err != nil {
				log.Fatalf("Error reading value: %v\n", err)
			}
		} else {
			val_buffer = make([]byte, 0)
		}

		// Store the record
		record := Record{
			length: append([]byte{}, len_buffer...),
			key:    append([]byte{}, key_buffer...),
			value:  append([]byte{}, val_buffer...),
		}
		records = append(records, record)
	}

	/* Sort the records */
	sort.Slice(records, func(i, j int) bool {
		return bytes.Compare(records[i].key, records[j].key) < 0
	})

	/* Write sorted records to output file */
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	for _, record := range records {
		_, err = outFile.Write(record.length)
		if err != nil {
			log.Fatalf("Error writing length: %v\n", err)
		}
		_, err = outFile.Write(record.key)
		if err != nil {
			log.Fatalf("Error writing key: %v\n", err)
		}
		_, err = outFile.Write(record.value)
		if err != nil {
			log.Fatalf("Error writing record: %v\n", err)
		}
	}
}
