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
	STORECOUNT = 500
	CONNRETRY  = 3
	DATACOUNT  = 9
)

var Headers = []string{"temp", "light"}

func getData(reader *bufio.Reader, records *[]record) error {
	for i := 0; i < DATACOUNT; i++ {
		reply, err := reader.ReadBytes('\x0a')
		fmt.Print(".")
		if err != nil {
			return err
		}

		var record = record{}
		record["time"] = time.Now().Format("2006-01-02T15:04:05Z")

		line := strings.TrimSpace(string(reply))
		for j, v := range strings.Split(line, " ") {
			record[Headers[j]] = v
		}
		*records = append(*records, record)
	}
	return nil
}

func upload(jsonstr []byte) error {
	url := "https://api.myjson.com/bins/3wczx"

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}

	return err
}

func main() {
	var records []record

	c := &serial.Config{Name: DEVICE, Baud: 9600}
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
		rets, _ := json.MarshalIndent(records, "", "    ")

		now := time.Now().Format("2006-01-02T15:04:05Z")
		fmt.Print(" ", now, " Uploading")
		for i := 0; i < CONNRETRY; i++ {
			err := upload(rets)
			fmt.Print(".")
			if err == nil {
				fmt.Println("Done")
				break
			} else if i == CONNRETRY-1 {
				fmt.Println("Fail")
			}
		}
	}
}
