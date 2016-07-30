package main

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"time"

	"github.com/xiam/exif"
)

var dateTags = [2]string{
	"Date and Time (Original)",
	"Date and Time"}

func main() {
	data, err := exif.Read("example1.jpg")

	if err != nil {
		panic(err)
	}

	for key, val := range data.Tags {
		fmt.Printf("%s = %s\n", key, val)
	}

	fmt.Printf("GetTAG: %s\n", findDate(data.Tags))
	hash, _ := hashFileCrc32("example1.jpg", 0xedb88320)
	fmt.Printf("HASH: %s\n", hash)
}

func findDate(tags map[string]string) time.Time {
	locSP, _ := time.LoadLocation("America/Sao_Paulo")
	var dateStr string
	for _, val := range dateTags {
		dateStr = tags[val]
		if dateStr != "" {
			date, _ := time.ParseInLocation("2006:01:02 15:04:05", dateStr, locSP)
			return date
		}
	}
	return time.Now()
}

func hashFileCrc32(filePath string, polynomial uint32) (string, error) {
	var returnCRC32String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnCRC32String, err
	}
	defer file.Close()
	tablePolynomial := crc32.MakeTable(polynomial)
	hash := crc32.New(tablePolynomial)
	if _, err := io.Copy(hash, file); err != nil {
		return returnCRC32String, err
	}
	hashInBytes := hash.Sum(nil)[:]
	returnCRC32String = hex.EncodeToString(hashInBytes)
	return returnCRC32String, nil

}
