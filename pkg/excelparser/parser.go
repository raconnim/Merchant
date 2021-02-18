package excelparser

import (
	"api/pkg/statistics"
	"fmt"

	"github.com/tealeg/xlsx/v3"
)

type ExcelParser struct {
	VendorID int
	WB       *xlsx.File
}

func (ep *ExcelParser) SetData(vendor int, wb *xlsx.File) {
	ep.WB = wb
	ep.VendorID = vendor
}

func (ep *ExcelParser) GetData() (*statistics.StatisticList, error) {

	st := statistics.NewStatisticList()
	//xlsx.
	//wb, err := xlsx.OpenFile(ep.Link)
	// if err != nil {
	// 	return nil, fmt.Errorf("no open file")
	// }

	for _, sh := range ep.WB.Sheets {
		maxRow := sh.MaxRow
		fmt.Println(maxRow)
		for i := 1; i < maxRow; i++ {
			st.GetProduct(sh, i, ep.VendorID)
		}
	}

	return st, nil
}
