package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	Ingredients []Ingredient
}

func (r *Recipe) Display() {
	for _, ing := range r.Ingredients {
		fmt.Println(ing.raw)
	}
}

type Ingredient struct {
	raw string
}

func NewRecipe(url string) *Recipe {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	if strings.Contains(url, "ricardocuisine.com") {
		return parseRicardoRecipe(res.Body)
	} else if strings.Contains(url, "recettes.qc.ca") {
		return parseRecetteDuQuebecRecipe(res.Body)
	}
	return nil
}

func parseRicardoRecipe(rdr io.Reader) *Recipe {
	var ingredients []Ingredient

	doc, err := goquery.NewDocumentFromReader(rdr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Looking for the ingredients section")
	doc.Find("#ingredients").Each(func(_ int, s *goquery.Selection) {
		s.Find("#formIngredients").Each(func(_ int, sl *goquery.Selection) {
			sl.Find("ul").Each(func(_ int, sel *goquery.Selection) {
				sel.Find("li").Each(func(_ int, sele *goquery.Selection) {
					ingredient := strings.Replace(sele.Find("label").Text(), "	", " ", -1)
					ingredients = append(ingredients, Ingredient{raw: ingredient})
				})
			})
		})
	})

	return &Recipe{Ingredients: ingredients}
}

func parseRecetteDuQuebecRecipe(rdr io.Reader) *Recipe { return nil }
