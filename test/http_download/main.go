package main

import (
	"io"
	"log"
	"os"

	download "github.com/joeybloggs/go-download"
)

func main() {

	// no options specified so will default to 10 concurrent download by default

	url := "https://download.bls.gov/pub/time.series/la/la.data.64.County"
	f, err := download.Open(url, &download.Options{
		Concurrency: func(size int64) int {
			return 0
		},
		Proxy:       nil,
		Client:      nil,
		Request:     nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	freal, err := os.Create("data.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer freal.Close()

	piper, pipew := io.Pipe()

	// write in writer end of pipe
	go func() {
		defer pipew.Close()
		io.Copy(pipew, f)
	}()

	// read from reader end of pipe.
	io.Copy(freal, piper)
	piper.Close()

	// f implements io.Reader, write file somewhere or do some other sort of work with it
}
