# Interact ethereum smart contract with go

## Install geth & tools

Download geth & tools from https://geth.ethereum.org/downloads/

e.g. `Geth & Tools 1.9.24`

Extract binaries from tar.gz file to `/usr/local/bin`

```bash
tar vzfx geth-alltools-linux-386-1.9.24-cc05b050.tar.gz
cd geth-alltools-linux-386-1.9.24-cc05b050
sudo mv * /usr/local/bin
```

Verify install

```bash
geth -version
abigen -version
```

## Install solc

Install solc by solc-select, please follow my article [here](https://liushiming.cn/article/solidity-version-manager-solc-select-introduction.html)

It helps to download all the solc versions, so you can switch as you like by `solc use version_you_like`.

```bash
$ solc --versions
0.7.4
0.7.3
0.7.2
0.7.1
0.7.0
...
```

## Generate go file from contract

Take HB token contract as an example, the contract code is at `https://etherscan.io/address/0x46705dfff24256421a05d056c29e81bdc09723b8#code

Create a working dir

```bash
mkdir ~/usdt
cd ~/usdt
```

Copy the contract to `TetherToken.sol`
 
```bash
vi TetherToken.sol
```

The contract version is `0.4.19` (according to the first line of contract `pragma solidity ^0.4.19;`) , so we need use 0.4.19 version of `solc` to compile it.
 
```bash
solc use 0.4.19
```

Use `solc` to generate `abi` & `bin` artifacts.

```bash
solc --abi HBtoken.sol --out abi
solc --bin HBtoken.sol --out bin
```

Use `abigen` to generate go file.

```bash
abigen --bin=bin/HBtoken.bin --abi=abi/HBtoken.abi --pkg=hbtoken --out=hbtoken.go
```

Flag explain: 
- pkg: go package name
- out: output filename

## ethClient

Connect to ethereum `client/client.go`

```go
var(
	EthClient *ethclient.Client
)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	var err error
	EthClient, err = ethclient.DialContext(ctx, "https://mainnet.infura.io/v3/e9d43fcc8b60466c9b8c6c5b8215475c")
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
		fmt.Println("tx cost: ", tx.Cost())
	}
}
```

## Interact with contract

Interact with contract by the contract go file we generated.

```go
func ContractInvoke(){
	address := common.HexToAddress("0x6f259637dcD74C767781E37Bc6133cd6A68aa161")
	instance, err := hbtoken.NewHbtoken(address, EthClient)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal("get token name failed: ", err.Error())
	}
	fmt.Println("token name: ", name)

	balance, err := instance.BalanceOf(&bind.CallOpts{Pending: false}, common.HexToAddress("0x46705dfff24256421a05d056c29e81bdc09723b8"))
	if err != nil {
		log.Fatal("contract interaction failed: ", err.Error())
	}
	fmt.Println("HB token balance: ",balance.Div(balance,big.NewInt(1e18)))
}
```

## Run it

```bash
go run main.go
```

output sample:

```bash
latest blockNum:  11249327
tx cost:  1140000087540000
token name:  HuobiToken
HB token balance:  1487479

```
