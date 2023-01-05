# DSSCollection Smart Contract

## DSSCollection Contract Addresses
| Network   | Address |                                                                                     |
| ----------|:-------:|-------------------------------------------------------------------------------------|
| Testnet   |  XXXX   | [Flow View Source](https://flow-view-source.com/testnet/) |
| Mainnet   |  XXXX   | [Flow View Source](https://flow-view-source.com/mainnet/) |


## Transactions
- MintNFT: Mints a nft from the collection group.

## Flow Integration

### Setup
You will need to create a flow.json file to hold accounts and contracts. This file is in .gitignore because it contains private keys. Please see the author of this README for file, or run flo-init and fill in yourself with the output of these scripts.
```
flow emulator --persist

flow accounts create --key 2488befee0c3873d2a0dcd291ad27d0cd061a750dc5137af3f3ffee7ff70f528a229ed1039ce1de23986d4506d4671df096881e5d60c18e93c8df321a180adac

flow project deploy --network=emulator --update
```

### Mint NFT
```
flow transactions send ./transactions/mintBrand.cdc 0x179b6b1cb6755e31 "nft-1"  --signer brand-acct
```

### Check Total Supply
```
flow scripts execute scripts/totalSupply.cdc
```
### Create Collections
```
// create collection for max-acct /storage/brandCollection
flow transactions send ./transactions/setupAccount.cdc --signer max-acct

```

