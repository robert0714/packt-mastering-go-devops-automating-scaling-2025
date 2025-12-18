package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	data := [][]string{
		{"Alice", "30", "85.6"},
		{"Bob", "24", "92.3"},
		{"Charlie", "29", "88.1"},
	}
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader([]string{"Name", "Age", "Score"}),
	)
	table.Bulk(data)
	table.Render()
}
