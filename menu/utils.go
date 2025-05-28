package menu

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/product"
)

func mapProductsByYearMonth(products []product.Product) map[int]map[int][]int {
	result := make(map[int]map[int][]int)
	for idx, p := range products {
		endDate := p.CreatedAt.AddDate(0, p.Installments-1, 0)
		currentDate := p.CreatedAt

		for !currentDate.After(endDate) {
			currentYear := currentDate.Year()
			currentMonth := int(currentDate.Month())

			if _, ok := result[currentYear]; !ok {
				result[currentYear] = make(map[int][]int)
			}
			result[currentYear][currentMonth] = append(result[currentYear][currentMonth], idx)

			currentDate = currentDate.AddDate(0, 1, 0)
		}
	}
	return result
}

func selectProductByYearMonth(reader *bufio.Reader, products []product.Product) (int, bool) {
	byYearMonth := mapProductsByYearMonth(products)
	if len(byYearMonth) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return -1, false
	}

	years := make([]int, 0, len(byYearMonth))
	for y := range byYearMonth {
		years = append(years, y)
	}
	sort.Ints(years)

	fmt.Println("\nSelecione o ano:")
	for i, y := range years {
		fmt.Printf("%d. %d\n", i+1, y)
	}
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)
	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(years) {
		fmt.Println("Ano inválido.")
		return -1, false
	}
	year := years[yearIdx-1]

	monthsMap := byYearMonth[year]
	months := make([]int, 0, len(monthsMap))
	for m := range monthsMap {
		months = append(months, m)
	}
	sort.Ints(months)

	fmt.Println("\nSelecione o mês:")
	for i, m := range months {
		fmt.Printf("%d. %s\n", i+1, monthNames[m-1])
	}
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)
	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(months) {
		fmt.Println("Mês inválido.")
		return -1, false
	}
	month := months[monthIdx-1]

	prodIndexes := monthsMap[month]

	uniqueIndexes := make([]int, 0)
	seen := make(map[int]bool)
	for _, idx := range prodIndexes {
		if !seen[idx] {
			seen[idx] = true
			uniqueIndexes = append(uniqueIndexes, idx)
		}
	}

	fmt.Println("\nSelecione o produto:")
	for i, idx := range uniqueIndexes {
		p := products[idx]
		fmt.Printf("%d. %s | Total: R$%.2f | Parcelas: %d | Adicionado em: %s\n",
			i+1, p.Name, p.TotalValue, p.Installments, p.CreatedAt.Format("02/01/2006"))
	}
	fmt.Print("Produto: ")
	prodStr, _ := reader.ReadString('\n')
	prodStr = strings.TrimSpace(prodStr)
	prodIdx, err := strconv.Atoi(prodStr)
	if err != nil || prodIdx < 1 || prodIdx > len(uniqueIndexes) {
		fmt.Println("Produto inválido.")
		return -1, false
	}
	return uniqueIndexes[prodIdx-1], true
}

func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Print(prompt)
	valueStr, _ := reader.ReadString('\n')
	valueStr = strings.TrimSpace(valueStr)
	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	return strconv.ParseFloat(valueStr, 64)
}

func isProductActiveInMonth(p product.Product, targetYear, targetMonth int) bool {
	startDate := p.CreatedAt
	endDate := startDate.AddDate(0, p.Installments-1, 0)

	targetDate := time.Date(targetYear, time.Month(targetMonth), 1, 0, 0, 0, 0, time.UTC)
	return !targetDate.Before(startDate) && !targetDate.After(endDate)
}

func getInstallmentNumber(p product.Product, targetYear, targetMonth int) int {
	startDate := p.CreatedAt
	yearDiff := targetYear - startDate.Year()
	monthDiff := targetMonth - int(startDate.Month())
	totalMonthDiff := yearDiff*12 + monthDiff + 1

	if totalMonthDiff < 1 {
		return 1
	}
	if totalMonthDiff > p.Installments {
		return p.Installments
	}
	return totalMonthDiff
}
