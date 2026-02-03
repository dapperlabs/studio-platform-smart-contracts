import NFTStorefront from 0x94b06cfca1d8a476

// removes an NFT listing from NFTStorefront
transaction() {
    let storefront: auth(NFTStorefront.RemoveListing) &NFTStorefront.Storefront

    prepare(acct: auth(Storage) &Account) {
        self.storefront = acct.storage.borrow<auth(NFTStorefront.RemoveListing) &NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefront.Storefront")
    }

    execute {
        self.storefront.removeListing(listingResourceID: {{.listingResourceID}})
    }
}