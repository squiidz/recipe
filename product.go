package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GlobalDB is global db pointer
var GlobalDB = newDB("product.db")

// Product is basic structure to display product information
type Product struct {
	ID       uint
	Name     string
	Price    float64
	Quantity string
	*DB
}

func newProduct(id, name, price, quant string) *Product {
	p := Product{}
	if id, err := strconv.Atoi(id); err == nil {
		p.ID = uint(id)
	}
	if price, err := strconv.ParseFloat(price, 64); err == nil {
		p.Price = math.Round(price*100) / 100
	}
	p.Name = name
	p.Quantity = quant
	p.DB = GlobalDB
	return &p
}

// Display print the product
func (p *Product) Display() {
	if strings.Contains(p.Quantity, " g") {
		fmt.Printf("%s (id: %d) are %.2f$ for %s\n", p.Name, p.ID, p.Price, p.Quantity)
		return
	}
	fmt.Printf("%s (id: %d) are %.2f$/unit\n", p.Name, p.ID, p.Price)
}

func (p *Product) save() error {
	if p.ID == 0 {
		return errors.New("Invalid Product")
	}
	tx, err := p.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into product(id, name, price, quantity) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID, p.Name, p.Price, p.Quantity)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// SearchProduct make a request to the metro.ca search url
func SearchProduct(productName string) (string, error) {
	res, err := http.Get(fmt.Sprintf("%ssearch?free-text=%s", baseURL, strings.Replace(productName, " ", "-", -1)))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return extractProductLink(res.Body)
}

// LookupProduct extract the useful information from the product page
func LookupProduct(link string) (*Product, error) {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, link))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ID string
	if splits := strings.Split(link, "/"); len(splits) != 0 {
		ID = splits[len(splits)-1]
	}
	title := doc.Find(".pi--title").Text()
	price, _ := doc.Find(".pi--prices--first-line").Attr("data-main-price")
	weight := doc.Find(".pi--weight").Text()

	p := newProduct(ID, title, price, weight)
	if err = p.save(); err != nil {
		return nil, err
	}
	return p, nil
}

func extractProductLink(searchReader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(searchReader)
	if err != nil {
		return "", err
	}

	link, _ := doc.Find(".product-details-link").Attr("href")
	return link, nil
}
