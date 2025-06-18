import FungibleToken from 0x{{.FungibleTokenContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}
import DapperUtilityCoin from 0x{{.DapperUtilityCoinContractAddress}}
import PackNFT from 0x{{.NFTContractAddress}}
import NFTStorefront from 0x{{.NFTStorefrontV1ContractAddress}}

/// Transaction facilitates the purchase of listed NFT.
/// It takes the storefront address, listing resource that need
/// to be purchased & a address that will takeaway the commission.
///
/// Buyer of the listing (,i.e. underling NFT) would authorize and sign the
/// transaction and if purchase happens then transacted NFT would store in
/// buyer's collection.

transaction() {
    let paymentVault: @{FungibleToken.Vault}
    let PackNFTCollection: &{NonFungibleToken.CollectionPublic}
    let storefront: &NFTStorefront.Storefront
    let listings: [&{NFTStorefront.ListingPublic}]

    prepare(universalDucPayer: auth(Storage, Capabilities) &Account, user: auth(Storage, Capabilities) &Account) {
        let sellerAddress: Address = 0x{{.StorefrontAddress}}
        let buyerAddress: Address = 0x{{.RecipientAddress}}
        let listingIds: [UInt64] = [{{.ListingIds}}]
        let buyer = getAccount(buyerAddress)

        assert(sellerAddress != buyerAddress, message : "Buyer and seller can not be same")

        // Ensure the user has an AllDay collection set up
        if user.storage.borrow<&AllDay.Collection>(from: AllDay.CollectionStoragePath) == nil {
            // Create a new collection and save it to the account storage
            user.storage.save(<- AllDay.createEmptyCollection(nftType: Type<@AllDay.NFT>()), to: AllDay.CollectionStoragePath)

            // Create a public capability for the collection
            user.capabilities.unpublish(AllDay.CollectionPublicPath)
            user.capabilities.publish(
                user.capabilities.storage.issue<&AllDay.Collection>(AllDay.CollectionStoragePath),
                at: AllDay.CollectionPublicPath
            )
        }
        // Ensure the user has a PackNFT collection set up
        if user.storage.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath) == nil {
            // Create a new collection and save it to the account storage
            user.storage.save(<- PackNFT.createEmptyCollection(nftType: Type<@PackNFT.NFT>()), to: PackNFT.CollectionStoragePath)

            // Create a public capability for the collection
            user.capabilities.unpublish(PackNFT.CollectionPublicPath)
            user.capabilities.publish(
                user.capabilities.storage.issue<&PackNFT.Collection>(PackNFT.CollectionStoragePath),
                at: PackNFT.CollectionPublicPath
            )
        }

        // Access the storefront public resource of the seller to purchase the listing.
        self.storefront = getAccount(sellerAddress)
            .capabilities.borrow<&NFTStorefront.Storefront>(
                NFTStorefront.StorefrontPublicPath
            )
            ?? panic("Could not borrow Storefront from provided address")

        var salePrice: UFix64 = 0.0
        self.listings = []
        for listingID in listingIds {
            let listing = self.storefront.borrowListing(listingResourceID: listingID)
                ?? panic("No Listing with that ID in Storefront")
            self.listings.append(
                listing
            )
            salePrice = salePrice + listing.getDetails().salePrice
        }

        // Borrow a provider reference to the buyers vault
		let provider = universalDucPayer.storage.borrow<auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault>(from: /storage/dapperUtilityCoinVault)
		?? panic("Cannot borrow DUC vault from buyer account storage")

        // withdraw the purchase tokens from the vault
        self.paymentVault <- provider.withdraw(amount: salePrice)

        // Access the buyer's NFT collection to store the purchased NFT.
        self.PackNFTCollection = buyer.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath)
        ?? panic("Could not borrow Storefront from provided address")

    }

    execute {
        // Purchase the NFTs
        for listing in self.listings {
            let listingSalePrice: UFix64 = listing.getDetails().salePrice
            let listingPaymentVault <-self.paymentVault.withdraw(amount: listingSalePrice)
            let item <- listing.purchase(
                payment: <-listingPaymentVault
            )

            self.PackNFTCollection.deposit(token: <-item)
        }
        destroy self.paymentVault
    }
}