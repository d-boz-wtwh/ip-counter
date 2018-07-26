package main

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"regexp"
	"strconv"
	"io/ioutil"
	"fmt"
	"strings"
)

func main() {
	fileInfo, err := ioutil.ReadDir(".")
	if err != nil {
		log.Printf("error reading file %s", err)
		return
	}
	uniqueIPs := map[string]int{}
	for i, file := range fileInfo {
		if file.Name() == "count.txt"  || strings.Contains(file.Name(), ".go") || strings.Contains(file.Name(),
			"main"){
			continue
		}
		fmt.Println("Files Remaining: ", len(fileInfo)-i)
		if file.IsDir() {
			fmt.Println("File Is A Directory: ", file.Name())
			continue
		}

		if strings.Contains(file.Name(), ".gz") {
			fmt.Println("GZ File: ", file.Name())
			ReadGZipFile(file.Name(), &uniqueIPs)
		} else {
			fmt.Println("File: ", file.Name())
			ReadFile(file.Name(), &uniqueIPs)
		}

		fmt.Println("TOTAL UNIQUE IPS", len(uniqueIPs))
		WriteToFile(uniqueIPs)
		uniqueIPs = map[string]int{}
	}
}

// TODO: Write Regular Expression To Break Into IP Reads / Line
func ScanIPs(data []byte, atEOF bool) (advance int, token []byte, err error) {

	return 0, nil, nil
}

func ReadFile(fileName string, ips *map[string]int) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open: %s", fileName)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scan(scanner, ips)
}

func ReadGZipFile(fileName string, ips *map[string]int) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open %s", fileName)
		return
	}
	defer f.Close()

	reader, err := gzip.NewReader(f)
	if err != nil {
		log.Printf("Error creating new gzip file reader: %s", err)
		return
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	scan(scanner, ips)
}

func WriteToFile(m map[string]int) {
	f, err := os.OpenFile("./count.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Cant open count file: ", err)
		return
	}

	fMap := map[string]int{}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var validIP = regexp.MustCompile(` `)
	for scanner.Scan() {
		stringArr := validIP.Split(scanner.Text(), 2)
		i, err := strconv.Atoi(stringArr[1])
		if err != nil {
			log.Println("Can't convert...")
			continue
		}
		fMap[stringArr[0]] = i
	}

	f.Close()
	for k := range m {
		fMap[k]++
	}

	err = ioutil.WriteFile("./count.txt", []byte(""), 0666)
	if err != nil {
		log.Println("error removing all data from count.txt ", err)
	}

	newFile, err := os.OpenFile("./count.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Cant open count file: ", err)
		return
	}
	for k,v := range fMap {
		count := strconv.Itoa(v)
		newFile.WriteString(k)
		newFile.WriteString(" ")
		newFile.WriteString(count)
		newFile.WriteString("\n")
	}
	newFile.Close()
}

func scan(s *bufio.Scanner, m *map[string]int) {
	var validIP = regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3}`)
	for s.Scan() {
		str := validIP.Find([]byte(s.Text()))
		(*m)[string(str)]++
	}
}
