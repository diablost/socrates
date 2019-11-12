package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// personal.unlockAccount(eth.accounts[2], "discovery")

//var _greeting = "13.230.37.18:18488" ;var helloContract = web3.eth.contract([{"constant":true,"inputs":[],"name":"say","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_greeting","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]);var hello = helloContract.new(
//   _greeting,
//   {
//     from: web3.eth.accounts[2],
//     data: '0x60606040526040805190810160405280601481526020017f31313131313131313131313131313131313131310000000000000000000000008152506001908051906020019061004f929190610094565b50341561005b57600080fd5b60405161030438038061030483398101604052808051820191905050806000908051906020019061008d929190610094565b5050610139565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100d557805160ff1916838001178555610103565b82800160010185558215610103579182015b828111156101025782518255916020019190600101906100e7565b5b5090506101109190610114565b5090565b61013691905b8082111561013257600081600090555060010161011a565b5090565b90565b6101bc806101486000396000f300606060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063954ab4b214610046575b600080fd5b341561005157600080fd5b6100596100d4565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009957808201518184015260208101905061007e565b50505050905090810190601f1680156100c65780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100dc61017c565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156101725780601f1061014757610100808354040283529160200191610172565b820191906000526020600020905b81548152906001019060200180831161015557829003601f168201915b5050505050905090565b6020604051908101604052806000815250905600a165627a7a7230582027c66a908efda524a19d64dbfece11adcb46615a7957d92bfa6af8456805ef940029',
//     gas: '4700000'
//   }, function (e, contract){
//    console.log(e, contract);
//    if (typeof contract.address !== 'undefined') {
//         console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
//    }
// })

var ethaddr = "http://13.230.37.18:18545"
var contract = "0x0ef8f01b7b4445e472982641c71ab3d0ada638f0"
var account = "0x77a5ffdca2a406bd4f8ac99e4ea695165df10ac0"
var ssTemplate = "ss://AEAD_CHACHA20_POLY1305:test1234@%s"

func Discovery(c chan int) []string {

	var hosts []string
	cli, err := ethclient.Dial(ethaddr)
	if err != nil {
		log.Fatal(err)
	}

	byteData := getContractData(cli)
	// remove "0x" and "\x0e" and "$"(0x24) and "0"
	//log.Printf("byteData:%v", strings.TrimRight(strings.TrimRight(strings.TrimRight(strings.TrimLeft(common.ToHex(byteData), "0x"), "0e"), "24"), "0"))
	newHost, err := hex.DecodeString(strings.TrimRight(strings.TrimRight(strings.TrimRight(strings.TrimLeft(common.ToHex(byteData), "0x"), "0e"), "24"), "0"))
	if err != nil {
		log.Fatal(err)
	}

	newAddr := fmt.Sprintf(ssTemplate, string(newHost))
	hosts = append(hosts, string(newAddr))
	log.Println("get proxy server from eth contract:", hosts)
	return hosts
}

func getContractData(cli *ethclient.Client) []byte {

	contractid := common.HexToAddress(contract)
	//pos := common.HexToHash("0")
	pos := common.HexToHash("0")

	ret, err := cli.PendingStorageAt(context.Background(), contractid, pos)
	if err != nil {
		log.Println("get contract data error occurs:", err)
	}

	return ret
}

func accountBalance(cli *ethclient.Client) {
	account := common.HexToAddress(account)
	balance, err := cli.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Println("get account balance error occurs:", err)
	}

	log.Println(balance)
}

func getHeader(cli *ethclient.Client) (string, error) {
	header, err := cli.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("get block error occurs:", err)
		return "", err
	}
	log.Println(header.Number.String())
	return header.Number.String(), err
}

func updateContract() {

}
