package executor

import (
	"github.com/alexkreidler/wiz/api/processors"
	"github.com/alexkreidler/wiz/client"
	"github.com/alexkreidler/wiz/utils"
	"github.com/imdario/mergo"
	"log"
	"sync/atomic"
)

// createOutputDataChunk simply creates a default output data chunk given an input data chunk
func createOutputDataChunk(d processors.Data) processors.Data {
	chunk := processors.Data{
		ChunkID:             d.AssociatedChunkID,
		Format:              0, // the format should change, it should be set by the output of the processor
		Type:                processors.DataTypeOUTPUT,
		State:               processors.DataChunkStateWAITING,
		RawData:             nil,
		FilesystemReference: processors.FilesystemReference{},
		AssociatedChunkID:   d.ChunkID,
	}

	return chunk
}

//rManager is a goroutine that transforms updates on the Processor's state into updates in the Executor's memory which can then respond to HTTP requests
// TODO: reevaluate later for too much locking?
// data is the data to run on, r is the runProcessor object, and baseProcessor is the configured ChunkProcessor that will be the base for the new spawned processor
func rManager(data processors.Data, r *runProcessor) {
	log.Println("starting runManger")
	outputChunk := createOutputDataChunk(data)

	r.dataLock.Lock()
	if r.workers == nil {
		r.workers = make(map[string]*Worker)
	}

	w := Worker{
		p:   r.baseProcessor.New(),
		in:  data,
		out: outputChunk,
	}

	r.workers[data.ChunkID] = &w

	r.dataLock.Unlock()

	// Handle the different data formats: TODO figure out others
	switch data.Format {
	case processors.DataFormatRAW:
		go w.p.Run(data.RawData)
		break
	default:
		log.Println("unsupported data type")
		return
	}
	log.Println("Processor has started")

	// handle state updates from the channel
	// this range stmt requires that the processor close its channel
	for state := range w.p.State() {
		log.Println("state change", state)
		r.dataLock.Lock()
		r.workers[data.ChunkID].in.State = state
		r.dataLock.Unlock()
	}
	// handle completion of this chunk
	out := w.p.Output()
	log.Println("chunk", w.in.ChunkID, "has completed. Got output:", out)

	r.dataLock.Lock()
	r.workers[data.ChunkID].out.Format = processors.DataFormatRAW
	r.workers[data.ChunkID].out.State = r.workers[data.ChunkID].in.State
	r.workers[data.ChunkID].out.RawData = out
	prevDataChunk := r.workers[data.ChunkID].out
	r.dataLock.Unlock()

	if r.run.Configuration.ExecutorConfig.SendDownstream {
		for _, downstream := range r.run.Configuration.ExecutorConfig.DownstreamLocations {
			err := sendToDownstream(prevDataChunk, downstream)
			log.Printf("failed to send chunk %s downstream: %v", prevDataChunk.ChunkID, err)
		}
	}
	atomic.AddUint32(&r.numCompleted, 1)
	if r.numCompleted == uint32(r.run.Configuration.ExpectedData.NumChunks) {
		handleAllChunksCompleted(r)
	}
}

func sendToDownstream(prevDataChunk processors.Data, downstream processors.DownstreamDataLocation) error {
	c := client.NewClient(downstream.Hostname)

	newOutputID := utils.GenID()
	// The previous output now becomes the input, so we need to generate a new output ID
	// In the future we may add provenance builtin to record every chunk that it has progressed from
	newDataChunk := processors.Data{Type: processors.DataTypeINPUT, State: processors.DataChunkStateWAITING, AssociatedChunkID: newOutputID}

	err := mergo.Merge(&newDataChunk, prevDataChunk)
	if err != nil {
		return err
	}
	return c.AddData(downstream.ProcID, downstream.ProcID, newDataChunk)
}

func setRunState(r *runProcessor, state processors.State) {
	r.runLock.Lock()
	r.run.CurrentState = state
	r.runLock.Unlock()
}

//handleAllChunksCompleted gets called whenever a processor completes (aka all chunks are done) to recalculate the Run state
func handleAllChunksCompleted(r *runProcessor) {
	r.dataLock.RLock()
	for _, v := range r.workers {
		if v.in.State == processors.DataChunkStateFAILED {
			setRunState(r, processors.StateERRORED)
			return
		} else if v.in.State != processors.DataChunkStateSUCCEEDED {
			log.Println("error: expected all chunks to be succeeded, but", v.in.ChunkID, "is not")
		}
	}
	r.dataLock.RUnlock()

	log.Print("all chunks succeeded")
	setRunState(r, processors.StateSUCCEEDED)
}
