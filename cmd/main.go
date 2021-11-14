package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/wcharczuk/go-chart/v2"
)

func main() {
	data := getData()

	dates := getDates(data)

	petrol95Series := chart.TimeSeries{
		Name:    "Petrol (95)",
		XValues: dates,
		YValues: getDataField(data, "petrol95"),
	}

	ethanolSeries := chart.TimeSeries{
		Name:    "Ethanol (e85)",
		XValues: dates,
		YValues: getDataField(data, "ethanol"),
	}

	dieselSeries := chart.TimeSeries{
		Name:    "Diesel (B7)",
		XValues: dates,
		YValues: getDataField(data, "diesel"),
	}

	graph := chart.Chart{
		Title: "Fuel price history in Sweden",
		YAxis: chart.YAxis{
			Name: "Price in SEK inc tax",
		},
		Series: []chart.Series{
			petrol95Series,
			ethanolSeries,
			dieselSeries,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 81,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	file, err := os.Create("output")
	if err != nil {
		panic("Unable to create file")
	}

	err = graph.Render(chart.PNG, file)
	if err != nil {
		panic("Failed to save error")
	}
}

func getDates(points []DataPoint) []time.Time {
	var dates []time.Time

	for _, point := range points {
		parsed, _ := time.Parse(chart.DefaultDateFormat, point.Date)
		dates = append(dates, parsed)
	}

	return dates
}

func getDataField(points []DataPoint, field string) []float64 {
	var result []float64

	for _, point := range points {
		var fieldData float64

		switch field {
		case "petrol95":
			fieldData = point.Petrol95
			break
		case "ethanol":
			fieldData = point.Ethanol
			break
		case "diesel":
			fieldData = point.Diesel
			break
		case "hvo100":
			fieldData = point.HVO100
			break
		}

		result = append(result, fieldData)
	}

	return result
}

type DataPoint struct {
	Date     string  `json:"date"`
	Petrol95 float64 `json:"95"`
	Ethanol  float64 `json:"e85"`
	Diesel   float64 `json:"diesel"`
	HVO100   float64 `json:"hvo100"`
}

func getData() []DataPoint {
	var data []DataPoint

	req, err := http.NewRequest(http.MethodGet, "https://tanka.se/api/prices", nil)
	if err != nil {
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil
	}

	return data
}
