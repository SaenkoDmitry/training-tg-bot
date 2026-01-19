package docgenerator

import (
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

func AutoFitColumns(
	f *excelize.File,
	sheet string,
	fromCol, toCol int,
) {
	for col := fromCol; col <= toCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		maxLen := 0

		rows, _ := f.GetRows(sheet)
		for _, row := range rows {
			if col-1 < len(row) {
				l := utf8.RuneCountInString(row[col-1])
				if l > maxLen {
					maxLen = l
				}
			}
		}

		width := float64(maxLen) + 2
		if width > 60 {
			width = 60
		}

		_ = f.SetColWidth(sheet, colName, colName, width)
	}
}
