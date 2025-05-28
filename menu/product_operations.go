package menu

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/product"
)

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
		time.Sleep(2 * time.Second)
		return
	}

	totalValue, err := readFloat(reader, "Valor total do produto (R$): ")
	if err != nil {
		fmt.Println("Valor invalido.")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Print("Em quantas vezes será parcelado: ")
	installmentsStr, _ := reader.ReadString('\n')
	installmentsStr = strings.TrimSpace(installmentsStr)
	installments, err := strconv.Atoi(installmentsStr)
	if err != nil || installments < 1 {
		fmt.Println("Número de parcelas inválido.")
		time.Sleep(2 * time.Second)
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

	time.Sleep(2 * time.Second)
}

func removeProduct(reader *bufio.Reader, list *product.ProductList) {
	title := " REMOVER PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para remover.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonth(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("\nTem certeza que deseja remover '%s'? (s/n): ", list.Products[idx].Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "s" && confirm != "sim" {
		fmt.Println("Operação cancelada.")
		time.Sleep(2 * time.Second)
		return
	}

	list.Products = append(list.Products[:idx], list.Products[idx+1:]...)

	fmt.Println(divider)
	fmt.Println("✅ Produto removido!")
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func editProduct(reader *bufio.Reader, list *product.ProductList) {
	title := " EDITAR PRODUTO "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para editar.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonth(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &list.Products[idx]

	fmt.Printf("Nome atual: %s. Novo nome (ou Enter para manter): ", p.Name)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)
	if newName != "" {
		p.Name = newName
	}

	fmt.Printf("Valor total atual: R$%.2f. Novo valor (ou Enter para manter): ", p.TotalValue)
	totalValueStr, _ := reader.ReadString('\n')
	totalValueStr = strings.TrimSpace(totalValueStr)
	if totalValueStr != "" {
		totalValueStr = strings.ReplaceAll(totalValueStr, ",", ".")
		totalValue, err := strconv.ParseFloat(totalValueStr, 64)
		if err == nil && totalValue > 0 {
			p.TotalValue = totalValue
		}
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

	time.Sleep(2 * time.Second)
}

func anticipateInstallments(reader *bufio.Reader, list *product.ProductList) {
	title := " ANTECIPAR PARCELAS "
	divider := strings.Repeat("-", 40)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	if len(list.Products) == 0 {
		fmt.Println("Nenhum produto para antecipar parcelas.")
		time.Sleep(2 * time.Second)
		return
	}

	idx, ok := selectProductByYearMonth(reader, list.Products)
	if !ok {
		time.Sleep(2 * time.Second)
		return
	}

	p := &list.Products[idx]

	now := time.Now()
	monthsElapsed := int(now.Sub(p.CreatedAt).Hours()/24/30) + 1
	if monthsElapsed > p.Installments {
		monthsElapsed = p.Installments
	}

	remainingInstallments := p.Installments - monthsElapsed + 1
	if remainingInstallments <= 0 {
		fmt.Println("Este produto já foi totalmente pago.")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Printf("Quantas parcelas deseja antecipar? (Total restante: %d): ", remainingInstallments)
	anticipateStr, _ := reader.ReadString('\n')
	anticipateStr = strings.TrimSpace(anticipateStr)
	anticipate, err := strconv.Atoi(anticipateStr)
	if err != nil || anticipate < 1 || anticipate > remainingInstallments {
		fmt.Println("Quantidade inválida.")
		time.Sleep(2 * time.Second)
		return
	}

	valorTotal := float64(anticipate) * p.Parcel

	fmt.Println(divider)
	fmt.Printf("Valor total para antecipar %d parcelas: R$%.2f\n", anticipate, valorTotal)
	fmt.Println(divider)

	fmt.Print("Deseja confirmar a antecipação? (s/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "s" && confirm != "sim" {
		fmt.Println("Operação cancelada.")
		time.Sleep(2 * time.Second)
		return
	}

	p.Installments -= anticipate

	fmt.Println("✅ Parcelas antecipadas com sucesso!")
	time.Sleep(2 * time.Second)
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
		time.Sleep(2 * time.Second)
		return
	}
	list.MonthlyProfit = profit
	list.Month = int(time.Now().Month())
	list.Year = time.Now().Year()

	fmt.Println(divider)
	fmt.Println("✅ Lucro mensal atualizado!")
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}

func configureSafePercentage(reader *bufio.Reader, list *product.ProductList) {
	title := " CONFIGURAR PORCENTAGEM SEGURA "
	divider := strings.Repeat("-", 50)

	fmt.Println("\n" + divider)
	fmt.Println(title)
	fmt.Println(divider)

	fmt.Println("A porcentagem segura define quanto do seu lucro mensal deve estar disponível para reinvestimento.")
	fmt.Println("Recomendação: Mantenha pelo menos 70% do seu lucro disponível para reinvestimento.")
	fmt.Printf("Porcentagem atual: %.0f%%\n", list.SafePercentage)

	fmt.Print("Nova porcentagem segura (ou Enter para manter): ")
	percentageStr, _ := reader.ReadString('\n')
	percentageStr = strings.TrimSpace(percentageStr)
	if percentageStr == "" {
		fmt.Println("Mantendo a porcentagem atual.")
		time.Sleep(2 * time.Second)
		return
	}

	percentageStr = strings.ReplaceAll(percentageStr, ",", ".")
	percentage, err := strconv.ParseFloat(percentageStr, 64)
	if err != nil || percentage <= 0 || percentage > 100 {
		fmt.Println("Valor inválido. Mantendo a porcentagem atual.")
		time.Sleep(2 * time.Second)
		return
	}

	list.SafePercentage = percentage

	fmt.Println(divider)
	fmt.Printf("✅ Porcentagem segura atualizada para %.0f%%!\n", percentage)
	fmt.Println(divider)

	time.Sleep(2 * time.Second)
}
