package main

import "strings"

// Translation is the search/value pair
type Translation struct {
	SearchTerm string
	ValueTerm  string
}

// SaveTerm add a term with is value
func SaveTerm(st, vt string) error {
	ct := processTerm(st)
	stmt, err := GlobalDB.Prepare("INSERT INTO translation(search, value) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ct, vt)
	if err != nil {
		return err
	}
	return nil
}

// FindTerm search for value of the search term
func FindTerm(st string) *Translation {
	stmt, err := GlobalDB.Prepare("SELECT value FROM translation WHERE search like ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	var trsl Translation
	err = stmt.QueryRow(st).Scan(&trsl.ValueTerm)
	if err != nil {
		return nil
	}
	trsl.SearchTerm = st
	return &trsl
}

func processTerm(st string) string {
	if strings.Contains(st, ")") {
		ss := strings.SplitAfter(st, ") de ")
		return ss[len(ss)-1]
	}
	return st
}
