//go:generate go-enum -f=$GOFILE --marshal

package processor

// DataType represents the type of data chunk
/*
ENUM(
RAW=1, // raw generic data
FILESYSTEM_REF, // a reference to a file or folder accessible to the processor
)
*/
type DataType int32