import NFTStorefrontV2 from 0x2d55b98eb200daef

// removes an NFT listing from NFTStorefrontV2

transaction() {
    let storefront: auth(NFTStorefrontV2.RemoveListing) &NFTStorefrontV2.Storefront

    prepare(acct: auth(Storage) &Account) {
        self.storefront = acct.storage.borrow<auth(NFTStorefrontV2.RemoveListing) &NFTStorefrontV2.Storefront>(from: NFTStorefrontV2.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefrontV2.Storefront resource")
    }

    execute {
        self.storefront.removeListing(listingResourceID: {{.listingResourceID}})
    }
}z