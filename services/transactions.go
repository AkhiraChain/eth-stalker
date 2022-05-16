package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	coingecko "github.com/superoo7/go-gecko/v3"
)

const (
	// WEI is "wei"
	WEI string = "wei"
	// KWEI is "ada", "kwei", "kilowei", "femtoether"
	KWEI string = "kwei"
	// MWEI is "babbage", "mwei", "megawei", "picoether"
	MWEI string = "mwei"
	// GWEI is "shannon", "gwei", "gigawei", "nanoether", "nano"
	GWEI string = "gwei"
	// MICRO is "szazbo", "micro", "microether"
	MICRO string = "micro"
	// MILLI is "finney", "milli", "milliether"
	MILLI string = "milli"
	// ETH is "ether", "eth"
	ETH string = "ether"
	// KILO is "einstein", "kilo", "kiloether", "kether", "grand"
	KILO string = "kilo"
	// MEGA is "mega", "megaether", "mether"
	MEGA string = "mega"
	// GIGA is giga", "gigaether", "gether"
	GIGA string = "giga"
	// TERA is "tera", "teraether", "tether"
	TERA string = "tera"

	// USD is US Dollar
	USD string = "usd"
	// AED is United Arab Emirates Dirham
	AED string = "aed"
	// ARS is Argentine Peso
	ARS string = "ars"
)

// TransactionsResponse is returned by GetLatestBlockTransactions
type TransactionsResponse struct {
	Header       *types.Header
	Transactions []*types.Transaction
}

func weiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

// ValidUnit returns a normalized metric unit or an error
func ValidUnit(unit string) (string, error) {
	switch strings.ToLower(unit) {
	case "wei":
		return WEI, nil
	case "ada", "kwei", "kilowei", "femtoether":
		return KWEI, nil
	case "babbage", "mwei", "megawei", "picoether":
		return MWEI, nil
	case "shannon", "gwei", "gigawei", "nanoether", "nano":
		return GWEI, nil
	case "szazbo", "micro", "microether":
		return MICRO, nil
	case "finney", "milli", "milliether":
		return MILLI, nil
	case "ether", "eth":
		return ETH, nil
	case "einstein", "kilo", "kiloether", "kether", "grand":
		return KILO, nil
	case "mega", "megaether", "mether":
		return MEGA, nil
	case "giga", "gigaether", "gether":
		return GIGA, nil
	case "tera", "teraether", "tether":
		return TERA, nil
	case "usd", "USD":
		return USD, nil
	}
	return "", fmt.Errorf("Unknown unit %s", unit)
}

// ToWeiMultiplier returns the multipler to convert a unit to wei
func ToWeiMultiplier(normalizeUnit string) decimal.Decimal {
	var multiplier decimal.Decimal
	switch normalizeUnit {
	case WEI:
		multiplier, _ = decimal.NewFromString("1")
	case KWEI:
		multiplier, _ = decimal.NewFromString("1000")
	case MWEI:
		multiplier, _ = decimal.NewFromString("1000000")
	case GWEI:
		multiplier, _ = decimal.NewFromString("1000000000")
	case MICRO:
		multiplier, _ = decimal.NewFromString("1000000000000")
	case MILLI:
		multiplier, _ = decimal.NewFromString("1000000000000000")
	case ETH:
		multiplier, _ = decimal.NewFromString("1000000000000000000")
	case KILO:
		multiplier, _ = decimal.NewFromString("1000000000000000000000")
	case MEGA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000")
	case GIGA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000000")
	case TERA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000000000")
	}
	return multiplier
}

// FromWeiMultiplier returns the multipler to convert a unit to wei
func FromWeiMultiplier(normalizeUnit string) decimal.Decimal {
	var multiplier decimal.Decimal
	switch normalizeUnit {
	case WEI:
		multiplier, _ = decimal.NewFromString("1")
	case KWEI:
		multiplier, _ = decimal.NewFromString("0.001")
	case MWEI:
		multiplier, _ = decimal.NewFromString("0.000001")
	case GWEI:
		multiplier, _ = decimal.NewFromString("0.000000001")
	case MICRO:
		multiplier, _ = decimal.NewFromString("0.000000000001")
	case MILLI:
		multiplier, _ = decimal.NewFromString("0.000000000000001")
	case ETH:
		multiplier, _ = decimal.NewFromString("0.000000000000000001")
	case KILO:
		multiplier, _ = decimal.NewFromString("0.000000000000000000001")
	case MEGA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000001")
	case GIGA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000000001")
	case TERA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000000000001")
	}
	return multiplier
}

// ConvertToWei will convert any valid ethereum unit to wei
func ConvertToWei(normalizeUnit string, amount decimal.Decimal) decimal.Decimal {
	var result decimal.Decimal
	result = amount.Mul(ToWeiMultiplier(normalizeUnit))
	return result
}

// ConvertFromWei will convert any valid ethereum unit to wei
func ConvertFromWei(normalizeUnit string, amount decimal.Decimal) decimal.Decimal {
	var result decimal.Decimal
	result = amount.Mul(FromWeiMultiplier(normalizeUnit))
	return result
}

// ConvertToUSD uses Coinmarketcap to estimate value of ETH in USD
func ConvertToUSD(amountInWei string) (decimal.Decimal, error) {
	zero, _ := decimal.NewFromString("0")
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	client := coingecko.NewClient(httpClient)
	singlePrice, err := client.SimpleSinglePrice("ethereum", "usd")
	if err != nil {
		log.Fatal(err)
	}
	balanceInWei, err := decimal.NewFromString(amountInWei)
	if err != nil {
		return zero, err
	}
	price := decimal.NewFromFloat32(singlePrice.MarketPrice)
	balanceInETH := ConvertFromWei(ETH, balanceInWei)
	exchangeValue := price.Mul(balanceInETH)
	// trim to 2 decimal places
	return exchangeValue.Truncate(2), nil
}

func getEthClient() *ethclient.Client {
	apiKey := "wss://eth-mainnet.alchemyapi.io/v2/%s" + os.Getenv("ALCHEMY_API_KEY")
	client, err := ethclient.Dial(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getHeader(ctx context.Context, client *ethclient.Client) *types.Header {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return header
}

func getLatestBlock(ctx context.Context, client *ethclient.Client) uint64 {
	latestBlockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return latestBlockNumber
}

func inspectBlock(ctx context.Context, client *ethclient.Client) *types.Block {
	header := getHeader(ctx, client)
	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(blockNumber, getLatestBlock(ctx, client))
	fmt.Println("Block Number:", block.Number().String())
	fmt.Println("Block Time:", block.Time())
	fmt.Println("Block Difficulty:", block.Difficulty().Uint64())
	fmt.Println("Block Hash:", block.Hash().Hex())
	fmt.Println("Block Transactions:", len(block.Transactions()), block.Transactions().Len())

	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Count:", count)
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
	}
	return block
}

func getMinerFromHeaderBlock(block *types.Block) map[string]interface{} {
	var headerInfo map[string]interface{}
	jsonByte, err := block.Header().MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(jsonByte), &headerInfo)
	return headerInfo
}

func getBalanceFromAccount(client *ethclient.Client, account common.Address) (balance *big.Int) {
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	return balance
}

func getTransactionByHash(block *types.Block, hash string) *types.Transaction {
	var targetTxn *types.Transaction
	txns := block.Transactions().Len()
	if targetTxn == nil {
		for i := 0; i < txns; i++ {
			fmt.Println(block.Transactions()[i].Hash().Hex())
			wanted := block.Transaction(block.Transactions()[i].Hash())
			targetTxn = wanted
		}
	}

	return targetTxn
}

// GetLatestBlockTransactions responds with the latest block with all transactions as JSON.
func GetLatestBlockTransactions(c *gin.Context) {
	ctx := context.Background()
	client := getEthClient()
	block := inspectBlock(ctx, client)
	minerInfo := getMinerFromHeaderBlock(block)
	account := common.HexToAddress(fmt.Sprintf("%v", minerInfo["miner"]))
	minerBalance := getBalanceFromAccount(client, account)
	fmt.Println("Miner account:", account, ", balance: ", weiToEther(minerBalance))
	m := TransactionsResponse{block.Header(), block.Body().Transactions}
	empJSON, err := json.MarshalIndent(block.Header(), "", "  ")
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(empJSON))
	c.IndentedJSON(http.StatusOK, m)
}

// StreamBlockTransactions streams block mined events as JSON
func StreamBlockTransactions(c *gin.Context) {
	ctx := context.Background()
	headers := make(chan *types.Header)
	chanStream := make(chan TransactionsResponse)
	client := getEthClient()
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer close(chanStream)
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case header := <-headers:

				block, err := client.BlockByHash(ctx, header.Hash())
				if err != nil {
					log.Fatal(err)
				}
				headerInfo := getMinerFromHeaderBlock(block)
				txnInfo := getTransactionByHash(block, header.TxHash.Hex())
				fmt.Println("Block Number", block.Number().Uint64())
				fmt.Println("Block Hash:", header.Hash().Hex())
				fmt.Println("Block Miner", headerInfo["miner"])
				fmt.Println("Block Transactions:", len(block.Transactions()))

				fmt.Println("Block Transaction: TX IS GOING TO  ", txnInfo.To().Hex())
				fmt.Println("Block Transaction: Txn Hash ", txnInfo.Hash().Hex())
				price, _ := ConvertToUSD(txnInfo.Value().String())
				fmt.Println("Block Transaction: VALUE $:", price)

				fmt.Println("------------------------------------------------------------------------------")
				m := TransactionsResponse{block.Header(), block.Body().Transactions}
				chanStream <- m
			}
		}
	}()
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			c.IndentedJSON(http.StatusOK, msg)
			return true
		}
		return false
	})
}
