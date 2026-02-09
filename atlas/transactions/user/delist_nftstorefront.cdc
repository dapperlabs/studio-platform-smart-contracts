import NFTStorefront from 0x{{.NFTStorefrontAddress}}

// removes an NFT listing from NFTStorefront
transaction() {
    let storefront: auth(NFTStorefront.RemoveListing) &NFTStorefront.Storefront

    prepare(acct: auth(Storage) &Account) {
        self.storefront = acct.storage.borrow<auth(NFTStorefront.RemoveListing) &NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefront.Storefront")
    }

    execute {
        let listingResourceIDs: [Uint64] = [{{.ListingResourceIDs}}]
        for resourceID in listingResourceIDs {
            self.storefront.removeListing(listingResourceID: resourceID)
        }
    }
}