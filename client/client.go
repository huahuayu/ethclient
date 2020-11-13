package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/huahuayu/ethclient/contract/hbtoken"
	"log"
	"time"
)

var(
	EthClient *ethclient.Client
)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	EthClient, err := ethclient.DialContext(ctx, "https://mainnet.infura.io/v3/e9d43fcc8b60466c9b8c6c5b8215475c")
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	number, err := EthClient.BlockNumber(ctx)
	if err != nil {
		log.Fatal("ethereum connect failed")
	}
	fmt.Println("latest blockNum: ", number)

	tx, pending, _ := EthClient.TransactionByHash(ctx, common.HexToHash("0xbff5fa4aa3b503b9ae2b2b89332bb1cec736bac96c9eba30fb7f54522496a570"))
	if !pending {
		fmt.Println("tx cost: ",tx.Cost())
	}
}

func ContractInvoke(){
	address := common.HexToAddress("0x6f259637dcD74C767781E37Bc6133cd6A68aa161")
	instance, err := hbtoken.NewHbtoken(address,EthClient)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil{
		log.Fatal("get token name failed: " ,err.Error())
	}
	fmt.Println("token name: ", name)

	balance, err := instance.BalanceOf(&bind.CallOpts{Pending: false,}, common.HexToAddress("0x46705dfff24256421a05d056c29e81bdc09723b8"))
	if err != nil{
		log.Fatal("contract interaction failed: ",err.Error())
	}
	fmt.Println(balance)
}
