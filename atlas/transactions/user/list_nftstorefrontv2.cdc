import FungibleToken from 0x{{.FungibleTokenContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}
import DapperUtilityCoin from 0x{{.DapperUtilityCoinContractAddress}}
import {{.NFTProductName}} from 0x{{.NFTContractAddress}}
import NFTStorefrontV2 from 0x{{.NFTStorefrontV2ContractAddress}}
import MetadataViews from 0x{{.MetadataViewsAddress}}  
import TokenForwarding from 0x{{.TokenForwardingContractAddress}} 

transaction() {
    // NFT IDs and buyback prices
    let nftIDs: [UInt64]
    let prices: [UFix64]

    var ftReceiver: Capability<&{FungibleToken.Receiver}>?
    let nftProvider: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>?
    let storefront: auth(NFTStorefrontV2.CreateListing) &NFTStorefrontV2.Storefront

    var allSaleCuts: [[NFTStorefrontV2.SaleCut]]                                                                                                                                                                                  
    var marketplacesCapability: [Capability<&{FungibleToken.Receiver}>]

    // 'customID' - Optional string to represent identifier of the dapp.
    let customID: String
    // 'commissionAmount' - Commission amount that will be taken away by the purchase facilitator i.e marketplacesAddress.
    let commissionAmount: UFix64
    // 'marketplacesAddress' - List of addresses that are allowed to get the commission.
    let marketplaceAddress: [Address]
    // we only ever want to use DapperUtilityCoin
    let universalDucReceiver: Address


    prepare(acct: auth(Storage, Capabilities) &Account) {
        // Initialize NFT IDs and buyback prices
        self.nftIDs = [{{.NFTIDs}}]
        self.prices = [{{.Prices}}]

        self.customID = "DAPPER_MARKETPLACE"
        self.commissionAmount = {{.SaleCommissionAmount}}
        self.marketplaceAddress = [0x{{.NFTContractAddress}}]
        // we only ever want to use DapperUtilityCoin
        self.universalDucReceiver = 0x{{.DapperUtilityCoinContractAddress}}


        self.allSaleCuts = []                                                                                                                                                                                                     
        self.marketplacesCapability = []


        // ************************* Handling of DUC Recevier *************************** //
        
        // Fetch the capability of the universal DUC receiver
        let recipient = getAccount(self.universalDucReceiver).capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)!
        assert(recipient.borrow() != nil, message: "Missing or mis-typed Fungible Token receiver for the DUC recipient")

        self.ftReceiver = acct.capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)

        // Validate the marketplaces capability before submiting to 'createListing'
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

        let PrivateCollectionPath = /storage/{{.NFTProductName}}CollectionProviderForNFTStorefront

        // Temporary variable to handle capability assignment
        var provider: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>? =
            acct.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>>(from: PrivateCollectionPath)

        if provider == nil {
            provider = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>({{.NFTProductName}}.CollectionStoragePath)
            acct.capabilities.storage.getController(byCapabilityID: provider!.id)!.setTag("{{.NFTProductName}}CollectionProviderForNFTStorefront")
            // Save the capability to the account storage
            acct.storage.save(provider!, to: PrivateCollectionPath)
        }

        self.nftProvider = provider
        assert(self.nftProvider?.borrow() != nil, message: "Missing or mis-typed {{.NFTProductName}}.Collection provider")


        let collectionRef = acct
        .capabilities.borrow<&{{.NFTProductName}}.Collection>({{.NFTProductName}}.CollectionPublicPath)
        ?? panic("Could not borrow a reference to the collection")



        // Pre-calculate sale cuts for each NFT                                                                                                                                                                                   
        for i, nftID in self.nftIDs {
            var saleCutsForThisNFT: [NFTStorefrontV2.SaleCut] = []         
            var totalRoyaltyCut = 0.0

            let nft = collectionRef.borrowNFT(self.nftIDs[i])!
            let effectiveSaleItemPrice = self.prices[i] - self.commissionAmount

            // Check whether the NFT implements the MetadataResolver or not.
            if nft.getViews().contains(Type<MetadataViews.Royalties>()) {
                let royaltiesRef = nft.resolveView(Type<MetadataViews.Royalties>()) ?? panic("Unable to retrieve the royalties")
                let royalties = (royaltiesRef as! MetadataViews.Royalties).getRoyalties()
                for royalty in royalties {
                    let royaltyReceiver = royalty.receiver
                    assert(royaltyReceiver.borrow() != nil && royaltyReceiver.borrow()!.isInstance(Type<@TokenForwarding.Forwarder>()), message: "Royalty receiver does not have a valid receiver type")

                    let royaltyAmount = royalty.cut * effectiveSaleItemPrice         
                    saleCutsForThisNFT.append(NFTStorefrontV2.SaleCut(receiver: royalty.receiver, amount: royaltyAmount))
                    totalRoyaltyCut = totalRoyaltyCut + royaltyAmount
                }
            }
            
            // Add seller cut (remaining after royalties)                                                                                                                                                                         
            saleCutsForThisNFT.append(NFTStorefrontV2.SaleCut(                                                                                                                                                                    
                receiver: self.ftReceiver!,                                                                                                                                                                                        
                amount: effectiveSaleItemPrice - totalRoyaltyCut                                                                                                                                                                  
            ))                                                                                                                                                                                                                    
                                                                                                                                                                                                                                
            self.allSaleCuts.append(saleCutsForThisNFT)     
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

        self.storefront = acct.storage.borrow<auth(NFTStorefrontV2.CreateListing) &NFTStorefrontV2.Storefront>(from: NFTStorefrontV2.StorefrontStoragePath)
            ?? panic("Missing or mis-typed NFTStorefrontV2 Storefront")
    }


    pre {
        self.nftIDs.length == self.prices.length: "NFTs/prices length mismatch"
        self.allSaleCuts.length == self.nftIDs.length: "Sale cuts/NFTs length mismatch"                                                                                                                                           
    }

    execute {
        for i, nftID in self.nftIDs {
            // List NFT for sale
            self.storefront.createListing(
                nftProviderCapability: self.nftProvider!,
                nftType: Type<@{{.NFTProductName}}.NFT>(),
                nftID: nftID,
                salePaymentVaultType: Type<@DapperUtilityCoin.Vault>(),
                saleCuts: self.allSaleCuts[i],
                marketplacesCapability: self.marketplacesCapability.length == 0 ? nil : self.marketplacesCapability,
                customID: self.customID,
                commissionAmount: self.commissionAmount,
                expiry: {{.Expiry}}
            )
        }
    }
}