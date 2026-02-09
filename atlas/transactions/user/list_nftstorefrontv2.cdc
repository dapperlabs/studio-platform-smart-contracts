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
    // 'commissionAmounts' - Commission amount that will be taken away by the purchase facilitator i.e marketplacesAddress.
    let commissionAmounts: [UFix64]
    // 'marketplacesAddress' - List of addresses that are allowed to get the commission.
    let marketplaceAddress: [Address]
    // we only ever want to use DapperUtilityCoin
    let universalDucReceiver: Address


    prepare(acct: auth(Storage, Capabilities) &Account) {
        // Initialize NFT IDs and buyback prices
        self.nftIDs = [{{.NFTIDsString}}]
        self.prices = [{{.PricesString}}]

        self.customID = "DAPPER_MARKETPLACE"
        self.marketplaceAddress = [0x{{.NFTContractAddress}}]
        // we only ever want to use DapperUtilityCoin
        self.universalDucReceiver = 0x{{.DapperUtilityCoinContractAddress}}


        self.allSaleCuts = []                                                                                                                                                                                                     
        self.marketplacesCapability = []
        self.commissionAmounts = []

        let commissionPercent = {{.SaleCommissionPercentString}}

        let collectionDataOpt = {{.NFTProductName}}.resolveContractView(resourceType: Type<@{{.NFTProductName}}.NFT>(), viewType: Type<MetadataViews.NFTCollectionData>())
            ?? panic("Missing collection data")
        let collectionData = collectionDataOpt as! MetadataViews.NFTCollectionData


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
        if !acct.capabilities.get<&{{.NFTProductName}}.Collection>(collectionData.publicPath)!.check() {
            acct.capabilities.unpublish(collectionData.publicPath)
            acct.capabilities.publish(
                acct.capabilities.storage.issue<&{{.NFTProductName}}.Collection>(collectionData.storagePath),
                at: collectionData.publicPath
            )
        }

        var nftProviderCap: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>? = nil
        // check if there is an existing capability/capability controller for the storage path
        let nftCollectionControllers = acct.capabilities.storage.getControllers(forPath: collectionData.storagePath)
        for controller in nftCollectionControllers {
            if let maybeProviderCap = controller.capability as? Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>? {
                nftProviderCap = maybeProviderCap
                break
            }
        }

        // if there are no capabilities created for that storage path
        // or if existing capability is no longer valid, issue a new one
        if nftProviderCap == nil || nftProviderCap?.check() ?? false {
            nftProviderCap = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
                collectionData.storagePath
            )
        }
        assert(nftProviderCap?.check() ?? false, message: "Could not assign Provider Capability")

        self.nftProvider = nftProviderCap!

        let collectionRef = acct
        .capabilities.borrow<&{{.NFTProductName}}.Collection>(collectionData.publicPath)
        ?? panic("Could not borrow a reference to the collection")



        // Pre-calculate sale cuts for each NFT                                                                                                                                                                                   
        for i, nftID in self.nftIDs {
            var saleCutsForThisNFT: [NFTStorefrontV2.SaleCut] = []         
            var totalRoyaltyCut = 0.0
            var commissionAmount = commissionPercent * self.prices[i]

            let nft = collectionRef.borrowNFT(self.nftIDs[i])!
            let effectiveSaleItemPrice = self.prices[i] - commissionAmount

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
            self.commissionAmounts.append(commissionAmount)
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
                commissionAmount: self.commissionAmounts[i] ,
                expiry: {{.Expiry}}
            )
        }
    }
}