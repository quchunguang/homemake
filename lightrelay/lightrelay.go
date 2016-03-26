package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io/ioutil"
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

func getData(reader *bufio.Reader, records *[]record) {
	fmt.Print(time.Now().Format("2006-01-02T15:04:05 "))

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
			fmt.Print("D")
			continue
		}

		var record = record{}
		record["time"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		s := record["time"]
		line := strings.TrimSpace(string(reply))
		if strings.HasPrefix(line, "[") {
			s = s + " " + line + "\n"
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
	fmt.Println("")
}

func downloadList(url string, records *[]record) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, records)
	return err
}

func upload(jsonstr []byte, url string) error {
	var err error

	for i := 0; i < CONNRETRY; i++ {
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonstr))
		if err != nil {
			continue
		}
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			err = errors.New(resp.Status)
			continue
		}
		resp.Body.Close()
		return nil // OK
	}
	return err
}

func uploadList(records []record, url string) {
	if len(records) > STORECOUNT {
		records = records[len(records)-STORECOUNT:]
	}
	jsonstr, _ := json.Marshal(records)
	err := upload(jsonstr, url)
	if err != nil {
		fmt.Println("U")
	} else {
		fmt.Print("u")
	}
}

func main() {
	var seconds, minutes, hours []record

	// Initialize reading data from web
	downloadList(URL_SECOND, &seconds)
	downloadList(URL_MINUTE, &minutes)
	downloadList(URL_HOUR, &hours)

	// Initialize serial port for reading
	c := &serial.Config{Name: DEVICE, Baud: BAUD}
	s, err := serial.OpenPort(c)
	defer s.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not connect serial port", DEVICE)
		os.Exit(1)
	}

	reader := bufio.NewReader(s)
	for i := 0; ; i++ {
		// Get a batch of data from serial port and add to `second` list
		getData(reader, &seconds)

		go func() {
			// Upload every-second-pt data to web
			// Test: About 16s every for-loop. About 1.7s every data point
			uploadList(seconds, URL_SECOND)

			// Upload every-minute-pt data to web
			// Test: 4 loops makes a minute
			if i%4 == 0 {
				minutes = append(minutes, seconds[len(seconds)-1])
				uploadList(minutes, URL_MINUTE)
			}

			// Upload every-hour-pt data to web
			// Test: 4*60 loops makes a hour
			if i%240 == 0 {
				hours = append(hours, seconds[len(seconds)-1])
				uploadList(hours, URL_HOUR)
			}
		}()
	}
}
