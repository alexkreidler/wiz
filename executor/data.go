package executor

import (
	"github.com/alexkreidler/wiz/api"
	procApi "github.com/alexkreidler/wiz/processors/processor"
	"log"
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

// data is the data to run on, r is the runProcessor object, and baseProcessor is the configured ChunkProcessor that will be the base for the new spawned processor
func rManager(data api.Data, r *runProcessor, baseProcessor procApi.ChunkProcessor) {

	outputChunk := createOutputDataChunk(data)

	r.dataLock.Lock()
	if r.workers == nil {
		r.workers = make(map[string]*Worker)
	}

	w := Worker{
		p:   baseProcessor.New(),
		in:  data,
		out: outputChunk,
	}

	r.workers[data.ChunkID] = &w

	r.dataLock.Unlock()

	// Handle the different data formats: TODO figure out others
	switch data.Format {
	case api.DataFormatRAW:
		r.wg.Add(1)
		go w.p.Run(data)
		break
	default:
		log.Println("unsupported data type")
		return
	}

	// handle state updates from the channel
	// this range stmt requires that the processor close its channel
	for state := range w.p.State() {
		r.dataLock.Lock()
		r.workers[data.ChunkID].in.State = state
		r.dataLock.Unlock()
	}
	// even if all processors report state Succeeded, they may not exit their goroutines, so we use a WaitGroup. TODO: think about this
	r.wg.Wait()
	handleProcDone(r)
}

//handleProcDone gets called whenever a processor completes to recalculate the Run state
func handleProcDone(r *runProcessor) {
	//r.workers.
	r.dataLock.RLock()
	for k, v := range r.workers {
		if v.in.State == api.DataChunkStateFAILED {
			
		}
	}
}

