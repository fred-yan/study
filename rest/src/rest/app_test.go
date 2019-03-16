package rest

import (
	"testing"
	"log"
	"os"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		"root",
		"redis",
		"testdb")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Println(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test product", "price":11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 0.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 0.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)
	res, err := a.DB.Exec("INSERT INTO products(name, price) VALUES(?, ?)", "Product 2", 22.00)
	if err != nil {
		t.Errorf(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Errorf(err.Error())
	}

	req, _ := http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO products(name, price) VALUES(?, ?)",
			"Product "+strconv.Itoa(i), (i+1.0)*10)
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	res, err := a.DB.Exec("INSERT INTO products(name, price) VALUES(?, ?)", "Product 2", 22.00)
	if err != nil {
		t.Errorf(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Errorf(err.Error())
	}

	req, _ := http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	payload := []byte(`{"name":"test product - updated name","price":11.22}`)

	req, _ = http.NewRequest("PUT", "/product/"+strconv.FormatInt(id, 10), bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	res, err := a.DB.Exec("INSERT INTO products(name, price) VALUES(?, ?)", "Product 2", 22.00)
	if err != nil {
		t.Errorf(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Errorf(err.Error())
	}

	req, _ := http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
