import NonFungibleToken from "NonFungibleToken"
import IPackNFT from "IPackNFT"

/// The Pack Distribution Service (PDS) contract is responsible for creating and managing distributions of packs.
///
access(all) contract PDS{
    /// Entitlement that grants the ability to operate PDS functionalities.
    ///
    access(all) entitlement Operate

    access(all) var version: String
    access(all) let PackIssuerStoragePath: StoragePath
    access(all) let PackIssuerCapRecv: PublicPath
    access(all) let DistCreatorStoragePath: StoragePath
    access(all) let DistManagerStoragePath: StoragePath

    /// The next distribution ID to be used.
    ///
    access(all) var nextDistId: UInt64

    /// Dictionary that stores distribution IDs to distribution details in the contract state.
    ///
    access(contract) let Distributions: {UInt64: DistInfo}

    /// Dictionary that stores distribution IDs to shared capabilities in the contract state.
    ///
    access(contract) let DistSharedCap: @{UInt64: SharedCapabilities}

    /// Emitted when an issuer has created a distribution.
    ///
    access(all) event DistributionCreated(DistId: UInt64, title: String, metadata: {String: String}, state: UInt8)

    /// Emmitted when a distribution manager has updated a distribution state.
    ///
    access(all) event DistributionStateUpdated(DistId: UInt64, state: UInt8)

    /// Enum that defines the status of a Distribution.
    ///
    access(all) enum DistState: UInt8 {
        access(all) case Initialized
        access(all) case Invalid
        access(all) case Complete
    }

    /// Struct that defines the details of a Distribution.
    ///
    access(all) struct DistInfo {
        access(all) let title: String
        access(all) let metadata: {String: String}
        access(all) var state: PDS.DistState

        access(contract) fun setState(newState: PDS.DistState) {
            self.state = newState
        }

        /// DistInfo struct initializer.
        ///
        view init(title: String, metadata: {String: String}) {
            self.title = title
            self.metadata = metadata
            self.state = PDS.DistState.Initialized
        }
    }

    /// Struct that defines a Collectible.
    ///
    access(all) struct Collectible: IPackNFT.Collectible {
        access(all) let address: Address
        access(all) let contractName: String
        access(all) let id: UInt64

        // returning in string so that it is more readable and anyone can check the hash
        access(all) view fun hashString(): String {
            // address string is 16 characters long with 0x as prefix (for 8 bytes in hex)
            // example: ,f3fcd2c1a78f5ee.ExampleNFT.12
            let c = "A."
            var a = ""
            let addrStr = self.address.toString()
            if addrStr.length < 18 {
                let padding = 18 - addrStr.length
                let p = "0"
                var i = 0
                a = addrStr.slice(from: 2, upTo: addrStr.length)
                while i < padding {
                    a = p.concat(a)
                    i = i + 1
                }
            } else {
                a = addrStr.slice(from: 2, upTo: 18)
            }
            return c.concat(a).concat(".").concat(self.contractName).concat(".").concat(self.id.toString())
        }

        /// Collectible struct initializer.
        ///
        view init(address: Address, contractName: String, id: UInt64) {
            self.address = address
            self.contractName = contractName
            self.id = id
        }
    }

    /// Resource that defines the shared capabilities required for creating and managing Pack NFTs.
    ///
    access(all) resource SharedCapabilities {
        /// Capability to withdraw NFTs from the issuer.
        ///
        access(self) let withdrawCap: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>

        /// Capability to mint, reveal, and open Pack NFTs.
        ///
        access(self) let operatorCap: Capability<auth(IPackNFT.Operate) &{IPackNFT.IOperator}>

        /// Withdraw an NFT from the issuer.
        ///
        access(contract) fun withdrawFromIssuer(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let c = self.withdrawCap.borrow() ?? panic("no such cap")
            return <- c.withdraw(withdrawID: withdrawID)
        }

        /// Mint Pack NFTs.
        ///
        access(contract) fun mintPackNFT(distId: UInt64, commitHashes: [String], issuer: Address, recvCap: &{NonFungibleToken.CollectionPublic}) {
            var i = 0
            let c = self.operatorCap.borrow() ?? panic("no such cap")
            while i < commitHashes.length{
                let nft <- c.mint(distId: distId, commitHash: commitHashes[i], issuer: issuer)
                i = i + 1
                let n <- nft
                recvCap.deposit(token: <- n)
            }
        }

        /// Reveal Pack NFTs.
        ///
        access(contract) fun revealPackNFT(packId: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
            let c = self.operatorCap.borrow() ?? panic("no such cap")
            c.reveal(id: packId, nfts: nfts, salt: salt)
        }

        /// Open Pack NFTs.
        ///
        access(contract) fun openPackNFT(packId: UInt64, nfts: [{IPackNFT.Collectible}], recvCap: &{NonFungibleToken.CollectionPublic}, collectionStoragePath: StoragePath?) {
            let c = self.operatorCap.borrow() ?? panic("no such cap")
            let toReleaseNFTs: [UInt64] = []
            var i = 0
            while i < nfts.length {
                toReleaseNFTs.append(nfts[i].id)
                i = i + 1
            }
            c.open(id: packId, nfts: nfts)
            if collectionStoragePath == nil {
                self.fulfillFromIssuer(nftIds: toReleaseNFTs, recvCap: recvCap)
            } else {
                self.releaseEscrow(nftIds: toReleaseNFTs, recvCap: recvCap , collectionStoragePath: collectionStoragePath!)
            }
        }

        /// Release escrowed NFTs to the receiver.
        ///
        access(contract) fun releaseEscrow(nftIds: [UInt64], recvCap: &{NonFungibleToken.CollectionPublic}, collectionStoragePath: StoragePath) {
            let pdsCollection = PDS.account.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(from: collectionStoragePath)
                ?? panic("Unable to borrow PDS collection provider capability from private path")
            var i = 0
            while i < nftIds.length {
                recvCap.deposit(token: <- pdsCollection.withdraw(withdrawID: nftIds[i]))
                i = i + 1
            }
        }

        /// Release NFTs from the issuer to the receiver.
        ///
        access(contract) fun fulfillFromIssuer(nftIds: [UInt64], recvCap:  &{NonFungibleToken.CollectionPublic}) {
            let issuerCollection = self.withdrawCap.borrow() ?? panic("Unable to borrow withdrawCap")
            var i = 0
            while i < nftIds.length {
                recvCap.deposit(token: <- issuerCollection.withdraw(withdrawID: nftIds[i]))
                i = i + 1
            }
        }

        /// SharedCapabilities resource initializer.
        ///
        view init(
            withdrawCap: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>,
            operatorCap: Capability<auth(IPackNFT.Operate) &{IPackNFT.IOperator}>
        ) {
            self.withdrawCap = withdrawCap
            self.operatorCap = operatorCap
        }
    }


    // Included for backwards compatibility.
    access(all) resource interface PackIssuerCapReciever {}

    /// Resource that defines the issuer of a pack.
    ///
    access(all) resource PackIssuer: PackIssuerCapReciever {
        access(self) var cap: Capability<&DistributionCreator>?

        /// Set the capability to create a distribution; the function is publicly accessible but requires a capability argument to a DistrubutionCreator admin resource.
        ///
        access(all) fun setDistCap(cap: Capability<&DistributionCreator>) {
            pre {
                cap.check(): "Invalid capability"
            }
            self.cap = cap
        }

        access(Operate) fun createDist(sharedCap: @SharedCapabilities, title: String, metadata: {String: String}) {
            assert(title.length > 0, message: "Title must not be empty")
            let c = self.cap!.borrow()!
            c.createNewDist(sharedCap: <- sharedCap, title: title, metadata: metadata)
        }

        /// PackIssuer resource initializer.
        ///
        view init() {
            self.cap = nil
        }
    }

    // Included for backwards compatibility.
    access(all) resource interface IDistCreator {}

    /// Resource that defines the creator of a distribution.
    ///
    access(all) resource DistributionCreator: IDistCreator {
        access(contract) fun createNewDist(sharedCap: @SharedCapabilities, title: String, metadata: {String: String}) {
            let currentId = PDS.nextDistId
            PDS.DistSharedCap[currentId] <-! sharedCap
            PDS.Distributions[currentId] = DistInfo(title: title, metadata: metadata)
            PDS.nextDistId = currentId + 1
            emit DistributionCreated(DistId: currentId, title: title, metadata: metadata, state: 0)
        }
    }

    /// Resource that defines the manager of a distribution.
    ///
    access(all) resource DistributionManager {
        access(Operate) fun updateDistState(distId: UInt64, state: PDS.DistState) {
            let d = PDS.Distributions.remove(key: distId) ?? panic ("No such distribution")
            d.setState(newState: state)
            PDS.Distributions.insert(key: distId, d)
            emit DistributionStateUpdated(DistId: distId, state: state.rawValue)
        }

        access(Operate) fun withdraw(distId: UInt64, nftIDs: [UInt64], escrowCollectionPublic: PublicPath) {
            assert(PDS.DistSharedCap.containsKey(distId), message: "No such distribution")
            let d <- PDS.DistSharedCap.remove(key: distId)!
            let pdsCollection = PDS.getManagerCollectionCap(escrowCollectionPublic: escrowCollectionPublic).borrow()!
            var i = 0
            while i < nftIDs.length {
                let nft <- d.withdrawFromIssuer(withdrawID: nftIDs[i])
                pdsCollection.deposit(token:<-nft)
                i = i + 1
            }
            PDS.DistSharedCap[distId] <-! d
        }

        access(Operate) fun mintPackNFT(distId: UInt64, commitHashes: [String], issuer: Address, recvCap: &{NonFungibleToken.CollectionPublic}) {
            assert(PDS.DistSharedCap.containsKey(distId), message: "No such distribution")
            let d <- PDS.DistSharedCap.remove(key: distId)!
            d.mintPackNFT(distId: distId, commitHashes: commitHashes, issuer: issuer, recvCap: recvCap)
            PDS.DistSharedCap[distId] <-! d
        }

        access(Operate) fun revealPackNFT(distId: UInt64, packId: UInt64, nftContractAddrs: [Address], nftContractNames: [String], nftIds: [UInt64], salt: String) {
            assert(PDS.DistSharedCap.containsKey(distId), message: "No such distribution")
            assert(
                nftContractAddrs.length == nftContractNames.length &&
                nftContractNames.length == nftIds.length,
                message: "NFTs must be fully described"
            )
            let d <- PDS.DistSharedCap.remove(key: distId)!
            let arr: [{IPackNFT.Collectible}] = []
            var i = 0
            while i < nftContractAddrs.length {
                let s = Collectible(address: nftContractAddrs[i], contractName: nftContractNames[i], id: nftIds[i])
                arr.append(s)
                i = i + 1
            }
            d.revealPackNFT(packId: packId, nfts: arr, salt: salt)
            PDS.DistSharedCap[distId] <-! d
        }

        access(Operate) fun openPackNFT(
            distId: UInt64,
            packId: UInt64,
            nftContractAddrs: [Address],
            nftContractNames: [String],
            nftIds: [UInt64],
            recvCap: &{NonFungibleToken.CollectionPublic},
            collectionStoragePath: StoragePath?
        ) {
            assert(PDS.DistSharedCap.containsKey(distId), message: "No such distribution")
            let d <- PDS.DistSharedCap.remove(key: distId)!
            let arr: [{IPackNFT.Collectible}] = []
            var i = 0
            while i < nftContractAddrs.length {
                let s = Collectible(address: nftContractAddrs[i], contractName: nftContractNames[i], id: nftIds[i])
                arr.append(s)
                i = i + 1
            }
            d.openPackNFT(packId: packId, nfts: arr, recvCap: recvCap, collectionStoragePath: collectionStoragePath)
            PDS.DistSharedCap[distId] <-! d
        }

    }

    /// Returns the manager collection capability to receive NFTs to be escrowed.
    ///
    access(contract) view fun getManagerCollectionCap(escrowCollectionPublic: PublicPath): Capability<&{NonFungibleToken.CollectionPublic}> {
        let pdsCollection = self.account.capabilities.get<&{NonFungibleToken.CollectionPublic}>(escrowCollectionPublic)!
        assert(pdsCollection.check(), message: "Please ensure PDS has created and linked a Collection for recieving escrows")
        return pdsCollection
    }



    /// Create a PackIssuer resource and return it to the caller.
    access(all) fun createPackIssuer(): @PackIssuer{
        return <- create PackIssuer()
    }

    /// Create a SharedCapabilities resource and return it to the caller.
    ///
    access(all) fun createSharedCapabilities(
        withdrawCap: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>,
        operatorCap: Capability<auth(IPackNFT.Operate) &{IPackNFT.IOperator}>
    ): @SharedCapabilities {
        return <- create SharedCapabilities(
            withdrawCap: withdrawCap,
            operatorCap: operatorCap
        )
    }

    /// Returns the details of a distribution if it exists, nil otherwise.
    ///
    access(all) view fun getDistInfo(distId: UInt64): DistInfo? {
        return PDS.Distributions[distId]
    }

    /// PDS contract initializer.
    ///
    init(
        PackIssuerStoragePath: StoragePath,
        PackIssuerCapRecv: PublicPath,
        DistCreatorStoragePath: StoragePath,
        DistManagerStoragePath: StoragePath,
        version: String
    ) {
        self.nextDistId = 1
        self.DistSharedCap <- {}
        self.Distributions = {}
        self.PackIssuerStoragePath = PackIssuerStoragePath
        self.PackIssuerCapRecv = PackIssuerCapRecv
        self.DistCreatorStoragePath = DistCreatorStoragePath
        self.DistManagerStoragePath = DistManagerStoragePath
        self.version = version

        // Create a DistributionCreator resource to share create capability with PackIssuer.
        self.account.storage.save(<- create DistributionCreator(), to: self.DistCreatorStoragePath)

        // Create a DistributionManager resource to manager distributions (withdraw for escrow, mint PackNFT todo: reveal / transfer).
        self.account.storage.save(<- create DistributionManager(), to: self.DistManagerStoragePath)
    }
}
