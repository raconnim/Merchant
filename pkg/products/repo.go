package products

import (
	"database/sql"
	"fmt"
	"strings"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{DB: db}
}

func (repo *ItemRepository) GetAll() ([]*Product, error) {
	items := []*Product{}
	rows, err := repo.DB.Query("SELECT vendor_id, offer_id, name, price, quantity FROM product")
	if err != nil {
		return nil, fmt.Errorf("----%v", err)
	}
	defer rows.Close() // надо закрывать соединение, иначе будет течь
	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.VendorID, &product.OfferID, &product.NameProduct, &product.Price, &product.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}

func (repo *ItemRepository) GetByID(vendorID, offerID int) (*Product, error) {
	product := &Product{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow("SELECT vendor_id, offer_id, name, price, quantity FROM product WHERE vendor_id = $1 and offer_id = $2", vendorID, offerID).
		Scan(&product.VendorID, &product.OfferID, &product.NameProduct, &product.Price, &product.Quantity)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *ItemRepository) GetProduct(vendorID, offerID int, subName string) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := repo.DB.Query("SELECT vendor_id, offer_id, name, price, quantity FROM product "+
		"WHERE vendor_id = $1 and offer_id = $2", vendorID, offerID)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.VendorID, &product.OfferID, &product.NameProduct, &product.Price, &product.Quantity)
		if err != nil {
			return nil, err
		}
		if strings.Contains(product.NameProduct, subName) {
			items = append(items, product)
		}
	}
	return items, nil
}

func (repo *ItemRepository) Add(elem *Product) (int64, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO product "+
			"(vendor_id, offer_id, name, price, quantity)"+
			"VALUES ($1, $2, $3, $4, $5)",
		elem.VendorID,
		elem.OfferID,
		elem.NameProduct,
		elem.Price,
		elem.Quantity,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

//result, err := db.Exec("update Products set price = $1 where id = $2", 69000, 1)
func (repo *ItemRepository) Update(elem *Product) (int64, error) {
	result, err := repo.DB.Exec(
		"UPDATE items SET "+
			"name = $1"+
			", price = $2"+
			", quantity = $3"+
			"WHERE vendor_id = $4 and offer_id = $5",
		elem.NameProduct,
		elem.Price,
		elem.Quantity,
		elem.VendorID,
		elem.OfferID,
	)
	if err != nil {
		return 0, fmt.Errorf("------%v", err)
	}
	return result.RowsAffected()
}

func (repo *ItemRepository) Delete(id int, offer_id int) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM product WHERE vendor_id = $1 and offer_id = $2",
		id,
		offer_id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
