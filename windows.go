package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"io/ioutil"

	"fmt"
	"sort"
)

const (
	size     = 5
	interval = 25
)

type DelayCalculator interface {
	AddInterval(value int)
	GetMedian() int
}

type Window struct {
	arr      []int
	size     int
	head     int
	tail     int
	interval int
}

func New(size, interval int) *Window {
	if size < 1 || interval < 1 {
		panic("Must have positive size and interval")
	}
	return &Window{
		arr:      make([]int, size),
		size:     size,
		interval: interval,
	}
}

func (w *Window) AddInterval(interval int) {
	w.interval = interval
}

type Res struct {
	median   string
	position int
}

func (w *Window) GetMedian() int {
	var r int
	a := w.Slice()
	if len(a) == 1 {

		return -1
	}
	//if len is even (M) = value of [((n)/2)th item term + ((n)/2 + 1)th item term ] /2
	if l := len(a); l%2 == 0 {
		var err error
		r, err = evenLen(a[(len(a)/2)-1], a[(len(a)/2)])
		if err != nil {
			panic(err)
		}
	} else {
		//If n is odd then Median (M) = value of ((n + 1)/2)th item from a sorted array of length n.
		sort.Ints(a)
		r = a[(len(a)-1)/2]
	}

	return r
}

func evenLen(a, b int) (median int, err error) {
	median = (a + b) / 2
	return median, nil
}

func main() {
	files, err := ioutil.ReadDir("./getMedian-attachments")
	if err != nil {
		log.Fatal(err)
	}
	//Sliding window size for Test 2 is 100
	//Sliding window size for Test 3 is 1000
	//Sliding window size for Test 4 is 10000
	var sizeStep = 100
	for _, f := range files {
		lines, err := readCSV("./getMedian-attachments/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		start := time.Now()
		win := New(sizeStep, interval)
		var resultArray []string
		for i := 0; i < len(lines); i++ {
			win.Push(lines[i])
			resultArray = append(resultArray, strconv.Itoa(win.GetMedian())+"\r\n")
		}
		endTime := time.Since(start)
		f, err := os.Create(f.Name() + "output.txt")
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range resultArray {
			fmt.Fprintf(f, s)
		}
		log.Println(f.Name()+", \n Slide window size:", sizeStep, " \n Sync Result time: ", endTime)
		sizeStep = sizeStep * 10
	}

}

func readCSV(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return []int{}, err
	}

	reader := csv.NewReader(file)
	var lines []int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err == nil && len(record) > 0 {
			i, err := strconv.Atoi(strings.Trim(record[0], "\r"))
			if err != nil {
				return []int{}, err
			}
			lines = append(lines, i)
		}
	}
	return lines, nil
}

func (w *Window) Push(v int) {
	if w.tail == len(w.arr) {
		w.swap()
	}
	w.arr[w.tail] = v
	if w.tail-w.head >= w.size {
		w.head++
	}
	w.tail++
}

func (w *Window) swap() {
	l := len(w.arr)
	for i := 0; i < w.size-1; i++ {
		w.arr[i] = w.arr[l-w.size+i+1]
	}
	w.head, w.tail = 0, w.size-1
}

func (w *Window) Slice() []int {
	return w.arr[w.head:w.tail]
}
