//go:generate go-enum -f=$GOFILE --marshal

package api

// State represents the state of the processor
/*
ENUM(
UNINITIALIZED=1,
CONFIGURED, // processor has been configured successfully
RUNNING, // processor has received and begun processing at least one chunk
SUCCEEDED, // processor has successfully processed all data chunks
COMPLETED, // processor finished with at least one individual chunk failing but no fatal errors
FAILED // processor hit an irrecoverable error and terminated
)
*/
type State int32

// DataChunkState represents the state of one data chunk
/*
ENUM(
WAITING=1, // the processor is waiting for the chunk to arrive, either from an external or regular source
VALIDATING, // the processor is validating that the chunk can work with the current configuration
DETERMINING, // the processor is determining the External State of the chunk
RUNNING, // the processor is running on the chunk
SUCCEEDED, // processor has successfully processed the data chunk
FAILED, // processor hit an error and terminated processing for the chunk.
)
*/
//FATAL // the processor hit a fatal error. This will stop all processing for the run TODO: figure out whether the distinction between regular failure and fatal failure is necessary. For now, just ignore all errors, including fatal
type DataChunkState int32
