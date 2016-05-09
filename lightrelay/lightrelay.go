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
	CONNRETRY  = 5
	DATACOUNT  = 10
	URL_SECOND = "https://api.myjson.com/bins/3wczx"
	URL_MINUTE = "https://api.myjson.com/bins/11yle"
	URL_HOUR   = "https://api.myjson.com/bins/4xhqq"
	DATAFILE   = "lightrelay.tsv"
)

// Headers of data fields
var Headers = []string{"temp", "light"}

func getData(reader *bufio.Reader, records *[]record) {
	fmt.Print(time.Now().Format("2006-01-02T15:04:05 "))

	// Open tsv file for append data
	tsvfile, err := os.OpenFile(DATAFILE,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer tsvfile.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not open data file for append", DATAFILE)
		os.Exit(2)
	}

	// Get DATACOUNT points of data
	for i := 0; i < DATACOUNT; i++ {
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			fmt.Print("D")
			i--
			continue
		}

		// Add `time` column, with UTC standard time format
		var record = record{}
		record["time"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		s := record["time"]
		line := strings.TrimSpace(string(reply))

		// Process data
		for j, v := range strings.Split(line, "\t") {
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
	// Reference at http://myjson.com/api
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
		// PUT method to update data entirely
		// Reference at http://myjson.com/api
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
	return err // Retry times out, give up
}

func uploadList(records []record, url string, maxlen int) {
	// Keep length of list less than maxlen
	if len(records) > maxlen {
		records = records[len(records)-maxlen:]
	}

	// Convert list to JSON string and upload
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
			// Upload second-pt data to web every 10 seconds
			uploadList(seconds, URL_SECOND, 500)

			// Insert a data pt every minute, 6 loops
			if i%6 == 0 {
				minutes = append(minutes, seconds[len(seconds)-1])
			}

			// Upload minute-pt data to web every 10 minutes
			if i%60 == 0 {
				uploadList(minutes, URL_MINUTE, 2000)
			}

			// Insert a data pt every hour, 6*60 loops
			if i%360 == 0 {
				hours = append(hours, seconds[len(seconds)-1])
			}

			// Upload hour-pt data to web every 10 hours
			if i%3600 == 0 {
				uploadList(hours, URL_HOUR, 10000)
			}
		}()
	}
}
