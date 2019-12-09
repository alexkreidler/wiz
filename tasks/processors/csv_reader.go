package processors

import (
	"encoding/csv"
	"encoding/json"
	"github.com/alexkreidler/ratchet/data"
	"github.com/alexkreidler/ratchet/util"
	"os"
)

// CSVReader opens and reads the contents of the given filename.
type CSVReader struct {
	filename string
}

// NewFileReader returns a new CSVReader that will read the entire contents
// of the given file path and send it at once. For buffered or line-by-line
// reading try using IoReader.
func NewCSVReader(filename string) *CSVReader {
	return &CSVReader{filename: filename}
}

// ProcessData reads a file and sends its contents to outputChan
func (r *CSVReader) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	f, err := os.Open(r.filename)
	util.KillPipelineIfErr(err, killChan)

	csvReader := csv.NewReader(f)
	csvData, err := csvReader.ReadAll()
	util.KillPipelineIfErr(err, killChan)

	d, err = json.Marshal(csvData)

	util.KillPipelineIfErr(err, killChan)
	outputChan <- d
}

// Finish - see interface for documentation.
func (r *CSVReader) Finish(outputChan chan data.JSON, killChan chan error) {
}

func (r *CSVReader) String() string {
	return "CSVReader"
}

