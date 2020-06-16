package util

import "github.com/olekukonko/tablewriter"

// Table - A wrapper around tablewriter.Table, for behavior customization
type Table struct {
	*tablewriter.Table
}

// NewTable - Constructor method with default behavior
func NewTable() *Table {
	return &Table{}
}

// Output - Render the table
func (t *Table) Output() {

}

// SetColumns - Set the headers (and their widths) for a table
func (t *Table) SetColumns(names []string, widths []int) {

}
