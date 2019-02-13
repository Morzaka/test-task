package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	URL         = "http://localhost:8080/"
	contentType = "application/json"
)

var (
	quantReq int
)

func main() {

	//threads := runtime.NumCPU()
	//var wg sync.WaitGroup
	//wg.Add(threads)

	done := make(chan bool)
	inputQuantity()
	fmt.Println("You enter number", quantReq)

	// set period of sending events
	ticker := time.NewTicker(30 * time.Microsecond)
	for p := 0; p < runtime.NumCPU(); p++ {
		go func() {
			i := 0
			for t := range ticker.C {
				MakeRequest(i)
				fmt.Println("Tick at", t)
				fmt.Println(runtime.NumGoroutine())
				fmt.Println(runtime.NumCPU())
				i++
			}
			done <- true
		}()
	}

	go func() {
		<- done
	}()

	time.Sleep(60 * time.Second)
	ticker.Stop()
	fmt.Println("Ticker stopped")

	//for i := 0; i < quantReq; i++ {
	//	MakeRequest(i)
	//}
}

func inputQuantity() {
	fmt.Println("Enter the number of quantity requests")
	_, err := fmt.Scanf("%d", &quantReq)
	if err != nil {
		fmt.Println(err)
		fmt.Println(quantReq)
		inputQuantity()
	}
}

func MakeRequest(i int) {

	message := map[string]int{
		"another":   24,
		"Goroutine": 0,
	}

	message["another"] = i
	message["Goroutine"] = runtime.NumGoroutine()

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(URL, contentType, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	fmt.Println("Event NUMBER____", i)
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
}
