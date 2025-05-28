package menu

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/product"
	"github.com/pedrorcruzz/smart-spending-checker/storage"
	"github.com/pedrorcruzz/smart-spending-checker/utils"
)

var monthNames = []string{
	"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
	"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
}

func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	fmt.Print(prompt)
	valueStr, _ := reader.ReadString('\n')
	valueStr = strings.TrimSpace(valueStr)
	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	return strconv.ParseFloat(valueStr, 64)
}

func ShowMenu() {
	reader := bufio.NewReader(os.Stdin)
	list, _ := storage.LoadProducts()
	for {
		utils.ClearTerminal()
		now := time.Now()
		title := fmt.Sprintf(" Gestor Inteligente de Gastos (%02d/%d) ", now.Month(), now.Year())
		divider := strings.Repeat("=", len(title)+10)

		fmt.Println("\n" + divider)
		fmt.Println(strings.Repeat(" ", 5) + title)
		fmt.Println(divider)

		if list.MonthlyProfit == 0 {
			fmt.Println("\nPor favor, defina seu lucro mensal antes de adicionar produtos.")
			updateMonthlyProfit(reader, &list)
			continue
		}
		showSummary(list)

		menuDivider := strings.Repeat("-", 40)
		fmt.Println("\n" + menuDivider)
		fmt.Println(" MENU PRINCIPAL")
		fmt.Println(menuDivider)
		fmt.Println("1. Adicionar produto")
		fmt.Println("2. Remover produto")
		fmt.Println("3. Listar meses")
		fmt.Println("4. Atualizar lucro mensal")
		fmt.Println("5. Editar produto")
		fmt.Println("6. Antecipar parcelas")
		fmt.Println("7. Sair")
		fmt.Println(menuDivider)
		fmt.Print("Escolha uma opcao: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			utils.ClearTerminal()
			addProduct(reader, &list)
		case "2":
			utils.ClearTerminal()
			removeProduct(reader, &list)
		case "3":
			utils.ClearTerminal()
			listMonths(reader, list)
		case "4":
			utils.ClearTerminal()
			updateMonthlyProfit(reader, &list)
		case "5":
			utils.ClearTerminal()
			editProduct(reader, &list)
		case "6":
			utils.ClearTerminal()
			anticipateInstallments(reader, &list)
		case "7":
			storage.SaveProducts(list)
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opcao invalida.")
		}
		storage.SaveProducts(list)
	}
}

func addProduct(reader *bufio.Reader, list *product.ProductList) {
	title := " ADICIONAR PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	fmt.Print("Nome do produto: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if _, err := strconv.Atoi(name); err == nil {
		fmt.Println("Nome invalido.")
		return
	}

	totalValue, err := readFloat(reader, "Valor total do produto (R$): ")
	if err != nil {
		fmt.Println("Valor invalido.")
		return
	}

	fmt.Print("Em quantas vezes será parcelado: ")
	installmentsStr, _ := reader.ReadString('\n')
	installmentsStr = strings.TrimSpace(installmentsStr)
	installments, err := strconv.Atoi(installmentsStr)
	if err != nil || installments < 1 {
		fmt.Println("Número de parcelas inválido.")
		return
	}

	parcel := totalValue / float64(installments)

	list.Products = append(list.Products, product.Product{
		Name:         name,
		Parcel:       parcel,
		TotalValue:   totalValue,
		Installments: installments,
		CreatedAt:    time.Now(),
	})
	list.Month = int(time.Now().Month())
	list.Year = time.Now().Year()

	fmt.Println(divider)
	fmt.Printf("✅ Produto adicionado! Parcela mensal: R$%.2f\n", parcel)
	fmt.Println(divider)
}

func removeProduct(reader *bufio.Reader, list *product.ProductList) {
	title := " REMOVER PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para remover.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("\nDigite o numero do produto para remover: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(list.Products) {
		fmt.Println("Numero invalido.")
		return
	}
	list.Products = append(list.Products[:num-1], list.Products[num:]...)

	fmt.Println(divider)
	fmt.Println("✅ Produto removido!")
	fmt.Println(divider)
}

func listMonths(reader *bufio.Reader, list product.ProductList) {
	title := " LISTAR PRODUTOS POR MÊS "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	fmt.Println("1. Mês atual")
	for i := 1; i <= 12; i++ {
		fmt.Printf("%d. %s\n", i+1, monthNames[i-1])
	}
	fmt.Println(divider)
	fmt.Print("Escolha um mês: ")
	monthStr, _ := reader.ReadString('\n')
	monthStr = strings.TrimSpace(monthStr)
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 13 {
		fmt.Println("Mês inválido.")
		return
	}

	var selectedMonth, selectedYear int
	if month == 1 {
		now := time.Now()
		selectedMonth = int(now.Month())
		selectedYear = now.Year()
	} else {
		selectedMonth = month - 1
		selectedYear = time.Now().Year()
	}

	listProducts(reader, list, selectedMonth, selectedYear)
}

func listProducts(reader *bufio.Reader, list product.ProductList, month int, year int) {
	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return
	}

	var monthlyProducts []product.Product
	var totalParcel float64

	for i := range list.Products {
		if month == 0 || (list.Products[i].CreatedAt.Month() == time.Month(month) && list.Products[i].CreatedAt.Year() == year) {
			monthlyProducts = append(monthlyProducts, list.Products[i])
			totalParcel += list.Products[i].Parcel
		}
	}

	if len(monthlyProducts) == 0 {
		fmt.Println("Nenhum produto cadastrado para este mês.")
		return
	}

	divider := strings.Repeat("-", 60)

	if month > 0 {
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

			if leftPercent >= 70 {
				fmt.Println("✅ Você pode usar seu lucro.")
			} else {
				fmt.Println("❌ Não recomendado. Crie uma caixinha separada para alguns produtos.")
				suggestProductsToSeparate(monthlyProducts, list.MonthlyProfit)
			}
		}
	}

	productsTitle := " PRODUTOS CADASTRADOS "
	fmt.Println("\n" + divider)
	fmt.Println(productsTitle)
	fmt.Println(divider)

	for i, p := range monthlyProducts {
		fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d vezes) | Adicionado em: %s\n",
			i+1, p.Name, p.TotalValue, p.Parcel, p.Installments, p.CreatedAt.Format("02/01/2006"))
	}
	fmt.Println(divider)
}

func updateMonthlyProfit(reader *bufio.Reader, list *product.ProductList) {
	title := " ATUALIZAR LUCRO MENSAL "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	profit, err := readFloat(reader, "Novo lucro mensal (R$): ")
	if err != nil {
		fmt.Println("Valor invalido.")
		return
	}
	list.MonthlyProfit = profit
	list.Month = int(time.Now().Month())
	list.Year = time.Now().Year()

	fmt.Println(divider)
	fmt.Println("✅ Lucro mensal atualizado!")
	fmt.Println(divider)
}

func editProduct(reader *bufio.Reader, list *product.ProductList) {
	title := " EDITAR PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para editar.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("\nDigite o numero do produto para editar: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(list.Products) {
		fmt.Println("Numero invalido.")
		return
	}
	p := &list.Products[num-1]

	fmt.Printf("Nome atual: %s. Novo nome (ou Enter para manter): ", p.Name)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)
	if newName != "" {
		p.Name = newName
	}

	totalValue, err := readFloat(reader, fmt.Sprintf("Valor total atual: R$%.2f. Novo valor (ou Enter para manter): ", p.TotalValue))
	if err == nil && totalValue > 0 {
		p.TotalValue = totalValue
	}

	fmt.Printf("Parcelas atuais: %d. Novo número de parcelas (ou Enter para manter): ", p.Installments)
	installmentsStr, _ := reader.ReadString('\n')
	installmentsStr = strings.TrimSpace(installmentsStr)
	if installmentsStr != "" {
		installments, err := strconv.Atoi(installmentsStr)
		if err == nil && installments > 0 {
			p.Installments = installments
		}
	}

	p.Parcel = p.TotalValue / float64(p.Installments)

	fmt.Println(divider)
	fmt.Println("✅ Produto atualizado!")
	fmt.Println(divider)
}

func anticipateInstallments(reader *bufio.Reader, list *product.ProductList) {
	title := " ANTECIPAR PARCELAS "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para antecipar parcelas.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("\nDigite o numero do produto para antecipar parcelas: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(list.Products) {
		fmt.Println("Numero invalido.")
		return
	}
	p := &list.Products[num-1]

	fmt.Printf("Quantas parcelas deseja antecipar? (Total restante: %d): ", p.Installments)
	anticipateStr, _ := reader.ReadString('\n')
	anticipateStr = strings.TrimSpace(anticipateStr)
	anticipate, err := strconv.Atoi(anticipateStr)
	if err != nil || anticipate < 1 || anticipate > p.Installments {
		fmt.Println("Quantidade inválida.")
		return
	}

	valorTotal := float64(anticipate) * p.Parcel

	fmt.Println(divider)
	fmt.Printf("Valor total para antecipar %d parcelas: R$%.2f\n", anticipate, valorTotal)
	fmt.Println(divider)
}

func suggestProductsToSeparate(products []product.Product, monthlyProfit float64) {
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

	targetParcel := totalParcel - (monthlyProfit * 0.7)
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

func showSummary(list product.ProductList) {
	var totalParcel float64
	for _, p := range list.Products {
		totalParcel += p.Parcel
	}
	usedPercent := (totalParcel / list.MonthlyProfit) * 100
	leftPercent := 100 - usedPercent

	monthName := monthNames[list.Month-1]

	summaryDivider := strings.Repeat("-", 60)
	title := fmt.Sprintf(" RESUMO DO MÊS (%02d/%d - %s) ", list.Month, list.Year, monthName)

	fmt.Println("\n" + summaryDivider)
	fmt.Println(title)
	fmt.Println(summaryDivider)

	fmt.Printf("Lucro mensal: R$%.2f\n", list.MonthlyProfit)
	fmt.Printf("Total de parcelas: R$%.2f\n", totalParcel)
	fmt.Printf("Usado: %.2f%% | Para reinvestir: %.2f%%\n", usedPercent, leftPercent)
	if leftPercent >= 70 {
		fmt.Println("✅ Você pode usar seu lucro.")
	} else {
		fmt.Println("❌ Não recomendado. Crie uma caixinha separada para alguns produtos.")
		suggestProductsToSeparate(list.Products, list.MonthlyProfit)
	}

	if len(list.Products) > 0 {
		productsTitle := " PRODUTOS ATIVOS "
		fmt.Println("\n" + summaryDivider)
		fmt.Println(productsTitle)
		fmt.Println(summaryDivider)

		for i, p := range list.Products {
			fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d vezes)\n",
				i+1, p.Name, p.TotalValue, p.Parcel, p.Installments)
		}
		fmt.Println(summaryDivider)
	}
}
