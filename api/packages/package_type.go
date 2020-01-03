package packages
//go:generate go-enum -f=$GOFILE --marshal

// PackageType represents the type of data chunk, either input or output
/*
ENUM(
DATA=1,
CODE,
)
*/
type PackageType int32
