import NFTStorefrontV2 from 0x{{.NFTStorefrontV2Address}}

// removes an NFT listing from NFTStorefrontV2
transaction() {
    let storefront: auth(NFTStorefrontV2.RemoveListing) &NFTStorefrontV2.Storefront

    prepare(acct: auth(Storage) &Account) {
        self.storefront = acct.storage.borrow<auth(NFTStorefrontV2.RemoveListing) &NFTStorefrontV2.Storefront>(from: NFTStorefrontV2.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefrontV2.Storefront resource")
    }

    execute {
        let listingResourceIDs: [Uint64] = [{{.ListingResourceIDs}}]
        for resourceID in listingResourceIDs {
            self.storefront.removeListing(listingResourceID: resourceID)
        }
    }
}z