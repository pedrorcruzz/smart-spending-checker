package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/storage"
	"github.com/pedrorcruzz/smart-spending-checker/utils"
)

var monthNames = []string{
	"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
	"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
}

func ShowMenu() {
	reader := bufio.NewReader(os.Stdin)
	list, _ := storage.LoadProducts()

	if list.SafePercentage == 0 {
		list.SafePercentage = 70
	}

	for {
		utils.ClearTerminal()
		title := " Gestor Inteligente de Gastos "
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
		fmt.Println("7. Configurar porcentagem segura")
		fmt.Println("8. Sair")
		fmt.Println(menuDivider)
		fmt.Print("Escolha uma opcão: ")
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
			utils.ClearTerminal()
			configureSafePercentage(reader, &list)
		case "8":
			storage.SaveProducts(list)
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opcão inválida.")
			time.Sleep(1 * time.Second)
		}
		storage.SaveProducts(list)
	}
}
