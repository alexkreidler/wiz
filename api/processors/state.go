//go:generate go-enum -f=$GOFILE --marshal

package processors

// State represents the state of the processor
/*
ENUM(
UNINITIALIZED=1,
CONFIGURED, // processor has been configured successfully
RUNNING, // processor has received and begun processing at least one chunk
SUCCEEDED, // processor has successfully processed all data chunks
ERRORED, // processor finished with at least one individual chunk failing
)
*/
type State int32

// DataChunkState represents the state of one data chunk
/*
ENUM(
WAITING=1, // the processor is waiting for the chunk to arrive, either from an external or regular source
RUNNING, // the processor is running on the chunk
SUCCEEDED, // processor has successfully processed the data chunk
FAILED, // processor hit an error and terminated processing for the chunk.
)
*/
type DataChunkState int32
