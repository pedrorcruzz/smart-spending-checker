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
		fmt.Printf("\n--- Gestor Inteligente de Gastos (%02d/%d) ---\n", now.Month(), now.Year())
		if list.MonthlyProfit == 0 {
			fmt.Println("Por favor, defina seu lucro mensal antes de adicionar produtos.")
			updateMonthlyProfit(reader, &list)
			continue
		}
		showSummary(list)

		fmt.Println("\n1. Adicionar produto")
		fmt.Println("2. Remover produto")
		fmt.Println("3. Listar meses")
		fmt.Println("4. Atualizar lucro mensal")
		fmt.Println("5. Editar produto")
		fmt.Println("6. Antecipar parcelas")
		fmt.Println("7. Sair")
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
	fmt.Printf("Produto adicionado! Parcela mensal: R$%.2f\n", parcel)
}

func removeProduct(reader *bufio.Reader, list *product.ProductList) {
	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para remover.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("Digite o numero do produto para remover: ")
	numStr, _ := reader.ReadString('\n')
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(list.Products) {
		fmt.Println("Numero invalido.")
		return
	}
	list.Products = append(list.Products[:num-1], list.Products[num:]...)
	fmt.Println("Produto removido!")
}

func listMonths(reader *bufio.Reader, list product.ProductList) {
	fmt.Println("\nListar produtos por mês:")
	fmt.Println("1. Mês atual")
	for i := 1; i <= 12; i++ {
		fmt.Printf("%d. %s\n", i+1, monthNames[i-1])
	}
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

	if month > 0 {
		monthName := monthNames[month-1]
		fmt.Printf("\nResumo do mes (%02d/%d - %s):\n", month, year, monthName)
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

	fmt.Println("\nProdutos cadastrados:")
	for i, p := range monthlyProducts {
		fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d vezes) | Adicionado em: %s\n",
			i+1, p.Name, p.TotalValue, p.Parcel, p.Installments, p.CreatedAt.Format("02/01/2006"))
	}
}

func updateMonthlyProfit(reader *bufio.Reader, list *product.ProductList) {
	profit, err := readFloat(reader, "Novo lucro mensal (R$): ")
	if err != nil {
		fmt.Println("Valor invalido.")
		return
	}
	list.MonthlyProfit = profit
	list.Month = int(time.Now().Month())
	list.Year = time.Now().Year()
	fmt.Println("Lucro mensal atualizado!")
}

func editProduct(reader *bufio.Reader, list *product.ProductList) {
	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para editar.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("Digite o numero do produto para editar: ")
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
	fmt.Println("Produto atualizado!")
}

func anticipateInstallments(reader *bufio.Reader, list *product.ProductList) {
	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para antecipar parcelas.")
		return
	}
	listProducts(reader, *list, 0, 0)
	fmt.Print("Digite o numero do produto para antecipar parcelas: ")
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
	fmt.Printf("Valor total para antecipar %d parcelas: R$%.2f\n", anticipate, valorTotal)
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

	if len(suggestedProducts) == 1 {
		fmt.Printf("Sugestão: Separe o produto '%s' (Parcela: R$%.2f) em uma caixinha separada.\n",
			suggestedProducts[0].Name, suggestedProducts[0].Parcel)
	} else if len(suggestedProducts) > 1 {
		fmt.Println("Sugestão: Separe os seguintes produtos em uma caixinha:")
		for i, p := range suggestedProducts {
			fmt.Printf("  %d. %s (Parcela: R$%.2f)\n", i+1, p.Name, p.Parcel)
		}
		fmt.Printf("Total a separar: R$%.2f\n", suggestedParcelSum)
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

	fmt.Printf("\nResumo do mes (%02d/%d - %s):\n", list.Month, list.Year, monthName)
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
		fmt.Println("\nProdutos ativos:")
		for i, p := range list.Products {
			fmt.Printf("%d. %s | Total: R$%.2f | Parcela: R$%.2f (%d vezes)\n",
				i+1, p.Name, p.TotalValue, p.Parcel, p.Installments)
		}
	}
}
