package services

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func Test_weiToEther(t *testing.T) {
	type args struct {
		wei *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Float
	}{
		// TODO: Add test cases.
		{
			name: "weiToEther",
			args: args{
				wei: big.NewInt(1000000000000000000),
			},
			want: big.NewFloat(1),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := weiToEther(tt.args.wei); got == tt.want {
				t.Errorf("weiToEther() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLatestBlock(t *testing.T) {
	ethc, _ := ethclient.Dial("wss://eth-mainnet.alchemyapi.io/v2/o_bo9q2LMtGvYqr7jsyYSpUrE_azdh9x")
	ef, _ := ethc.BlockByNumber(context.Background(), nil)
	type args struct {
		ctx    context.Context
		client *ethclient.Client
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
		{
			name: "getLatestBlock",
			args: args{
				ctx:    context.Background(),
				client: ethc,
			},
			want: ef.Number().Uint64(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getLatestBlock(tt.args.ctx, tt.args.client); got != tt.want {
				t.Errorf("getLatestBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func Test_getMinerFromHeaderBlock(t *testing.T) {
	var headerInfo map[string]interface{}
	ethc , _ := ethclient.Dial("wss://eth-mainnet.alchemyapi.io/v2/o_bo9q2LMtGvYqr7jsyYSpUrE_azdh9x")
	blocc, _ := ethc.BlockByNumber(context.Background(), nil)
	type args struct {
		block *types.Block
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
		{
			name: "getMinerFromHeaderBlock",
			args: args{
				block: blocc,
			},
			want: headerInfo,
	},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getMinerFromHeaderBlock(tt.args.block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMinerFromHeaderBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func Test_getBalanceFromAccount(t *testing.T) {
	ethc, _ := ethclient.Dial("wss://eth-mainnet.alchemyapi.io/v2/o_bo9q2LMtGvYqr7jsyYSpUrE_azdh9x")
	bal, _ := ethc.BalanceAt(context.Background(), common.HexToAddress("0x0d1d4e623d10f9fba5db95830f7d3839406c6af2"), nil)
	account := common.HexToAddress("0x0d1d4e623d10f9fba5db95830f7d3839406c6af2")

	type args struct {
		client  *ethclient.Client
		account common.Address
	}
	tests := []struct {
		name        string
		args        args
		wantBalance *big.Int
	}{
		{
			name: "getBalanceFromAccount",
			args: args{
				client:  ethc,
				account: account,
			},
			wantBalance: bal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotBalance := getBalanceFromAccount(tt.args.client, tt.args.account); !reflect.DeepEqual(gotBalance, tt.wantBalance) {
				t.Errorf("getBalanceFromAccount() = %v, want %v", gotBalance, tt.wantBalance)
			}
		})
	}
}

func TestGetLatestBlockTransactions(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetLatestBlockTransactions",
			args: args{
				c: &gin.Context{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			GetLatestBlockTransactions(tt.args.c)
		})
	}
}

func TestGetTransactions(t *testing.T) {
	ethc, _ := ethclient.Dial(
		"wss://eth-mainnet.alchemyapi.io/v2/o_bo9q2LMtGvYqr7jsyYSpUrE_azdh9x")
	blocc, _ := ethc.BlockByNumber(context.Background(), nil)
	hashes := blocc.Transactions().Len()
	var want *types.Transaction
	if want == nil {
		for i := 0; i < hashes; i++ {
			fmt.Println(blocc.Transactions()[i].Hash().Hex())
			wanted := blocc.Transaction(blocc.Transactions()[i].Hash())
			want = wanted
		}
	}
	t.Log(want)
	type args struct {
		block *types.Block
		hash  string
	}
	tests := []struct {
		name    string
		args    args
		wantTxn *types.Transaction
	}{
		{
			name: "GetLatestBlockTransactions",
			args: args{
				block: blocc,
				hash:  want.Hash().Hex(),
			},
			wantTxn: want,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			getTransactionByHash(tt.args.block, tt.args.hash)
		})
	}
	t.Log(want)
}
