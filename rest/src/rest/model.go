package rest

import (
	"database/sql"
	"log"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {
	err := db.QueryRow("SELECT name, price FROM products WHERE id=?", p.ID).Scan(&p.Name, &p.Price)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=?, price=? WHERE id=?", p.Name, p.Price, p.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=?", p.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT products SET name=?, price=?")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = stmt.Exec(p.Name, p.Price)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func getProducts(db *sql.DB, count int) ([]product, error) {
	rows, err := db.Query("SELECT id, name,  price FROM products LIMIT ?", count)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
