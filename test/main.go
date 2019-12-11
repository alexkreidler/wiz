package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
)

func main() {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", "http://ipv4.download.thinkbroadband.com/1GB.zip")

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.DoBatch(6, req)
	r := (<-resp)
	//r.HTTPResponse.H
	spew.Dump(r.HTTPResponse)
	fmt.Printf("  %v\n", r.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				r.BytesComplete(),
				r.Size,
				100*r.Progress())

		case <-r.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := (<-resp).Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", (<-resp).Filename)

	// Output:
	// Downloading http://www.golang-book.com/public/pdf/gobook.pdf...
	//   200 OK
	//   transferred 42970 / 2893557 bytes (1.49%)
	//   transferred 1207474 / 2893557 bytes (41.73%)
	//   transferred 2758210 / 2893557 bytes (95.32%)
	// Download saved to ./gobook.pdf
}