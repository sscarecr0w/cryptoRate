package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const TOKEN = "" // put here you token from https://www.coingecko.com/en/developers/dashboard

var CryptoIDS = []string{
	"bitcoin",
	"the-open-network",
	"gram-2",
	"ton-fish-memecoin",
}

type Data struct {
	ID         string     `json:"id"`
	Symbol     string     `json:"symbol"`
	Name       string     `json:"name"`
	WebSlug    string     `json:"web_slug"`
	MarketData MarketData `json:"market_data"`
}

type MarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}

func main() {

	coinsInfo := make(chan Data, len(CryptoIDS))

	var wg sync.WaitGroup
	for _, id := range CryptoIDS {
		time.Sleep(time.Millisecond * 200)

		wg.Add(1)

		address := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", id)
		go MakeRequest(coinsInfo, address, &wg)

	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Price"})

	wg.Wait()
	close(coinsInfo)

	for coin := range coinsInfo {
		name := coin.Name + "(" + coin.Symbol + ")"
		price := fmt.Sprintf("$%v", coin.MarketData.CurrentPrice["usd"])
		t.AppendRow(table.Row{name, price})
	}

	fmt.Println(t.Render())

}

func MakeRequest(coinsInfo chan Data, address string, wg *sync.WaitGroup) {

	defer wg.Done()

	client := http.DefaultClient

	u, err := url.Parse(address)

	if err != nil {
		log.Fatal(err)
	}

	headers := http.Header{}
	headers.Set("Accepts", "application/json")
	headers.Add("x-cg-api-key", TOKEN)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: headers,
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	var coin Data

	err = json.Unmarshal(data, &coin)

	if err != nil {
		log.Fatal(err)
	}

	coinsInfo <- coin
}
