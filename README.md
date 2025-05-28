# Smart Spending Checker

Smart Spending Checker is a simple system to help you control your monthly spending, allowing you to add installment purchases and track their impact on your budget.

## Features

*   **Add product:** Register a new product with installments, specifying the name, total value, and number of installments.
*   **Remove product:** Delete a product from your spending list.
*   **List months:** View products registered for each month.
*   **Update monthly profit:** Set or change your monthly profit to calculate the percentage used by your expenses.
*   **Edit product:** Modify information for an existing product, such as name, total value, and number of installments.
*   **Anticipate installments:** Calculate the total amount to pay if you want to anticipate a specific number of installments for a product.
*   **Monthly summary:** See a summary for the month, including your monthly profit, total installments, percentage used, and a strategy recommendation.

## How to Use

### Prerequisites

*   [Go](https://golang.org/dl/) installed and configured on your system.

### Installation

1.  Clone the repository:

    ```bash
    git clone [your-repository-url]
    cd [your-repository-name]
    ```

2.  Initialize the Go module:

    ```bash
    go mod init [your-module-name]
    ```

3.  Run the program:

    ```bash
    go run main.go
    ```

### Setup

1.  **Create the `data` folder:**

    ```bash
    mkdir data
    ```

2.  **Create the `products.json` file inside the `data` folder:**

    The `products.json` file is where the program stores your products and monthly profit data. If the file does not exist, the program will automatically create it with default values.

    Example content for `products.json`:

    ```json
    {
      "products": [],
      "monthly_profit": 0,
      "month": 5,
      "year": 2025
    }
    ```

    You can create the file manually or let the program create it on first run.

### Usage

1.  Run the program:

    ```bash
    go run main.go
    ```

2.  Follow the menu instructions to add, remove, list, or edit products, update your monthly profit, and anticipate installments.

## Important Note on Strategy

The program provides a recommendation based on the percentage of your monthly profit used for installments. **It's generally recommended to keep your spending below 70% of your monthly profit.** If the percentage exceeds this threshold, the program will advise against using your profit to pay for it and suggest creating a separate fund.

## Notes

*   The program saves data in the `data/products.json` file, so make sure the `data` folder and `products.json` file exist and have the correct permissions.
*   The program accepts both comma (`,`) and dot (`.`) as decimal separators when entering values.

## Contributing

Contributions are welcome! Feel free to open issues and submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
