import FungibleToken from 0x{{.FungibleTokenContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}
import DapperUtilityCoin from 0x{{.DapperUtilityCoinContractAddress}}
import {{.NFTProductName}} from 0x{{.NFTContractAddress}}
import NFTStorefrontV2 from 0x{{.NFTStorefrontV2ContractAddress}}
import TokenForwarding from {{.TokenForwardingContractAddress}}


transaction() {
    // NFT IDs and buyback prices
    let nftIDs: [UInt64]
    let prices: [UFix64]

    let ducReceiver: Capability<&{FungibleToken.Receiver}>
    let royaltyReceiver: Capability<&{FungibleToken.Receiver}>
    let NFTProvider: Capability<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>?
    let storefront: auth(NFTStorefrontV2.CreateListing) &NFTStorefrontV2.Storefront
    var marketplacesCapability: [Capability<&{FungibleToken.Receiver}>]

    // 'customID' - Optional string to represent identifier of the dapp.
    let customID: String
    // 'commissionAmount' - Commission amount that will be taken away by the purchase facilitator i.e marketplacesAddress.
    let commissionAmount: UFix64
    // 'marketplacesAddress' - List of addresses that are allowed to get the commission.
    let marketplaceAddress: [Address]



    prepare(acct: auth(Storage, Capabilities) &Account) {
        // Initialize NFT IDs and buyback prices
        self.nftIDs = [{{.NFTIDs}}]
        self.prices = [{{.Prices}}]

        self.customID = "DAPPER_MARKETPLACE"
        self.commissionAmount = {{.SaleCommissionAmount}}
        self.marketplaceAddress = [{{.NFTContractAddress}}]
        self.marketplacesCapability = []

        // Validate the marketplaces capability before submiting to 'createListing'.
        for mp in self.marketplaceAddress {
            let marketplaceReceiver = getAccount(mp).capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)
            assert(marketplaceReceiver.borrow() != nil && marketplaceReceiver.borrow()!.isInstance(Type<@TokenForwarding.Forwarder>()), message: "Marketplaces does not possess the valid receiver type for DUC")
            self.marketplacesCapability.append(marketplaceReceiver!)
        }

        // *************************** Seller account interactions  *************************** //

        // This checks for the public capability
        if !acct.capabilities.get<&{{.NFTProductName}}.Collection>({{.NFTProductName}}.CollectionPublicPath)!.check() {
            acct.capabilities.unpublish({{.NFTProductName}}.CollectionPublicPath)
            acct.capabilities.publish(
                acct.capabilities.storage.issue<&{{.NFTProductName}}.Collection>({{.NFTProductName}}.CollectionStoragePath),
                at: {{.NFTProductName}}.CollectionPublicPath
            )
        }

        // If the account doesn't already have a Storefront
        if acct.storage.borrow<&NFTStorefrontV2.Storefront>(from: NFTStorefrontV2.StorefrontStoragePath) == nil {

            // Create a new empty .Storefront
            let storefront <- NFTStorefrontV2.createStorefront() as! @NFTStorefrontV2.Storefront
            
            // save it to the account
            acct.storage.save(<-storefront, to: NFTStorefrontV2.StorefrontStoragePath)

            // create a public capability for the .Storefront
            acct.capabilities.publish(
                acct.capabilities.storage.issue<&NFTStorefrontV2.Storefront>(NFTStorefrontV2.StorefrontStoragePath),
                at: NFTStorefrontV2.StorefrontPublicPath
            )
        }

        self.ducReceiver = acct.capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)!
            assert(self.ducReceiver.borrow() != nil, message: "Missing or mis-typed DUC receiver")
        
        self.royaltyReceiver = getAccount(0x4dfd62c88d1b6462).capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)!
            assert(self.royaltyReceiver.borrow() != nil, message: "Missing or mis-typed fungible token receiver for {{.NFTProductName}} account")
        
        // previous: let AllDayNFTCollectionProviderPrivatePath = /storage/AllDayNFTCollectionProviderForNFTStorefront
        let NFTCollectionProviderPrivatePath = /storage/{{.NFTProductName}}NFTCollectionProviderNFTStorefront

        // Temporary variable to handle capability assignment
        var provider: Capability<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>? =
            acct.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>>(from: NFTCollectionProviderPrivatePath)
        
        if provider == nil {
            provider = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>({{.NFTProductName}}.CollectionStoragePath)
            acct.capabilities.storage.getController(byCapabilityID: provider!.id)!.setTag("{{.NFTProductName}}NFTCollectionProviderForNFTStorefront")
            // Save the capability to the account storage
            acct.storage.save(provider!, to: NFTCollectionProviderPrivatePath)
        }

        self.NFTProvider = provider
        assert(self.NFTProvider?.borrow() != nil, message: "Missing or mis-typed {{.NFTProductName}}.Collection provider")

        self.storefront = acct.storage.borrow<auth(NFTStorefrontV2.CreateListing) &NFTStorefrontV2.Storefront>(from: NFTStorefrontV2.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefrontV2 Storefront")
    }


    pre {
        self.nftIDs.length == self.prices.length: "NFTs/prices length mismatch"
    }

    execute {
        for i, nftID in self.nftIDs {
            // List NFT for sale
            let saleCut = NFTStorefrontV2.SaleCut(
                receiver: self.ducReceiver,
                amount: self.prices[i] * 0.95
            )
            let royaltyCut = NFTStorefrontV2.SaleCut(
                receiver: self.royaltyReceiver,
                amount: self.prices[i] * 0.05
            )
            self.storefront.createListing(
                nftProviderCapability: self.NFTProvider!,
                nftType: Type<@{{.NFTProductName}}.NFT>(),
                nftID: nftID,
                salePaymentVaultType: Type<@DapperUtilityCoin.Vault>(),
                saleCuts: [saleCut, royaltyCut],
                marketplacesCapability: self.marketplacesCapability.length == 0 ? nil : self.marketplacesCapability,
                customID: self.customID,
                commissionAmount: self.commissionAmount,
                expiry: {{.Expiry}}
            )
        }
    }
}