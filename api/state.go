//go:generate go-enum -f=$GOFILE --marshal

package api

// State represents the state of the processor
/*
ENUM(
UNINITIALIZED=1,
CONFIGURED, // processor has been configured successfully
RUNNING, // processor has recieved and begun processing at least one chunk
SUCCESS, // processor has successfully processed all data chunks
COMPLETED, // processor finished with at least one individual chunk failing but no fatal errors
FAILURE // processor hit an irrecoverable error and terminated
)
*/
type State int32

// DataChunkState represents the state of one data chunk
/*
ENUM(
VALIDATING, // the processor is validating that the chunk can work with the current configuration
DETERMINING=1, // the processor is determining the External State of the chunk
RUNNING // the processor is running on the chunk
)
*/
type DataChunkState int32
