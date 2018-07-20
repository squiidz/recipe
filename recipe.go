package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	prose "gopkg.in/jdkato/prose.v2"
)

type Recipe struct {
	Ingredients []Ingredient
}

func (r *Recipe) Display() {
	for _, ing := range r.Ingredients {
		fmt.Println(ing.Name)
	}
}

type Ingredient struct {
	raw  string
	Name string
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
					ing := Ingredient{raw: ingredient}
					ing.extractIngredientName()
					ingredients = append(ingredients, ing)
				})
			})
		})
	})

	return &Recipe{Ingredients: ingredients}
}

func parseRecetteDuQuebecRecipe(rdr io.Reader) *Recipe { return nil }

func (ing *Ingredient) extractIngredientName() {
	// Create a new document with the default configuration:
	doc, err := prose.NewDocument(ing.raw)
	if err != nil {
		log.Fatal(err)
	}

	var nouns []prose.Token
	var adjs []prose.Token
	for _, tok := range doc.Tokens() {
		if tok.Tag == "NN" || tok.Tag == "NNS" {
			nouns = append(nouns, tok)
		} else if tok.Tag == "JJ" {
			adjs = append(adjs, tok)
		}
	}

	if len(adjs) > 0 {
		ing.Name = fmt.Sprint(adjs[len(adjs)-1].Text, " ", nouns[len(nouns)-1].Text)
	} else {
		ing.Name = fmt.Sprint(nouns[len(nouns)-2].Text, " ", nouns[len(nouns)-1].Text)
	}
}
