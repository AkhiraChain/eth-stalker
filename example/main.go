package main

import (
	"fmt"
	//"log"
	"os"

	"github.com/AkhiraChain/eth-stalker/types"
)

var clientID string

func init() {
	clientID = os.Getenv("API_KEY")
}

func main() {
	c := types.New()
	c.APIKey = clientID

	fmt.Print(c.APIKey)
	resp, err := c.GetAddressEthAdv("ethereum", "0x3282791d6fd713f1e94f4bfd565eaa78b3a0599d", map[string]string{"limit": "3", "offset": "0"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(resp)
	for i := range resp.Data {

	fmt.Printf("Type: %v\n", resp.Data[i].Address.Type)
	fmt.Printf("Spent in USD: %v\n", resp.Data[i].Address.SpentUsd)
	fmt.Printf("Number of transactions: %v\n", resp.Data[i].Address.TransactionCount)
	}
	//fmt.Print(c)
	//fmt.Print(clientID)
	//fmt.Print(c)

	/*
	resp, err := c.GetAddressEthAdv("ethereum", "0x3282791d6fd713f1e94f4bfd565eaa78b3a0599d", map[string]string{"limit": "3", "offset": "0"})
	if err != nil {
		log.Fatalln(err)
	}

	for i := range resp.Data {
		fmt.Printf("Type: %v\n", resp.Data[i].Address.Type)
		fmt.Printf("Spent in USD: %v\n", resp.Data[i].Address.SpentUsd)
		fmt.Printf("Number of transactions: %v\n", resp.Data[i].Address.TransactionCount)
		for j := range resp.Data[i].Calls {
			fmt.Printf("\nTransaction number %v:\n", j+1)
			fmt.Printf("ID: %v\n", resp.Data[i].Calls[j].BlockID)
			fmt.Printf("Value in USD: %v\n", resp.Data[i].Calls[j].ValueUsd)
		}

	}
	*/
}
