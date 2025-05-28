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

func listMonths(reader *bufio.Reader, list product.ProductList) {
	title := " LISTAR MESES "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)
	fmt.Println("0. Voltar ao Menu")
	fmt.Println(divider)

	byYearMonth := mapProductsByYearMonth(list.Products)
	if len(byYearMonth) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		time.Sleep(2 * time.Second)
		return
	}

	years := make([]int, 0, len(byYearMonth))
	for y := range byYearMonth {
		years = append(years, y)
	}
	sort.Ints(years)

	fmt.Println("\nSelecione o ano (0 para voltar):")
	for i, y := range years {
		fmt.Printf("%d. %d\n", i+1, y)
	}
	fmt.Print("Ano: ")
	yearStr, _ := reader.ReadString('\n')
	yearStr = strings.TrimSpace(yearStr)

	if yearStr == "0" {
		return
	}

	yearIdx, err := strconv.Atoi(yearStr)
	if err != nil || yearIdx < 1 || yearIdx > len(years) {
		fmt.Println("Ano inválido.")
		time.Sleep(2 * time.Second)
		return
	}
	year := years[yearIdx-1]

	monthsMap := byYearMonth[year]
	months := make([]int, 0, len(monthsMap))
	for m := range monthsMap {
		months = append(months, m)
	}
	sort.Ints(months)

	fmt.Println("\nSelecione o mês (0 para voltar):")
	for i, m := range months {
		fmt.Printf("%d. %s\n", i+1, monthNames[m-1])
	}
	fmt.Print("Mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)

	if monthStr == "0" {
		return
	}

	monthIdx, err := strconv.Atoi(monthStr)
	if err != nil || monthIdx < 1 || monthIdx > len(months) {
		fmt.Println("Mês inválido.")
		time.Sleep(2 * time.Second)
		return
	}
	month := months[monthIdx-1]

	prodIndexes := monthsMap[month]
	if len(prodIndexes) == 0 {
		fmt.Printf("\nNenhum produto encontrado para %s/%d.\n", monthNames[month-1], year)
		time.Sleep(2 * time.Second)
		return
	}

	uniqueIndexes := make([]int, 0)
	seen := make(map[int]bool)
	for _, idx := range prodIndexes {
		if !seen[idx] {
			seen[idx] = true
			uniqueIndexes = append(uniqueIndexes, idx)
		}
	}

	productsTitle := fmt.Sprintf(" PRODUTOS DE %s/%d ", monthNames[month-1], year)
	fmt.Println("\n" + divider)
	fmt.Println(productsTitle)
	fmt.Println(divider)

	for i, idx := range uniqueIndexes {
		p := list.Products[idx]
		installmentNumber := getInstallmentNumber(p, year, month)
		fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d)\n",
			i+1, p.Name, p.TotalValue, p.Parcel, installmentNumber, p.Installments)
	}
	fmt.Println(divider)

	fmt.Print("\nPressione Enter para voltar...")
	reader.ReadString('\n')
}

func showProductsByMonth(list product.ProductList, month int, year int) {
	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return
	}

	var monthlyProducts []product.Product
	var totalParcel float64

	targetDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for _, p := range list.Products {
		startDate := p.CreatedAt
		endDate := startDate.AddDate(0, p.Installments-1, 0)

		if !targetDate.Before(startDate) && !targetDate.After(endDate) {
			monthlyProducts = append(monthlyProducts, p)
			totalParcel += p.Parcel
		}
	}

	if len(monthlyProducts) == 0 {
		monthName := monthNames[month-1]
		fmt.Printf("Nenhum produto ativo para %s de %d.\n", monthName, year)
		return
	}

	divider := strings.Repeat("-", 60)
	monthName := monthNames[month-1]
	title := fmt.Sprintf(" RESUMO DO MÊS (%02d/%d - %s) ", month, year, monthName)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	fmt.Printf("Total de parcelas: R$%.2f\n", totalParcel)

	if list.MonthlyProfit > 0 {
		usedPercent := (totalParcel / list.MonthlyProfit) * 100
		leftPercent := 100 - usedPercent
		fmt.Printf("Usado: %.2f%% | Para reinvestir: %.2f%%\n", usedPercent, leftPercent)

		if leftPercent >= list.SafePercentage {
			fmt.Println("✅ Você pode usar parte do seu lucro para pagar as parcelas!")
		} else {
			fmt.Println("❌ Não recomendado. Crie uma caixinha separada para alguns produtos!")
			suggestProductsToSeparate(monthlyProducts, list.MonthlyProfit, list.SafePercentage)
		}
	}

	productsTitle := " PRODUTOS ATIVOS NESTE MÊS "
	fmt.Println("\n" + divider)
	fmt.Println(productsTitle)
	fmt.Println(divider)

	for i, p := range monthlyProducts {
		installmentNumber := getInstallmentNumber(p, year, month)

		fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d) | Adicionado em: %s\n",
			i+1, p.Name, p.TotalValue, p.Parcel, installmentNumber, p.Installments, p.CreatedAt.Format("02/01/2006"))
	}
	fmt.Println(divider)
}
