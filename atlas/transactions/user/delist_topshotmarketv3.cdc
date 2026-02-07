import TopShot from 0x{{.TopShotContractAddress}}
import Market from 0x{{.TopShotMarketContractAddress}}
import TopShotMarketV3 from 0x{{.TopShotMarketContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}

// This transaction is for a user to stop a moment sale in their account

// Parameters
//
// tokenID: the ID of the moment whose sale is to be delisted

transaction() {

    let nftIds: [UInt64]

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.nftIds = {{.NftIds}}

        // borrow a reference to the owner's sale collection
        if let topshotSaleV3Collection = acct.storage.borrow<auth(TopShotMarketV3.Cancel) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) {

            // cancel the moments from the sale, thereby de-listing it
            for nftId in self.nftIds {
                topshotSaleV3Collection.cancelSale(tokenID: nftId)
            }
            
        } else if let topshotSaleCollection = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &Market.SaleCollection>(from: /storage/topshotSaleCollection) {
            // Borrow a reference to the NFT collection in the signers account
            let collectionRef = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
                ?? panic("Could not borrow from MomentCollection in storage")
        
            // withdraw the moments from the sale, thereby de-listing it
            for nftId in self.nftIds {
                let token <- topshotSaleCollection.withdraw(tokenID: nftId)
                
                // deposit the moment into the owner's collection
                collectionRef.deposit(token: <-token)
            }

        }
    }
}