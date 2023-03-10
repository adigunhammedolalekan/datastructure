package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type (
	Product struct {
		Client   string `json:"client"`
		Quantity int    `json:"quantity"`
		Color    string `json:"color"`
	}
)

func main() {

	var filePath = ""
	flag.StringVar(&filePath, "path", "test.txt", "go run main.go -path=test.txt")
	flag.Parse()

	data, err := read(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s - %v", filePath, err)
	}

	products := make(map[int][]Product)
	keys := make([]int, 0)
	for _, v := range data {
		values := strings.Split(v, ",")
		if len(values) >= 3 {
			price := toInt(values[2])
			keys = append(keys, price)
			product := Product{Client: values[0], Quantity: toInt(values[3]), Color: values[1]}
			items, ok := products[price]
			if !ok {
				items = make([]Product, 0)
			}

			items = append(items, product)
			products[price] = items
		}
	}

	// sort map
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	// use orderedmap because Go native map does not support ordering
	newValues := orderedmap.New()
	for _, k := range keys {
		newValues.Set(fmt.Sprintf("%d", k), products[k])
	}
	// encode to JSON and print to std output
	if err := json.NewEncoder(os.Stdout).Encode(newValues); err != nil {
		log.Fatal(err)
	}
}

func toInt(s string) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return 0
}

func read(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var response []string
	reader := bufio.NewScanner(f)

	for reader.Scan() {
		response = append(response, reader.Text())
	}

	return response, nil
}
