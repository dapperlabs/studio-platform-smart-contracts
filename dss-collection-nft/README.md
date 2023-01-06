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
flow emulator

flow accounts create --key 2488befee0c3873d2a0dcd291ad27d0cd061a750dc5137af3f3ffee7ff70f528a229ed1039ce1de23986d4506d4671df096881e5d60c18e93c8df321a180adac

flow project deploy --network=emulator --update
```

### Create Collections
```
// create & save DSSCollection 
flow transactions send ./transactions/setup_collection.cdc --signer emulator-account
```

### Admin
```
// create collection group
flow transactions send ./transactions/admin/create_collection_group.cdc "collection name" "product-name" --signer emulator-account

// mint
flow transactions send ./transactions/admin/mint_nft.cdc 0xf8d6e0586b0a20c7 1 "jer.ahrens"  --signer emulator-account

```


### Scripts
```
flow scripts execute scripts/totalSupply.cdc
```

