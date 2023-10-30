package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url string
	image string
	name string
	price string
}

func main() {
	
	var pokemonProducts []PokemonProduct

	c := colly.NewCollector()

	//c.Visit("https://scrapeme.live/shop/") 

	var wg sync.WaitGroup

	c.OnHTML("li.product", func (e *colly.HTMLElement)  {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = e.ChildAttr("a","href")
		pokemonProduct.image = e.ChildAttr("img","src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func (r *colly.Response)  {
		//Now the scrapping of the data is done so we can now write it in csv file
		defer wg.Done()
	})

	wg.Add(1)

	c.Visit("https://scrapeme.live/shop/") 

	wg.Wait()

	file, err := os.Create("products.csv")

	if err != nil {
		log.Fatalln("Failed to create the output csv file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	header := []string{
		"url",
		"iamge",
		"name",
		"price",
	}

	writer.Write(header)

	for _, pokemonProduct := range pokemonProducts {
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		writer.Write(record)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}