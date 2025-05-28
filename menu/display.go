package menu

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/product"
)

func showSummary(list product.ProductList) {
	var totalParcel float64

	now := time.Now()
	targetYear := now.Year()
	targetMonth := int(now.Month())

	var activeProducts []product.Product

	for _, p := range list.Products {
		startYear, startMonth := p.CreatedAt.Year(), int(p.CreatedAt.Month())
		endDate := p.CreatedAt.AddDate(0, p.Installments-1, 0)
		endYear, endMonth := endDate.Year(), int(endDate.Month())

		if (targetYear > startYear || (targetYear == startYear && targetMonth >= startMonth)) &&
			(targetYear < endYear || (targetYear == endYear && targetMonth <= endMonth)) {
			activeProducts = append(activeProducts, p)
			totalParcel += p.Parcel
		}
	}

	usedPercent := 0.0
	leftPercent := 100.0

	if list.MonthlyProfit > 0 {
		usedPercent = (totalParcel / list.MonthlyProfit) * 100
		leftPercent = 100 - usedPercent
	}

	monthName := monthNames[targetMonth-1]

	summaryDivider := strings.Repeat("-", 60)
	title := fmt.Sprintf(" RESUMO DO MÊS (%02d/%d - %s) ", targetMonth, targetYear, monthName)

	fmt.Println("\n" + summaryDivider)
	fmt.Println(title)
	fmt.Println(summaryDivider)

	fmt.Printf("Lucro mensal: R$%.2f\n", list.MonthlyProfit)
	fmt.Printf("Total de parcelas: R$%.2f\n", totalParcel)
	fmt.Printf("Usado: %.2f%% | Para reinvestir: %.2f%%\n", usedPercent, leftPercent)
	fmt.Printf("Porcentagem segura configurada: %.0f%%\n", list.SafePercentage)

	if leftPercent >= list.SafePercentage {
		fmt.Println("✅ Você pode usar seu lucro.")
	} else {
		fmt.Println("❌ Não recomendado. Crie uma caixinha separada para alguns produtos.")
		suggestProductsToSeparate(activeProducts, list.MonthlyProfit, list.SafePercentage)
	}

	if len(activeProducts) > 0 {
		productsTitle := " PRODUTOS ATIVOS NESTE MÊS "
		fmt.Println("\n" + summaryDivider)
		fmt.Println(productsTitle)
		fmt.Println(summaryDivider)

		for i, p := range activeProducts {
			installmentNumber := getInstallmentNumber(p, targetYear, targetMonth)
			fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d/%d)\n",
				i+1, p.Name, p.TotalValue, p.Parcel, installmentNumber, p.Installments)
		}
		fmt.Println(summaryDivider)
	}
}

func suggestProductsToSeparate(products []product.Product, monthlyProfit float64, safePercentage float64) {
	if len(products) == 0 {
		return
	}

	type ProductWithIndex struct {
		Index   int
		Product product.Product
	}

	productsWithIndex := make([]ProductWithIndex, len(products))
	for i, p := range products {
		productsWithIndex[i] = ProductWithIndex{i, p}
	}

	sort.Slice(productsWithIndex, func(i, j int) bool {
		return productsWithIndex[i].Product.Parcel > productsWithIndex[j].Product.Parcel
	})

	var totalParcel float64
	for _, p := range products {
		totalParcel += p.Parcel
	}

	targetParcel := totalParcel - (monthlyProfit * (safePercentage / 100))
	if targetParcel <= 0 {
		return
	}

	var suggestedProducts []product.Product
	var suggestedParcelSum float64

	for _, pwi := range productsWithIndex {
		if suggestedParcelSum >= targetParcel {
			break
		}
		suggestedProducts = append(suggestedProducts, pwi.Product)
		suggestedParcelSum += pwi.Product.Parcel
	}

	suggestionDivider := strings.Repeat("-", 50)

	if len(suggestedProducts) == 1 {
		fmt.Println(suggestionDivider)
		fmt.Printf("💡 Sugestão: Separe o produto '%s' (Parcela: R$%.2f) em uma caixinha separada.\n",
			suggestedProducts[0].Name, suggestedProducts[0].Parcel)
		fmt.Println(suggestionDivider)
	} else if len(suggestedProducts) > 1 {
		fmt.Println(suggestionDivider)
		fmt.Println("💡 Sugestão: Separe os seguintes produtos em uma caixinha:")
		for i, p := range suggestedProducts {
			fmt.Printf("  %d. %s (Parcela: R$%.2f)\n", i+1, p.Name, p.Parcel)
		}
		fmt.Printf("  Total a separar: R$%.2f\n", suggestedParcelSum)
		fmt.Println(suggestionDivider)
	}
}
