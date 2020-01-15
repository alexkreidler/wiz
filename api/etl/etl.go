package etl

type TextView interface {
	// View selects the lines from from to to and returns them as a string
	View(from int, to int) string
}

type TableView interface {
	// takes the row and cols from the from and to cells and returns an excel-style rectangular subset of the table
	View(from Cell, to Cell) Table
}


// TODO: think about these serialized data representations and if they are the most effective way to go across the wire to React frontend

// A cell represents a cell that can be accessed or set
type Cell struct {
	Row int
	Col int
	Val interface{}
}


// A table is a map from Row indexes to Column indexes to cell values
type Table map[int]map[int]Cell
