package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	// "strconv"
	"bytes"
	// "io/ioutil"
	"net/http"
	"strings"
	"time"
)

type record map[string]string

const MAXLENGTH = 500

var head = []string{"temp"}

func getData(reader *bufio.Reader, records *[]record) {
	for i := 1; i < 10; i++ {
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			break
		}

		var record = record{}
		record["time"] = time.Now().Format("2006-01-02T15:04:05")

		line := strings.TrimSpace(string(reply))
		for j, v := range strings.Split(line, " ") {
			record[head[j]] = v
		}
		*records = append(*records, record)
	}
}

func upload(jsonstr []byte) {
	now := time.Now().Format("2006-01-02T15:04:05")
	url := "https://api.myjson.com/bins/3wczx"
	fmt.Print(now, " ", url, " ... ")

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}

func main() {
	c := &serial.Config{Name: DEVICE, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}

	// s := strings.NewReader("10.1\n11.2\n10.1\n11.2\n10.1\n11.2\n")
	reader := bufio.NewReader(s)

	var records []record

	for {
		getData(reader, &records)
		if len(records) > MAXLENGTH {
			records = records[len(records)-MAXLENGTH:]
		}
		rets, _ := json.MarshalIndent(records, "", "    ")
		upload(rets)
	}

	s.Close()
}
