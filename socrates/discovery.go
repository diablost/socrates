package socrates

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// personal.unlockAccount(eth.accounts[2], "discovery")

//var registerContract = web3.eth.contract([{"constant":false,"inputs":[],"name":"kill","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getInfo","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_info","type":"string"}],"name":"setInfo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_info","type":"string"}],"name":"InfoChanged","type":"event"}]);var register = registerContract.new(
//{
//  from: web3.eth.accounts[0], 
//  data: '0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506104ad806100606000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806341c0e1b5146100675780635a9b0b891461007e5780638da5cb5b1461010e578063937f6e7714610165575b600080fd5b34801561007357600080fd5b5061007c6101ce565b005b34801561008a57600080fd5b5061009361025f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100d35780820151818401526020810190506100b8565b50505050905090810190601f1680156101005780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561011a57600080fd5b50610123610301565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561017157600080fd5b506101cc600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610326565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561025d576000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b565b606060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102f75780601f106102cc576101008083540402835291602001916102f7565b820191906000526020600020905b8154815290600101906020018083116102da57829003601f168201915b5050505050905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b806001908051906020019061033c9291906103dc565b507f53903e402ecb1380ca1f307cbedde216f9251ea107bb8ef363bfed600d5ed865816040518080602001828103825283818151815260200191508051906020019080838360005b8381101561039f578082015181840152602081019050610384565b50505050905090810190601f1680156103cc5780820380516001836020036101000a031916815260200191505b509250505060405180910390a150565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061041d57805160ff191683800117855561044b565b8280016001018555821561044b579182015b8281111561044a57825182559160200191906001019061042f565b5b509050610458919061045c565b5090565b61047e91905b8082111561047a576000816000905550600101610462565b5090565b905600a165627a7a723058204b7f02dea3e394db59081bd8e1d8a39e9cde32ea163a2380c127992537c203910029', 
//  gas: '4700000'
//}, function (e, contract){
// console.log(e, contract);
// if (typeof contract.address !== 'undefined') {
//	  console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
// }
//})

var ethaddr = "http://13.230.37.18:18545"
var contract = "0x863a5d8911988c6f42cd0ca0a344ea61ca998212"
var account = "0x77a5ffdca2a406bd4f8ac99e4ea695165df10ac0"
var ssTemplate = "ss://AEAD_CHACHA20_POLY1305:test1234@%s"
var pos0 = "290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"
//var pos1 = "0xb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6"

// Discovery  1
// input :c channel used for nodify
// return value: slice of discovered hosts
func Discovery(c chan int) []string {

	cli, err := ethclient.Dial(ethaddr)
	if err != nil {
		log.Fatal(err)
	}

	newHost := getContractData(cli)
	if err != nil {
		log.Fatal(err)
	}

	logf("get proxy server from eth contract:%v", newHost)
	return newHost
}

func getContractData(cli *ethclient.Client) []string {

	contractid := common.HexToAddress(contract)
	posShort := common.HexToHash("0")

	byteData, err := cli.PendingStorageAt(context.Background(), contractid, posShort)
	if err != nil {
		logf("get contract data error occurs:%v", err)
	}

	// remove "0x" and "\x0e" and "$"(0x24) and "0"
	//log.Printf("byteData:%v", strings.TrimRight(strings.TrimRight(strings.TrimRight(strings.TrimLeft(common.ToHex(byteData), "0x"), "0e"), "24"), "0"))
	decodeData, err := hex.DecodeString(strings.TrimRight(strings.TrimRight(strings.TrimRight(strings.TrimLeft(common.ToHex(byteData), "0x"), "0e"), "24"), "0"))
	if err != nil {
		logf("get contract data error occurs:%v", err)
	}
	if len(decodeData) < 5 {
		decodeData = []byte {}
		nextPos := pos0
		for {
			p, err := cli.PendingStorageAt(context.Background(), contractid, common.HexToHash(nextPos))
			if err != nil {
				break
			}
			strData := strings.TrimRight(strings.TrimRight(strings.TrimRight(strings.TrimLeft(common.ToHex(p), "0x"), "0e"), "24"), "0")
			q, _ := hex.DecodeString(strData)

			logf("get %v Contract storage pos:%v %v,%v,%v,%v", contractid, pos0, common.ToHex(p), strData, q, string(q))
			if strings.TrimRight(strData, "0") == "" {
				break
			} else {
				decodeData = append(decodeData, q...)
			}

			nextPos = hexAddition(nextPos)
		}
	}

	return strings.Split(string(decodeData), ",")
}

func hexAddition(s string) string {
	h := new(big.Int)
	h.SetString(s, 16)
	int1 := new(big.Int)
	int1.SetString("1", 16)
	h = h.Add(h, int1)
	return h.Text(16)
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
