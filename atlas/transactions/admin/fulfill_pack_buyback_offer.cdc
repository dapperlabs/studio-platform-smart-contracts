import FungibleToken from 0x{{.FungibleTokenContractAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenContractAddress}}
import DapperUtilityCoin from 0x{{.DapperUtilityCoinContractAddress}}
import {{.NFTProductName}} from 0x{{.NFTContractAddress}}
import NFTStorefront from 0x{{.NFTStorefrontV1ContractAddress}}

/// Fulfills a pack buyback offer: dapper pays DUC, user's NFTs are transferred to issuer via NFTStorefront.
///
transaction() {
    // NFT IDs and buyback prices
    let nftIDs: [UInt64]
    let prices: [UFix64]

    // User's storefront for listing NFTs
    let userStorefront: auth(NFTStorefront.CreateListing) &NFTStorefront.Storefront

    // Capability to withdraw user's NFTs for listing
    let userNFTWithdrawCap: Capability<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>?

    // Capability to receive DUC payment
    let userDUCReceiverCap: Capability<&{FungibleToken.Receiver}>

    // Dapper's DUC vault and initial balance
    let dapperVault: auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault
    let initialDapperBalance: UFix64

    // Issuer's NFT collection and address
    let issuerCollection: &NonFungibleToken.Collection
    let issuerAddress: Address

    prepare(user: auth(Storage, Capabilities) &Account, dapper: auth(BorrowValue) &Account) {
        // Initialize NFT IDs and buyback prices
        self.nftIDs = [{{.nftIDs}}]
        self.prices = [{{.prices}}]

        // Get or create user's storefront
        if user.storage.borrow<&NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath) == nil {
            user.storage.save(<- NFTStorefront.createStorefront(), to: NFTStorefront.StorefrontStoragePath)
            user.capabilities.publish(
                user.capabilities.storage.issue<&NFTStorefront.Storefront>(NFTStorefront.StorefrontStoragePath),
                at: NFTStorefront.StorefrontPublicPath
            )
        }
        self.userStorefront = user.storage.borrow<auth(NFTStorefront.CreateListing) &NFTStorefront.Storefront>(from: NFTStorefront.StorefrontStoragePath)
            ?? panic("Missing NFTStorefront")

        // Get or create user's NFT withdrawal capability
        let providerPath = /storage/{{.NFTProductName}}NFTCollectionProviderForNFTStorefront
        self.userNFTWithdrawCap = user.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>>(from: providerPath)
        if self.userNFTWithdrawCap == nil {
            self.userNFTWithdrawCap = user.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{{.NFTProductName}}.Collection>({{.NFTProductName}}.CollectionStoragePath)
            user.capabilities.storage.getController(byCapabilityID: self.userNFTWithdrawCap!.id)!.setTag("{{.NFTProductName}}CollectionProviderForNFTStorefront")
            user.storage.save(self.userNFTWithdrawCap!, to: providerPath)
        }

        // Get user's DUC receiver capability
        self.userDUCReceiverCap = user.capabilities.get<&{FungibleToken.Receiver}>(/public/dapperUtilityCoinReceiver)!
        assert(self.userDUCReceiverCap.borrow() != nil, message: "Missing DUC receiver cap")

        // Borrow issuer's NFT collection
        self.issuerCollection = issuer.storage.borrow<&NonFungibleToken.Collection>(from: {{.NFTProductName}}.CollectionStoragePath)
            ?? panic("Missing issuer NFT collection")

        // Borrow Dapper's DUC vault
        self.dapperVault = dapper.storage.borrow<auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault>(from: /storage/dapperUtilityCoinVault)
            ?? panic("Missing Dapper DUC vault")

        // Record Dapper's DUC balance for post-condition
        self.initialDapperBalance = self.dapperVault.balance

        // Record issuer's address for pre-condition
        self.issuerAddress = issuer.address
    }

    pre {
        self.nftIDs.length == self.prices.length: "NFTs/prices length mismatch"
        self.issuerAddress == 0x{{.NFTContractAddress}}: "Must be NFT contract owner"
    }

    execute {
        // Gather active listings in user's storefront
        let listingIDsByNFTID: {UInt64: UInt64} = {}
        for id in self.userStorefront.getListingIDs() {
            let details = self.userStorefront.borrowListing(listingResourceID: id)!.getDetails()
            if !details.purchased {
                listingIDsByNFTID[details.nftID] = id
            }
        }

        // For each NFT, list and buy back at the specified price
        for i, nftID in self.nftIDs {
            // Ensure no active listing exists for this NFT
            assert(!listingIDsByNFTID.containsKey(nftID),
                message: "NFT ".concat(nftID.toString()).concat(" already listed, id=").concat(listingIDsByNFTID[nftID]!.toString()))

            // List NFT for sale
            let listingID = self.userStorefront.createListing(
                nftProviderCapability: self.userNFTWithdrawCap!,
                nftType: Type<@{{.NFTProductName}}.NFT>(),
                nftID: nftID,
                salePaymentVaultType: Type<@DapperUtilityCoin.Vault>(),
                saleCuts: [NFTStorefront.SaleCut(receiver: self.userDUCReceiverCap, amount: self.prices[i])]
            )

            // Dapper purchases NFT and deposits in issuer's collection
            let listing = self.userStorefront.borrowListing(listingResourceID: listingID)!
            self.issuerCollection.deposit(token: <- listing.purchase(payment: <- self.dapperVault.withdraw(amount: self.prices[i])))
        }
    }

    post {
        self.dapperVault.balance == self.initialDapperBalance: "DUC balance must be restored (externally) before transaction completes"
    }
}
