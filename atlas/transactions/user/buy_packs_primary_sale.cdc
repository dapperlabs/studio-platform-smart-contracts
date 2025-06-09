import FungibleToken from 0x{{.FungibleTokenContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}
import DapperUtilityCoin from 0x{{.DapperUtilityCoinContractAddress}}
import {{.NFTProductName}}, AllDay from 0x{{.NFTContractAddress}}
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
    let {{.NFTProductName}}Collection: &{NonFungibleToken.CollectionPublic}
    let storefront: &NFTStorefront.Storefront
    let listings: [&{NFTStorefront.ListingPublic}]

    prepare(universalDucPayer: auth(Storage, Capabilities) &Account) {
        let sellerAddress: Address = 0x{{.StorefrontAddress}}
        let buyerAddress: Address = 0x{{.RecipientAddress}}
        let listingIds: [UInt64] = [{{.ListingIds}}]
        let buyer = getAccount(buyerAddress)

        assert(sellerAddress != buyerAddress, message : "Buyer and seller can not be same")
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
        self.{{.NFTProductName}}Collection = buyer.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>({{.NFTProductName}}.CollectionPublicPath)
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

            self.{{.NFTProductName}}Collection.deposit(token: <-item)
        }
        destroy self.paymentVault
    }
}