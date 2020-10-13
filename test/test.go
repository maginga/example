package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fileName := "/Users/hansonjang/Documents/Grandview/APTC/Log-500items.dat"
	file, err := os.Open(fileName)
	if err != nil {

	}
	defer file.Close()

	for {

		var tot uint32
		totBytes, e := readNextBytes(file, 4)
		if e != nil && e == io.EOF {
			fmt.Print(e)
			break
		}

		b1 := bytes.NewBuffer(totBytes)
		err = binary.Read(b1, binary.LittleEndian, &tot)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Println(tot)

		var datacount uint32
		dcBytes, _ := readNextBytes(file, 4)

		b2 := bytes.NewBuffer(dcBytes)
		err = binary.Read(b2, binary.LittleEndian, &datacount)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Println(datacount)

		var dtv [23]byte
		dtBytes, _ := readNextBytes(file, 23)

		b3 := bytes.NewBuffer(dtBytes)
		err = binary.Read(b3, binary.LittleEndian, &dtv)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		strByte := string(dtv[:23]) //BytesToString(datetime)
		fmt.Println(strByte)

		var null byte
		nullBytes, _ := readNextBytes(file, 1)
		nullBuffer := bytes.NewBuffer(nullBytes)
		err = binary.Read(nullBuffer, binary.LittleEndian, &null)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		itemCount := int(datacount)
		for i := 0; i < itemCount; i++ {
			var dataPoint float32
			dpBytes, _ := readNextBytes(file, 4)
			dpBuffer := bytes.NewBuffer(dpBytes)
			err = binary.Read(dpBuffer, binary.LittleEndian, &dataPoint)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}
			fmt.Println(dataPoint)
		}

		var crlf []byte
		crlfBytes, _ := readNextBytes(file, 2)
		crlfBuffer := bytes.NewBuffer(crlfBytes)
		err = binary.Read(crlfBuffer, binary.LittleEndian, &crlf)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		fmt.Println(crlf)
	}
}

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	return bytes, err
}

func BytesToString(data []byte) string {
	return string(data[:])
}
