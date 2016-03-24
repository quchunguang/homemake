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
	STORECOUNT = 500
	CONNRETRY  = 3
	DATACOUNT  = 9
	UPLOADURL  = "https://api.myjson.com/bins/3wczx"
	DATAFILE   = "lightrelay.tsv"
	LOGFILE    = "lightrelay.log"
)

var Headers = []string{"temp", "light"}

func getData(reader *bufio.Reader, records *[]record) error {
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
			return err
		}

		var record = record{}
		record["time"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		s := record["time"]
		line := strings.TrimSpace(string(reply))
		if strings.HasPrefix(line, "[") {
			fmt.Print("i")
			s = line + " " + s + "\n"
			logfile.Write([]byte(s))
			continue
		}
		fmt.Print(".")
		for j, v := range strings.Split(line, " ") {
			record[Headers[j]] = v
			s += "\t" + v
		}
		s += "\n"
		tsvfile.Write([]byte(s))
		*records = append(*records, record)
	}
	return nil
}

func upload(jsonstr []byte) error {
	req, err := http.NewRequest("PUT", UPLOADURL, bytes.NewBuffer(jsonstr))
	if err != nil {
		return err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	resp.Body.Close()

	return err
}

func main() {
	var records []record

	c := &serial.Config{Name: DEVICE, Baud: BAUD}
	s, err := serial.OpenPort(c)
	defer s.Close()
	if err != nil {
		fmt.Println("[ERROR] Can not connect serial port", DEVICE)
		os.Exit(1)
	}

	reader := bufio.NewReader(s)
	for {
		fmt.Print("Reading")
		err = getData(reader, &records)
		if err != nil {
			fmt.Println("Fail")
			continue
		}
		if len(records) > STORECOUNT {
			records = records[len(records)-STORECOUNT:]
		}
		rets, _ := json.Marshal(records)

		now := time.Now().Format("2006-01-02T15:04:05Z")
		fmt.Print(" ", now, " Uploading")
		for i := 0; i < CONNRETRY; i++ {
			fmt.Print(".")
			err := upload(rets)
			if err == nil {
				fmt.Println("Done")
				break
			} else if i == CONNRETRY-1 {
				fmt.Println("Fail")
			}
		}
	}
}
