# DSSCollection Smart Contract

## NonFungibleToken Contract Addresses
| Network   | Address |
| ----------|:-------:|
| Emulator   |  0xf8d6e0586b0a20c7  |
| Testnet   |  0x860da69cd099c1d2  |
| Mainnet   |  0x1d7e57aa55817448   |

## DSSCollection Contract Addresses
| Network   | Address |                                                                                     |
| ----------|:-------:|-------------------------------------------------------------------------------------|
| Testnet   |  0x9dae79d85e17d382   | [Flow View Source](https://flow-view-source.com/testnet/) |
| Mainnet   |  XXXX   | [Flow View Source](https://flow-view-source.com/mainnet/) |


## Transactions
- MintNFT: Mints a nft from the collection group.

## Flow Integration

### Setup
You will need to create a flow.json file to hold accounts and contracts. This file is in .gitignore because it contains private keys. Please see the author of this README for file, or run flo-init and fill in yourself with the output of these scripts.
```
flow emulator

flow keys generate

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
flow transactions send ./transactions/admin/create_collection_group.cdc "collection name" /public/AllDayNFTCollection --signer emulator-account

// create collection group with time range
flow transactions send ./transactions/admin/create_collection_group_time_range.cdc "collection name" /public/AllDayNFTCollection 1673299041.0 1674681434.0 --signer emulator-account

// mint
flow transactions send ./transactions/admin/mint_nft.cdc 0xf8d6e0586b0a20c7 1 "houseofhufflepuff"  --signer emulator-account

// close collection group
flow transactions send ./transactions/admin/close_collection_group.cdc 1 --signer emulator-account

// add item to slot
flow transactions send ./transactions/admin/create_item_in_slot.cdc 100 10 "edition.id" 1 --signer emulator-account
```


### Scripts
```
flow scripts execute scripts/total_supply.cdc

flow scripts execute scripts/get_collection_group.cdc 1
```


