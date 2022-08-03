# Edition Nft Smart Contracts

## AllDay Seasonal Contract Addresses
| Network   | Address |              |
| ----------|:-------:| -------------|
| Testnet   |  0x4dbd1602c43aae03   | [Flow View Source](https://flow-view-source.com/mainnet/account/0x4dbd1602c43aae03) |
| Mainnet   |  XXXX   | [Flow View Source](https://flow-view-source.com/testnet/account/0xe4cf4bdc1751c65d) |

## Entities

### Editions
Edition contain the metadata.

By default, edition is active. If it is closed, nfts can not be minted. An edition can be closed under following condition:
- The CloseEdition transaction is used

**Fields**
- ID
- Active
- Metadata

**Transactions**
- CreateEdition: Mints a new Edition on Flow.
- CloseEdition: Closes an Edition so no new moments can be minted from it. This is irreversible. The Edition is closed by setting active to be false.

### NFT
Nfts are minted out of editions. You can think of Editions as a "cookie cutter" for nfts. 

**Fields**
- ID
- EditionID

**Transactions**
- MintNFT: Mints a nft out of an edition
