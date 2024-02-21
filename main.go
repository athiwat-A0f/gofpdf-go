package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// strDelimit converts 'ABCDEFG' to, for example, 'A,BCD,EFG'
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	// Set current font directory
	// pdf.SetFontLocation("fonts/")
	// Add font based on filename
	// pdf.AddFont("Prompt", "", "Prompt-Regular.ttf")
	pdf.AddUTF8Font("prompt", "", "fonts/Prompt-Regular.ttf")
	pdf.AddUTF8Font("prompt-bold", "B", "fonts/Prompt-Bold.ttf")

	type countryType struct {
		districtStr, cityThStr, cityEnStr, popStr string
	}
	countryList := make([]countryType, 0, 8)
	// count := 0
	header := []string{"Order", "District (EN)", "District (TH)", "Province", "Po."}
	loadData := func(fileStr string) {
		fl, err := os.Open(fileStr)
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var c countryType
			for scanner.Scan() {
				// Austria;Vienna;83859;8075
				lineStr := scanner.Text()
				// print(lineStr)
				list := strings.Split(lineStr, ";")
				if len(list) == 4 {
					c.districtStr = list[0]
					c.cityThStr = list[1]
					c.cityEnStr = list[2]
					c.popStr = list[3]
					countryList = append(countryList, c)

				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			// fmt.Println(countryList)
			// count = len(countryList)
			if len(countryList) == 0 {
				err = fmt.Errorf("error loading data from %s", fileStr)
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}

	// Simple table
	basicTable := func() {
		// pages := math.Ceil(float64(count) / 40)
		// println(count)
		// println(pages)

		i := 0
		j := 0
		rows := 30
		left := (210.0 - 5*40) / 2
		for index, c := range countryList {
			// println(index)
			if index%rows == 0 {
				pdf.SetTopMargin(30)
				pdf.SetHeaderFuncMode(func() {
					pdf.Image("map.jpg", 5, 5, 20, 0, false, "", 0, "")
					pdf.SetY(5)
					pdf.SetFont("prompt-bold", "B", 14)
					// pdf.Cell(150, 10, "District In Thailand")
					pdf.CellFormat(0, 10, "District In Thailand", "", 0, "C", false, 0, "")
					pdf.Ln(20)
				}, true)

				pdf.SetFooterFunc(func() {
					pdf.SetY(-15)
					pdf.SetFont("prompt", "", 8)
					pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
						"", 0, "C", false, 0, "")
				})

				pdf.AddPage()
				pdf.SetX(left)

				for index2, str := range header {
					pdf.SetFont("prompt-bold", "B", 10)
					// print("index=", index2)
					if index2 == 0 || index2 == 4 {
						pdf.CellFormat(20, 7, str, "1", 0, "", false, 0, "")
					} else {
						pdf.CellFormat(50, 7, str, "1", 0, "", false, 0, "")
					}
				}
				pdf.Ln(-1)

				i++
				println("==", i)
			}
			j++

			pdf.SetFont("prompt", "", 10)
			pdf.SetX(left)

			pdf.CellFormat(20, 6, strconv.Itoa(j), "1", 0, "C", false, 0, "")
			pdf.CellFormat(50, 6, c.districtStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(50, 6, c.cityThStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(50, 6, c.cityEnStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(20, 6, c.popStr, "1", 0, "", false, 0, "")

			pdf.Ln(-1)

		}

	}
	// // Better table
	// improvedTable := func() {
	// 	// Column widths
	// 	w := []float64{40.0, 35.0, 40.0, 45.0}
	// 	wSum := 0.0
	// 	for _, v := range w {
	// 		wSum += v
	// 	}
	// 	left := (210 - wSum) / 2
	// 	// 	Header
	// 	pdf.SetX(left)
	// 	for j, str := range header {
	// 		pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
	// 	}
	// 	pdf.Ln(-1)
	// 	// Data
	// 	for _, c := range countryList {
	// 		pdf.SetX(left)
	// 		pdf.CellFormat(w[0], 6, c.districtStr, "LR", 0, "", false, 0, "")
	// 		pdf.CellFormat(w[1], 6, c.cityThStr, "LR", 0, "", false, 0, "")
	// 		pdf.CellFormat(w[2], 6, strDelimit(c.cityEnStr, ",", 3),
	// 			"LR", 0, "R", false, 0, "")
	// 		pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
	// 			"LR", 0, "R", false, 0, "")
	// 		pdf.Ln(-1)
	// 	}
	// 	pdf.SetX(left)
	// 	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	// }
	// // Colored table
	// fancyTable := func() {
	// 	// Colors, line width and bold font
	// 	pdf.SetFillColor(255, 0, 0)
	// 	pdf.SetTextColor(255, 255, 255)
	// 	pdf.SetDrawColor(128, 0, 0)
	// 	pdf.SetLineWidth(.3)
	// 	pdf.SetFont("", "B", 0)
	// 	// 	Header
	// 	w := []float64{40, 35, 40, 45}
	// 	wSum := 0.0
	// 	for _, v := range w {
	// 		wSum += v
	// 	}
	// 	left := (210 - wSum) / 2
	// 	pdf.SetX(left)
	// 	for j, str := range header {
	// 		pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
	// 	}
	// 	pdf.Ln(-1)
	// 	// Color and font restoration
	// 	pdf.SetFillColor(224, 235, 255)
	// 	pdf.SetTextColor(0, 0, 0)
	// 	pdf.SetFont("", "", 0)
	// 	// 	Data
	// 	fill := false
	// 	for _, c := range countryList {
	// 		pdf.SetX(left)
	// 		pdf.CellFormat(w[0], 6, c.districtStr, "LR", 0, "", fill, 0, "")
	// 		pdf.CellFormat(w[1], 6, c.cityThStr, "LR", 0, "", fill, 0, "")
	// 		pdf.CellFormat(w[2], 6, strDelimit(c.cityEnStr, ",", 3),
	// 			"LR", 0, "R", fill, 0, "")
	// 		pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
	// 			"LR", 0, "R", fill, 0, "")
	// 		pdf.Ln(-1)
	// 		fill = !fill
	// 	}
	// 	pdf.SetX(left)
	// 	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	// }
	loadData("countries.txt")
	// pdf.SetFont("Prompt", "", 14)

	// pdf.AddPage()
	basicTable()
	// pdf.AddPage()
	// improvedTable()
	// pdf.AddPage()
	// fancyTable()
	fileStr := "simple.pdf"
	err := pdf.OutputFileAndClose(fileStr)

	print(err)
}
