package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"api/pkg/excelparser"
	"api/pkg/products"
	"api/pkg/statistics"

	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"
)

//go:generate mockgen -source=items.go -destination=items_mock.go -package=handlers ItemRepositoryInterface

type ItemRepositoryInterface interface {
	GetAll() ([]*products.Product, error)
	GetByID(vendorID, offerID int) (*products.Product, error)
	Add(elem *products.Product) (int64, error)
	Update(elem *products.Product) (int64, error)
	Delete(id int, offer_id int) (int64, error)
	GetProduct(vendorID, offerID int, subName string) ([]*products.Product, error)
}

type ProductRepositoryInterface interface {
	GetData() ([]*products.Product, error)
	SetData(vendor int, Link string)
}

type ItemsHandler struct {
	Tmpl        *template.Template
	ProductRepo ProductRepositoryInterface
	ItemsRepo   ItemRepositoryInterface
	ParserRepo  excelparser.ExcelParser
	Logger      *zap.SugaredLogger
}

func (h *ItemsHandler) Show(w http.ResponseWriter, r *http.Request) {

	err := h.Tmpl.ExecuteTemplate(w, "show.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Index(w http.ResponseWriter, r *http.Request) {
	//r.FormFile()
	err := h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		h.Logger.Error("GetData err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "product.html", struct {
		Products []*products.Product
	}{
		Products: elems,
	})
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) ListProduct(w http.ResponseWriter, r *http.Request) {
	subName := r.FormValue("subname")

	value := r.FormValue("vendor")
	vendor, err := strconv.Atoi(value)
	if err != nil {
		h.Logger.Error("field vendor cannot be converted to int", err)
		http.Error(w, `field vendor cannot be converted to int`, http.StatusInternalServerError)
		return
	}

	value = r.FormValue("offerid")
	offerID, err := strconv.Atoi(value)
	if err != nil {
		h.Logger.Error("field offerid cannot be converted to int", err)
		http.Error(w, `field offerid cannot be converted to int`, http.StatusInternalServerError)
		return
	}
	elems, err := h.ItemsRepo.GetProduct(vendor, offerID, subName)
	if err != nil {
		h.Logger.Error("GetData err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "product.html", struct {
		Products []*products.Product
	}{
		Products: elems,
	})
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Upload(w http.ResponseWriter, r *http.Request) {

	err := h.Tmpl.ExecuteTemplate(w, "upload.html", nil)
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}

}

func (h *ItemsHandler) Statistic(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		h.Logger.Error("failed to upload file", err)
		http.Error(w, `failed to upload file`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	value := r.FormValue("id")
	id, err := strconv.Atoi(value)
	if err != nil {
		h.Logger.Error("field Продавец cannot be converted to int", err)
		http.Error(w, `field Продавец cannot be converted to int`, http.StatusInternalServerError)
		return
	}

	wb, err := xlsx.OpenReaderAt(file, 10000)
	if err != nil {
		h.Logger.Error("need xlsx format file", err)
		http.Error(w, `need xlsx format file`, http.StatusInternalServerError)
		return
	}

	h.ParserRepo.SetData(id, wb)
	elems, err := h.ParserRepo.GetData()
	if err != nil {
		h.Logger.Error("GetData err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	st := statistics.Statistic{ErrorRow: elems.ErrorRow}
	for _, row := range elems.AddRows {
		_, err := h.ItemsRepo.GetByID(row.VendorID, row.OfferID)
		if err != nil {
			st.CreateRow++
			_, err := h.ItemsRepo.Add(row)
			if err != nil {
				h.Logger.Error("Add Data err: ", err)
				http.Error(w, `DB err`, http.StatusInternalServerError)
				return
			}
		} else {
			st.UpdateRow++
			h.ItemsRepo.Update(row)
		}
	}
	for _, row := range elems.DropRows {
		_, err = h.ItemsRepo.Delete(row.VendorID, row.OfferID)
		if err == nil {
			st.DropRow++
		}
	}

	err = h.Tmpl.ExecuteTemplate(w, "statistic.html", st)
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}
