# Transaction monitoring - Ethereum JSON-RPC
Golang system that monitors incoming and outgoing transactions of specified Ethereum addresses and store the transaction into database.

## How to run program
**Option #1: Run on your local**
1. edit .env file
example
```bash
MONITORED_ADDRESS=0x28c6c06298d514db089934071355e5743bf21d60
JSONRPC_SERVER=https://mainnet.infura.io/v3/c37e80f7c9e646e3a3db8a6ddec8dcd4
DB_HOST=localhost
DB_DRIVER=postgres
DB_USER=admin
DB_PASSWORD=password
DB_NAME=postgres
DB_PORT=5432
```
2. sync dependencies

```bash
  go mod tidy
```

3. run program

```bash
  go run main.go
```

**Option #2: Run on docker**
```bash
  docker-compose build
  docker-compose up
```
---
    
## Program Explanation
The monitoring transaction system is designed to continuously check the Ethereum blockchain for new transactions related to a specific address and store them in a database. The program works as follows:

1. Upon starting, the program connects to the Ethereum JSON-RPC API to retrieve the latest block.
2. The program checks whether the latest block contains any transactions related to the specified address.
3. If a transaction related to the address is found, the program will store transaction data in the database.
4. The program repeats steps 1-3 indefinitely, allowing it to monitor the blockchain for new transactions in real-time.

## Transaction retrieval API
#### API design approch
To retrieve transaction data from the database, I have implemented a ***pagination technique*** that allows users to request the data in smaller chunks. This approach is necessary due to the large volume of transaction data stored in the database, which could cause performance issues if retrieved all at once.

When making a request to the Transaction Retrieval API, users can specify the number of transactions they want to retrieve along with another query critiria. The API returns a JSON response containing the requested transaction data, as well as 'next_cursor_id' indicating the next cursor id which user should get next.

By using pagination, users can retrieve transaction data in a more manageable and efficient manner, while also reducing the strain on the database and the API. This approach ensures that users can access the transaction data they need without encountering performance issues or data retrieval failures.

#### API spec

<summary><code>GET</code> <code><b>/transaction</b></code></summary>

##### Query Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `address` |  optional | string   | The specific address        |
> | `start_block` |  optional | integer   | The starting block number to get (including this block no.)       |
> | `end_block` |  optional | integer   | The ending block number to get (including this block no.)        |
> | `cursor_id` |  optional | integer   | Next cursor id to get (use the previos 'next_cursor_id' from response to get next data)        |
> | `limit` |  optional | integer   | The limit of transaction from response        |


##### Responses

> | http code     | content-type                      |
> |---------------|-----------------------------------|
> | `200`         | `text/plain;charset=UTF-8`        |
example response data
```
{
    "transactions": [
        {
            "ID": 3,
            "TransactionHash": "0xb762e3e841fac73814b5b88487b4157057d583dd762a6a31f038e8f944ab63f6",
            "Address": "0x28c6c06298d514db089934071355e5743bf21d60",
            "BlockNumber": 17174653,
            "Data": {
                "blockHash": "0x71381eb9931a01de1ed0f5bf4774ea1338cbb3329e969add94b7b0bcc18b61a6",
                "blockNumber": "0x106107d",
                "chainId": "0x1",
                "from": "0x28c6c06298d514db089934071355e5743bf21d60",
                "gas": "0x32918",
                "gasPrice": "0x1853a28b17",
                "hash": "0xb762e3e841fac73814b5b88487b4157057d583dd762a6a31f038e8f944ab63f6",
                "input": "0x",
                "maxFeePerGas": "0x1d9db5f800",
                "maxPriorityFeePerGas": "0x77359400",
                "nonce": "0x64844b",
                "r": "0xbc8dc6304f70346962301ec2934ad45786e21b8b8220989c8597222356dc2cc4",
                "s": "0x6379e88ea0c6dae9396461096258712c663ef830ac6536e557209f64b42b964f",
                "to": "0x1bf57b7e351755a98e38d1ef31acaca9a7effd11",
                "transactionIndex": "0x48",
                "type": "0x2",
                "v": "0x0",
                "value": "0x4d4fbdd81e54000"
            }
        }
    ],
    "next_cursor_id": 3
}
```
##### Example cURL

> ```javascript
>  curl --request GET 'http://localhost:8080/transaction?address=0x28c6c06298d514db089934071355e5743bf21d60&start_block=17174652&end_block=17174660&cursor_id=2&limit=1'
> ```
