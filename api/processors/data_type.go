//go:generate go-enum -f=$GOFILE --marshal

package processors

// DataFormat represents the format of data chunk, either raw JSON data or a filesystem reference
/*
ENUM(
RAW=1, // raw generic data
FILESYSTEM_REF, // a reference to a file or folder accessible to the processor
)
*/
type DataFormat int32

// DataType represents the type of data chunk, either input or output
/*
ENUM(
INPUT=1,
OUTPUT,
)
*/
type DataType int32
