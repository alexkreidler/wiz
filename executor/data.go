package executor

import (
	"github.com/alexkreidler/wiz/api"
	"log"
	"sync/atomic"
)

// createOutputDataChunk simply creates a default output data chunk given an input data chunk
func createOutputDataChunk(d api.Data) api.Data {
	chunk := api.Data{
		ChunkID:             d.AssociatedChunkID,
		Format:              0, // the format should change, it should be set by the output of the processor
		Type:                api.DataTypeOUTPUT,
		State:               api.DataChunkStateWAITING,
		RawData:             nil,
		FilesystemReference: api.FilesystemReference{},
		AssociatedChunkID:   d.ChunkID,
	}

	return chunk
}

//rManager is a goroutine that transforms updates on the Processor's state into updates in the Executor's memory which can then respond to HTTP requests
// TODO: reevaluate later for too much locking?
// data is the data to run on, r is the runProcessor object, and baseProcessor is the configured ChunkProcessor that will be the base for the new spawned processor
func rManager(data api.Data, r *runProcessor) {
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
	case api.DataFormatRAW:
		go w.p.Run(data)
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
	r.workers[data.ChunkID].out.Format = api.DataFormatRAW
	r.workers[data.ChunkID].out.State = api.DataChunkStateSUCCEEDED
	r.workers[data.ChunkID].out.RawData = out
	r.dataLock.Unlock()

	if r.run.Configuration.ExecutorConfig.SendDownstream {
	//	todo: send the output to the downstream processors
	}
	atomic.AddUint32(&r.numCompleted, 1)
	if r.numCompleted == uint32(r.run.Configuration.ExpectedData.NumChunks) {
		handleAllChunksCompleted(r)
	}
}

func setRunState(r *runProcessor, state api.State) {
	r.runLock.Lock()
	r.run.CurrentState = state
	r.runLock.Unlock()
}

//handleAllChunksCompleted gets called whenever a processor completes (aka all chunks are done) to recalculate the Run state
func handleAllChunksCompleted(r *runProcessor) {
	r.dataLock.RLock()
	for _, v := range r.workers {
		if v.in.State == api.DataChunkStateFAILED {
			setRunState(r, api.StateERRORED)
			return
		} else if v.in.State != api.DataChunkStateSUCCEEDED {
			log.Println("error: expected all chunks to be succeeded, but", v.in.ChunkID, "is not")
		}
	}
	r.dataLock.RUnlock()

	log.Print("all chunks succeeded")
	setRunState(r, api.StateSUCCEEDED)
}

