package statistics

import (
	"api/pkg/products"

	"github.com/tealeg/xlsx/v3"
)

type Statistic struct {
	CreateRow int
	UpdateRow int
	DropRow   int
	ErrorRow  int
}

type StatisticList struct {
	AddRows  []*products.Product
	DropRows []*products.Product
	ErrorRow int
}

func NewStatisticList() *StatisticList {
	return &StatisticList{
		AddRows:  make([]*products.Product, 0, 10),
		DropRows: make([]*products.Product, 0, 10),
	}
}

func (st *StatisticList) GetProduct(sh *xlsx.Sheet, row int, vendor int) {

	item := products.Product{VendorID: vendor}

	offer, err := sh.Cell(row, 0)
	if err != nil {
		st.ErrorRow++
		return
	}
	item.OfferID, err = offer.Int()
	if err != nil || item.OfferID < 0 {
		st.ErrorRow++
		return
	}
	name, err := sh.Cell(row, 1)
	if err != nil {
		st.ErrorRow++
		return
	}
	item.NameProduct = name.Value

	price, err := sh.Cell(row, 2)
	if err != nil {
		st.ErrorRow++
		return
	}
	item.Price, err = price.Float()
	if err != nil || int(item.Price) < 0 {
		st.ErrorRow++
		return
	}

	quantity, err := sh.Cell(row, 3)
	if err != nil {
		st.ErrorRow++
		return
	}
	item.Quantity, err = quantity.Int()
	if err != nil || item.Quantity < 0 {
		st.ErrorRow++
		return
	}

	available, err := sh.Cell(row, 4)
	if err != nil {
		st.ErrorRow++
		return
	}

	if available.Value == "false" {
		st.DropRows = append(st.DropRows, &item)
	} else if available.Value == "true" {
		st.AddRows = append(st.AddRows, &item)
	} else {
		st.ErrorRow++
	}
}
