package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"net/http"
	"os"
	"strings"
	"time"
)

type record map[string]string

const (
	BAUD       = 9600
	STORECOUNT = 1000
	CONNRETRY  = 5
	DATACOUNT  = 9
	URL_SECOND = "https://api.myjson.com/bins/3wczx"
	URL_MINUTE = "https://api.myjson.com/bins/11yle"
	URL_HOUR   = "https://api.myjson.com/bins/4xhqq"
	DATAFILE   = "lightrelay.tsv"
	LOGFILE    = "lightrelay.log"
)

var Headers = []string{"temp", "light"}

func getData(reader *bufio.Reader, records *[]record) error {
	fmt.Print("Reading")

	tsvfile, err := os.OpenFile(DATAFILE,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer tsvfile.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not open data file for append", DATAFILE)
		os.Exit(2)
	}

	logfile, err := os.OpenFile(LOGFILE,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer logfile.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not open log file for append", LOGFILE)
		os.Exit(3)
	}

	for i := 0; i < DATACOUNT; i++ {
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			fmt.Print("x")
			continue
		}

		var record = record{}
		record["time"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		s := record["time"]
		line := strings.TrimSpace(string(reply))
		if strings.HasPrefix(line, "[") {
			s = line + " " + s + "\n"
			logfile.Write([]byte(s))
			fmt.Print("i")
			continue
		}

		for j, v := range strings.Split(line, " ") {
			record[Headers[j]] = v
			s += "\t" + v
		}
		s += "\n"
		tsvfile.Write([]byte(s))
		*records = append(*records, record)
		fmt.Print(".")
	}
	return err
}

func upload(jsonstr []byte, url string) error {
	var err error
	now := time.Now().Format("2006-01-02T15:04:05.000Z")
	fmt.Print(" ", now, " Uploading")

	for i := 0; i < CONNRETRY; i++ {
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonstr))
		if err != nil {
			fmt.Print("x")
			continue
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("x")
			continue
		}

		if resp.StatusCode != 200 {
			err = errors.New(resp.Status)
			fmt.Print("x")
			continue
		}
		resp.Body.Close()
		fmt.Print(".")
		return nil // OK
	}
	return err
}

func uploadlist(records []record, url string) {
	if len(records) > STORECOUNT {
		records = records[len(records)-STORECOUNT:]
	}
	jsonstr, _ := json.Marshal(records)
	err := upload(jsonstr, url)
	if err != nil {
		fmt.Println("Fail")
	} else {
		fmt.Println("Done")
	}
}

func main() {
	var seconds, minutes, hours []record

	c := &serial.Config{Name: DEVICE, Baud: BAUD}
	s, err := serial.OpenPort(c)
	defer s.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not connect serial port", DEVICE)
		os.Exit(1)
	}

	reader := bufio.NewReader(s)
	for i := 0; ; i++ {
		err = getData(reader, &seconds)
		if err != nil {
			fmt.Println("Fail")
			continue
		}

		uploadlist(seconds, URL_SECOND)

		if i%60 == 0 {
			minutes = append(minutes, seconds[len(seconds)-1])
			uploadlist(minutes, URL_MINUTE)
		}

		if i%3600 == 0 {
			hours = append(hours, seconds[len(seconds)-1])
			uploadlist(hours, URL_HOUR)
		}
	}
}
