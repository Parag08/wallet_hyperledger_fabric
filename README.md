# WALLET

An application which can create, manage and destroy wallet's on bloackchain

## Specification

HyperLedger fabric 1.2.0 is used for creating the wallet. Chaincode is written in golang. All wallet related information will be stored on hyperledger blockchain and hyperLedger fabric will be used to interact  application.

## Getting Started

### Prerequisites

Make sure docker and docker-compose is installed on the system. You can check if they are installed using following commands

```bash
docker -v
docker-compose -v
```

### Terminal 1 - Build network & Copy project files
```bash
curl -sSL http://bit.ly/2ysbOFE | bash -s 1.2.0
cp -rf wallet  fabric-samples/chaincode/
cd fabric-samples/chaincode-docker-devmode/  
docker-compose -f docker-compose-simple.yaml up
```

[above commands will start the fabric network bare minimum requirements you can refer
[chaincode developer tutorial](https://hyperledger-fabric.readthedocs.io/en/release-1.2/chaincode4ade.html)]

### Terminal 2 - Build & start the chaincode
```bash
docker exec -it chaincode bash
cd wallet
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=wallet:0 ./wallet
```

### Terminal 3 - Use the chaincode
```bash
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/wallet -n wallet -v 0

# I don't understand significance of -C myc if anyone do let me know
peer chaincode instantiate -n wallet -v 0 -C myc
peer chaincode invoke -n wallet  -c '{"Args":["initWallet","masterpass"]}' -C myc
peer chaincode invoke -n wallet  -c '{"Args":["getWalletInfo","masterWallet","masterpass"]}' -C myc
```

## Usage

### 1. Wallet Structure

```golang
type wallet struct {
        Name string `json:"name"`
        Owner string `json:"owner"`
        Balance float64 `json:"balance"`
        Password string `json:"password"`
}
```

### 2. Init Wallet

This function setup the first wallet with balance 100000 coins. All subsequent wallet will always be initialised with balance zero. It is the only source of coin in the system hence making number of coin in system constant.
 
args:- 
masterpass

example:- 
```bash
peer chaincode invoke -n wallet  -c '{"Args":["initWallet","masterpass"]}' -C myc
```

response:-
chaincode result: status:200 message:nil

error response:- chaincode result: status:500 message:"masteWallet already exists"

### 3. Create Wallet

Create wallet is used to initialise a new wallet.

args:- 1."name of wallet"  2."name of owner" 3."password"

example:-
```bash
peer chaincode invoke -n wallet  -c '{"Args":["createWallet","wallet1","mike","mikepass"]}' -C myc
```

respose:- 
Chaincode invoke successful. result: status:200

error response:-
chaincode result: status:500 message:"Wallet already exists"

### 4. transaction

Transact coin between wallets.

args:- 1. fromWallet 2. towallet 3. amount  4. password

example :-
```bash
peer chaincode invoke -n wallet  -c '{"Args":["transaction","masterWallet","wallet1","10","masterpass"]}' -C myc
```

response:-
Chaincode invoke successful. result: status:200

### 5. Get Wallet Info

get wallet details

args:- 1. walletname 2. password

example:-
```bash
peer chaincode invoke -n wallet  -c '{"Args":["getWalletInfo","masterWallet","masterpass"]}' -C myc
```

response: -
Chaincode invoke successful. result: status:200 payload:"{\"name\":\"masterWallet\",\"owner\":\"admin\",\"balance\":99990,\"password\":\"854E973753957848BD76991343F840B1A34784B7AAC33ECB0964C89D0A6FC8CC\"}"

## Build With

[hyperledger-fabric](https://www.hyperledger.org/projects/fabric) - a blockchain framework implementation and one of the Hyperledger projects hosted by The Linux Foundation.
[golang 1.10](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software

## Contributing

Coming soon

## Authors

- Parag Rahangdale

## licence

This project is licensed under the Apache-2.0 Licence

## Future Works

- [ ] find a way to properly deploy [not in dev mode] the chaincode on hyperledger Fabric
- [ ] Create a Nodejs client to interact with chaincode
- [ ] Create a frontend for the wallet