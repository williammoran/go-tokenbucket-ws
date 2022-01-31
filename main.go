package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	confFile := flag.String("c", "tokenbucket.conf", "location of config file")
	listen := flag.String("l", ":8080", "listen address")
	flag.Parse()
	f, err := os.Open(*confFile)
	if err != nil {
		panic(err.Error())
	}
	buckets := make(map[string]*bucket)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		init, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		max, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		replenish, err := time.ParseDuration(record[3])
		if err != nil {
			panic(err.Error())
		}
		buckets[record[0]] = NewBucket(init, max, replenish)
	}
	f.Close()
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			p := strings.Split(r.URL.Path, "/")
			amount, err := strconv.ParseInt(p[len(p)-1], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			bucketName := p[len(p)-3]
			tokenName := p[len(p)-2]
			bucket, found := buckets[bucketName]
			if !found {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			use := bucket.Use(tokenName, amount)
			w.Header().Set("Content-Type", "application/json")
			jsonResp, err := json.Marshal(use)
			if err != nil {
				panic(err.Error())
			}
			w.Write(jsonResp)
		},
	)
	http.ListenAndServe(*listen, nil)
}
