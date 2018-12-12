package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"text/template"
	"time"
)

type Conf struct {
	TimeGenerated string // Time a Configuration was Created
	Zoos          []Zoo  // List of Zoos
}

type Zoo struct {
	Name    string   // Name of Zoo
	Climate string   // Climate of Area where Zoo is located
	Animals []Animal // Animals in Zoo
}

type Animal struct {
	Name     string   // Name of the Animal
	Climates []string // Climates where it can live
}

var (
	// Animals
	Alligator = Animal{Name: "Alligator", Climates: []string{"Tropical", "SubTropical"}}
	Crocodile = Animal{Name: "Crocodile", Climates: []string{"Tropical", "SubTropical"}}
	ArcticFox = Animal{Name: "ArcticFox", Climates: []string{"Arctic", "SubArctic"}}
	Puffin    = Animal{Name: "Puffin", Climates: []string{"Arctic", "SubArctic"}}

	// Zoos
	miamiZoo     = Zoo{Name: "MiamiZoo", Climate: "SubTropical", Animals: []Animal{Alligator, Crocodile, ArcticFox, Puffin}}
	reykjavikZoo = Zoo{Name: "ReykjavikZoo", Climate: "SubArctic", Animals: []Animal{Alligator, Crocodile, ArcticFox, Puffin}}

	// Function Map, notice how we define the getAcceptableAnimals key to a function of the same name
	funcMap = template.FuncMap{"getAcceptableAnimals": getAcceptableAnimals, "getUnacceptableAnimals": getUnacceptableAnimals}
)

func main() {
	conf := Conf{TimeGenerated: time.Now().UTC().String(), Zoos: []Zoo{miamiZoo, reykjavikZoo}}

	// Gets the directory this file is in
	// See: https://stackoverflow.com/questions/32163425/golang-how-to-get-the-directory-of-the-package-the-file-is-in-not-the-current-w
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	// The template name must match a file in ParseFiles which is the main template
	// See: https://stackoverflow.com/questions/10199219/go-template-function
	tmpl, err := template.New("zoo.tmpl").Funcs(funcMap).ParseFiles(path.Dir(filename) + "/zoo.tmpl")
	if err != nil {
		panic(err)
	}

	// Here we use the template and conf to make to generate textual output
	// We are using 'os.Stdout` to output to screen, a file can be used instead
	err = tmpl.Execute(os.Stdout, conf)
	if err != nil {
		panic(err)
	}

}

// Function used in the template to return a list of acceptable animals
func getAcceptableAnimals(a interface{}, b string) []Animal {
	animals, ok := a.([]Animal)
	if !ok {
		err := errors.New(fmt.Sprintf("expected an '[]*Animal' type but %T was returned", animals))
		panic(err)
	}

	acceptable := []Animal{}
	for _, animal := range animals {
		for _, climate := range animal.Climates {
			if b == climate {
				acceptable = append(acceptable, animal)
				break
			}
		}
	}

	return acceptable
}

// Function used in the template to return a list of unacceptable animals
func getUnacceptableAnimals(a interface{}, b string) []Animal {
	animals, ok := a.([]Animal)
	if !ok {
		err := errors.New(fmt.Sprintf("expected an '[]Animal' type but %T was returned", animals))
		panic(err)
	}

	unacceptable := []Animal{}
	for _, animal := range animals {
		for j, climate := range animal.Climates {
			if b == climate {
				break
			} else if j == (len(animal.Climates) - 1) { // If at last item and still not found
				unacceptable = append(unacceptable, animal)
			}
		}
	}

	return unacceptable
}