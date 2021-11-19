package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// run as :  vmstat 1 | go run main.go
// see cmd arguments as: go run main.go -h

func main() {
	reader := bufio.NewReader(os.Stdin)

	watch := flag.String("w", "free", "watch field")
	threshold := flag.Int("t", 194700, "max threshold")
	maxThresholdCount := flag.Int("c", 3, "max threshold count to show panic")
	flag.Parse()

	freeAlarm := NewAlarm(*watch, *threshold, *maxThresholdCount)
	freeAlarm.PrintInfo()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("some error occurred when reading line:", err)
			continue
		}

		if err := freeAlarm.Check(string(line)); err != nil {
			fmt.Println(err)
		}
	}
}

//alarm file

type Alarm struct {
	watch            string
	threshold        int
	maxThreshold     int
	thresholdCounter int
	t                time.Time
	watchIndex       int
}

func NewAlarm(watch string, threshold, maxThreshold int) *Alarm {
	return &Alarm{
		watch:            watch,
		threshold:        threshold,
		maxThreshold:     maxThreshold,
		thresholdCounter: 0,
		t:                time.Now(),
		watchIndex:       -1, //initial index not found
	}
}

func (a *Alarm) PrintInfo() {
	fmt.Printf("Alert watching for: %s, threshold: %d, max count: %d\n\n", a.watch, a.threshold, a.maxThreshold)
}

func (a *Alarm) Check(line string) error {
	data := strings.Trim(string(line), " ")
	if strings.HasPrefix(data, "procs") { //skip head line
		return nil
	} else if strings.HasPrefix(data, "r  b") { // set watch index
		arr := parseString(data)
		for ix, d := range arr {
			if d == a.watch {
				a.watchIndex = ix
				return nil
			}
		}
	}

	if a.watchIndex == -1 {
		panic("invalid alarm for " + a.watch)
	}
	return a.process(data)
}

func (a *Alarm) process(data string) error {
	arr := parseString(data)
	if len(arr) <= a.watchIndex {
		return errors.New("some error occurred when processing data:" + data)
	}

	d, err := strconv.Atoi(arr[a.watchIndex])
	if err != nil {
		return errors.New("some error occurred when converting data:" + err.Error())
	}

	if d >= a.threshold {
		a.thresholdCounter++
	} else {
		a.thresholdCounter = 0
	}

	switch {
	case a.thresholdCounter == 1:
		a.t = time.Now()
	case a.thresholdCounter >= a.maxThreshold:
		a.panic(d)
	}
	return nil
}

func (a *Alarm) panic(data int) {
	fmt.Printf("PANIC: '%s' threshold value reached at %d from last %v with max %d attempts\n", a.watch, data, time.Since(a.t), a.thresholdCounter)
}

func parseString(data string) []string {
	result := []string{}
	arr := strings.Split(data, " ")
	for _, d := range arr {
		if d != "" { //cleanup the empty strings
			result = append(result, d)
		}
	}
	return result
}
