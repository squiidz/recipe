package main

import (
	"log"
	"os"
	"sync"
)

var baseURL = "https://www.metro.ca/"
var productURL = "https://www.metro.ca/epicerie-en-ligne/allees/p/"

var testLink = "https://www.metro.ca/epicerie-en-ligne/allees/produits-laitiers-et-fromages/mon-fromager/fromages-a-pate-molle-et-frais/fromage-camembert/p/3161910238710"
var testRecipe = "https://www.ricardocuisine.com/recettes/5582-poulet-grille-et-gremolata-au-celeri-et-a-la-noix-de-coco"

func main() {
	arg := os.Args[len(os.Args)-1]
	recipe := NewRecipe(arg)
	var wg sync.WaitGroup
	for _, ing := range recipe.Ingredients {
		wg.Add(1)
		go func(ing Ingredient) {
			defer wg.Done()
			productName := processTerm(ing.raw)
			log.Println("Looking for", productName)
			if trsl := FindTerm(productName); trsl != nil {
				if product, err := GlobalDB.findLike(trsl.ValueTerm); err == nil {
					product.Display()
					return
				}
			}

			link, err := SearchProduct(productName)
			if err != nil {
				log.Println(err)
			}
			product, err := LookupProduct(link)
			if err != nil {
				log.Println(err)
				return
			}
			err = SaveTerm(productName, product.Name)
			if err != nil {
				log.Println(err)
			}
			product.Display()
		}(ing)
	}
	wg.Wait()
}
