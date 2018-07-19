package main

import (
	"testing"
)

var ricardoTestURL = "https://www.ricardocuisine.com/recettes/5582-poulet-grille-et-gremolata-au-celeri-et-a-la-noix-de-coco"

func TestParsingRicardoRecipe(t *testing.T) {
	recipe := NewRecipe(ricardoTestURL)
	if len(recipe.Ingredients) == 0 {
		t.Error("Parsing ricardo recipe failed")
	}
	"		" recipe.Display()
}
