import Crypto
import NonFungibleToken from "NonFungibleToken"
import IPackNFT from "IPackNFT"
import MetadataViews from "MetadataViews"

/// Contract that defines Pack NFTs.
///
access(all) contract PackNFT: NonFungibleToken, IPackNFT {

    access(all) var totalSupply: UInt64
    access(all) let version: String
    access(all) let CollectionStoragePath: StoragePath
    access(all) let CollectionPublicPath: PublicPath
    access(all) let CollectionIPackNFTPublicPath: PublicPath
    access(all) let OperatorStoragePath: StoragePath

    /// Dictionary that stores Pack resources in the contract state (i.e., Pack NFT representations to keep track of states).
    ///
    access(contract) let packs: @{UInt64: Pack}

    access(all) event RevealRequest(id: UInt64, openRequest: Bool)
    access(all) event OpenRequest(id: UInt64)
    access(all) event Revealed(id: UInt64, salt: [UInt8], nfts: String)
    access(all) event Opened(id: UInt64)
    access(all) event Minted(id: UInt64, hash: [UInt8], distId: UInt64)
    access(all) event ContractInitialized()

    // TODO: Consider removing 'Withdraw' and 'Deposit' events now that similar 'Withdrawn' and 'Deposited' events are emitted in NonFungibleToken contract interface
    access(all) event Withdraw(id: UInt64, from: Address?)
    access(all) event Deposit(id: UInt64, to: Address?)

    /// Enum that defines the status of a Pack resource.
    ///
    access(all) enum Status: UInt8 {
        access(all) case Sealed
        access(all) case Revealed
        access(all) case Opened
    }

    /// Resource that defines a Pack NFT Operator, responsible for:
    ///  - Minting Pack NFTs and the corresponding Pack resources that keep track of states,
    ///  - Revealing sealed Pack resources, and
    ///  - opening revealed Pack resources.
    ///
    access(all) resource PackNFTOperator: IPackNFT.IOperator {

        /// Mint a new Pack NFT resource and corresponding Pack resource; store the Pack resource in the contract's packs dictionary
        /// and return the Pack NFT resource to the caller.
        ///
        access(IPackNFT.Operate) fun mint(distId: UInt64, commitHash: String, issuer: Address): @{IPackNFT.NFT} {
            let nft <- create NFT(commitHash: commitHash, issuer: issuer)
            PackNFT.totalSupply = PackNFT.totalSupply + 1
            let p <- create Pack(commitHash: commitHash, issuer: issuer)
            PackNFT.packs[nft.id] <-! p
            emit Minted(id: nft.id, hash: commitHash.decodeHex(), distId: distId)
            return <- nft
        }

        /// Reveal a Sealed Pack resource.
        ///
        access(IPackNFT.Operate) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
            let p <- PackNFT.packs.remove(key: id) ?? panic("no such pack")
            p.reveal(id: id, nfts: nfts, salt: salt)
            PackNFT.packs[id] <-! p
        }

        /// Open a Revealed Pack NFT resource.
        ///
        access(IPackNFT.Operate) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}]) {
            let p <- PackNFT.packs.remove(key: id) ?? panic("no such pack")
            p.open(id: id, nfts: nfts)
            PackNFT.packs[id] <-! p
        }

        /// PackNFTOperator resource initializer.
        ///
        view init() {}
    }

    /// Resource that defines a Pack NFT.
    ///
    access(all) resource Pack {
        access(all) let hash: [UInt8]
        access(all) let issuer: Address
        access(all) var status: Status
        access(all) var salt: [UInt8]?

        access(all) view fun verify(nftString: String): Bool {
            assert(self.status != Status.Sealed, message: "Pack not revealed yet")
            var hashString = String.encodeHex(self.salt!)
            hashString = hashString.concat(",").concat(nftString)
            let hash = HashAlgorithm.SHA2_256.hash(hashString.utf8)
            assert(String.encodeHex(self.hash) == String.encodeHex(hash), message: "CommitHash was not verified")
            return true
        }

        access(self) fun _verify(nfts: [{IPackNFT.Collectible}], salt: String, commitHash: String): String {
            var hashString = salt
            var nftString = nfts[0].hashString()
            var i = 1
            while i < nfts.length {
                let s = nfts[i].hashString()
                nftString = nftString.concat(",").concat(s)
                i = i + 1
            }
            hashString = hashString.concat(",").concat(nftString)
            let hash = HashAlgorithm.SHA2_256.hash(hashString.utf8)
            assert(String.encodeHex(self.hash) == String.encodeHex(hash), message: "CommitHash was not verified")
            return nftString
        }

        access(contract) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
            assert(self.status == Status.Sealed, message: "Pack status is not Sealed")
            let v = self._verify(nfts: nfts, salt: salt, commitHash: String.encodeHex(self.hash))
            self.salt = salt.decodeHex()
            self.status = Status.Revealed
            emit Revealed(id: id, salt: salt.decodeHex(), nfts: v)
        }

        access(contract) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}]) {
            assert(self.status == Status.Revealed, message: "Pack status is not Revealed")
            self._verify(nfts: nfts, salt: String.encodeHex(self.salt!), commitHash: String.encodeHex(self.hash))
            self.status = Status.Opened
            emit Opened(id: id)
        }

        /// Pack resource initializer.
        ///
        view init(commitHash: String, issuer: Address) {
            // Set the hash and issuer from the arguments.
            self.hash = commitHash.decodeHex()
            self.issuer = issuer

            // Initial status is Sealed.
            self.status = Status.Sealed

            // Salt is nil until reveal.
            self.salt = nil
        }
    }

    /// Resource that defines a Pack NFT.
    ///
    access(all) resource NFT: NonFungibleToken.NFT, IPackNFT.NFT, IPackNFT.IPackNFTToken, IPackNFT.IPackNFTOwnerOperator {
        /// This NFT's unique ID.
        ///
        access(all) let id: UInt64

        /// This NFT's commit hash, used to verify the IDs of the NFTs in the Pack.
        ///
        access(all) let hash: [UInt8]

        /// This NFT's issuer.
        access(all) let issuer: Address

        /// Reveal a Sealed Pack resource.
        ///
        access(NonFungibleToken.Update) fun reveal(openRequest: Bool) {
            PackNFT.revealRequest(id: self.id, openRequest: openRequest)
        }

        /// Open a Revealed Pack resource.
        ///
        access(NonFungibleToken.Update) fun open() {
            PackNFT.openRequest(id: self.id)
        }

        /// Event emitted when a NFT is destroyed (replaces Burned event before Cadence 1.0 update)
        ///
        access(all) event ResourceDestroyed(id: UInt64 = self.id)

        /// Executed by calling the Burner contract's burn method (i.e., conforms to the Burnable interface)
        ///
        access(contract) fun burnCallback() {
            PackNFT.totalSupply = PackNFT.totalSupply - 1
            destroy <- PackNFT.packs.remove(key: self.id) ?? panic("no such pack")
        }

        /// NFT resource initializer.
        ///
        view init(commitHash: String, issuer: Address) {
            self.id = self.uuid
            self.hash = commitHash.decodeHex()
            self.issuer = issuer
        }

        /// Create an empty Collection for Pinnacle NFTs and return it to the caller
        ///
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- PackNFT.createEmptyCollection(nftType: Type<@NFT>())
        }

        /// Return the metadata view types available for this NFT.
        ///
        access(all) view fun getViews(): [Type] {
            return []
        }

        /// Resolve this NFT's metadata views.
        //// TODO: Implement metadata views as needed for the NFT the PackNFT is intended to contain
        ///
        access(all) view fun resolveView(_ view: Type): AnyStruct? {
            return nil
        }
    }

    /// Resource that defines a Collection of Pack NFTs.
    ///
    access(all) resource Collection: NonFungibleToken.Collection, IPackNFT.IPackNFTCollectionPublic {
        /// Dictionary of NFT conforming tokens.
        /// NFT is a resource type with a UInt64 ID field.
        ///
        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}

        /// Collection resource initializer.
        ///
        view init() {
            self.ownedNFTs <- {}
        }

        /// Remove an NFT from the collection and moves it to the caller.
        ///
        access(NonFungibleToken.Withdraw) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            // Withdrawn event emitted from NonFungibleToken contract interface.
            emit Withdraw(id: token.id, from: self.owner?.address) // TODO: Consider removing
            return <- token
        }

        /// Deposit an NFT into this Collection.
        ///
        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @NFT
            let id: UInt64 = token.id
            // Add the new token to the dictionary which removes the old one.
            let oldToken <- self.ownedNFTs[id] <- token

            // Deposited event emitted from NonFungibleToken contract interface.
            emit Deposit(id: id, to: self.owner?.address) // TODO: Consider removing
            destroy oldToken
        }

        /// Return an array of the IDs that are in the collection.
        ///
        access(all) view fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        /// Return the amount of NFTs stored in the collection.
        ///
        access(all) view fun getLength(): Int {
            return self.ownedNFTs.length
        }

        /// Return a list of NFT types that this receiver accepts.
        ///
        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@NFT>()] = true
            return supportedTypes
        }

        /// Return whether or not the given type is accepted by the collection.
        ///
        access(all) view fun isSupportedNFTType(type: Type): Bool {
            if type == Type<@NFT>() {
                return true
            }
            return false
        }

        /// Return a reference to an NFT in the Collection.
        ///
        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        /// Create an empty Collection of the same type and returns it to the caller.
        ///
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-PackNFT.createEmptyCollection(nftType: Type<@NFT>())
        }
    }

    /// Emit a RevealRequest event to signal a Sealed Pack NFT should be revealed.
    ///
    access(contract) fun revealRequest(id: UInt64, openRequest: Bool ) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        assert(p.status.rawValue == Status.Sealed.rawValue, message: "Pack status must be Sealed for reveal request")
        emit RevealRequest(id: id, openRequest: openRequest)
    }

    /// Emit an OpenRequest event to signal a Revealed Pack NFT should be opened.
    ///
    access(contract) fun openRequest(id: UInt64) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        assert(p.status.rawValue == Status.Revealed.rawValue, message: "Pack status must be Revealed for open request")
        emit OpenRequest(id: id)
    }

    /// Reveal a Sealed Pack NFT.
    ///
    access(all) fun publicReveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        p.reveal(id: id, nfts: nfts, salt: salt)
    }

    /// Return a reference to a Pack resource stored in the contract state.
    ///
    access(all) view fun borrowPackRepresentation(id: UInt64): &Pack? {
        return (&self.packs[id] as &Pack?)!
    }

    /// Create an empty Collection for Pack NFTs and return it to the caller.
    ///
    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        if nftType != Type<@NFT>() {
            panic("NFT type is not supported")
        }
        return <- create Collection()
    }

    /// Return the metadata views implemented by this contract.
    ///
    /// @return An array of Types defining the implemented views. This value will be used by
    ///         developers to know which parameter to pass to the resolveView() method.
    ///
    access(all) view fun getContractViews(resourceType: Type?): [Type] {
        return [
            Type<MetadataViews.NFTCollectionData>(),
            Type<MetadataViews.NFTCollectionDisplay>()
        ]
    }

    /// Resolve a metadata view for this contract.
    ///
    /// @param view: The Type of the desired view.
    /// @return A structure representing the requested view.
    ///
    access(all) view fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                let collectionData = MetadataViews.NFTCollectionData(
                    storagePath: /storage/exampleNFTCollection,
                    publicPath: /public/exampleNFTCollection,
                    publicCollection: Type<&Collection>(),
                    publicLinkedType: Type<&Collection>(),
                    createEmptyCollectionFunction: (fun(): @{NonFungibleToken.Collection} {
                        return <-PackNFT.createEmptyCollection(nftType: Type<@NFT>())
                    })
                )
                return collectionData
            case Type<MetadataViews.NFTCollectionDisplay>():
                let media = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://assets.website-files.com/5f6294c0c7a8cdd643b1c820/5f6294c0c7a8cda55cb1c936_Flow_Wordmark.svg"
                    ),
                    mediaType: "image/svg+xml"
                )
                return MetadataViews.NFTCollectionDisplay(
                    name: "The Example Collection",
                    description: "This collection is used as an example to help you develop your next Flow NFT.",
                    externalURL: MetadataViews.ExternalURL("https://example-nft.onflow.org"),
                    squareImage: media,
                    bannerImage: media,
                    socials: {
                        "twitter": MetadataViews.ExternalURL("https://twitter.com/flow_blockchain")
                    }
                )
        }
        return nil
    }

    /// PackNFT contract initializer.
    ///
    init(
        CollectionStoragePath: StoragePath,
        CollectionPublicPath: PublicPath,
        CollectionIPackNFTPublicPath: PublicPath,
        OperatorStoragePath: StoragePath,
        version: String
    ) {
        self.totalSupply = 0
        self.packs <- {}
        self.CollectionStoragePath = CollectionStoragePath
        self.CollectionPublicPath = CollectionPublicPath
        self.CollectionIPackNFTPublicPath = CollectionIPackNFTPublicPath
        self.OperatorStoragePath = OperatorStoragePath
        self.version = version

        // Create a collection to receive Pack NFTs and publish public receiver capabilities.
        self.account.storage.save(<- create Collection(), to: self.CollectionStoragePath)
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{NonFungibleToken.CollectionPublic}>(self.CollectionStoragePath),
            at: self.CollectionPublicPath
        )
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{IPackNFT.IPackNFTCollectionPublic}>(self.CollectionStoragePath),
            at: self.CollectionIPackNFTPublicPath
        )

        // Create a Pack NFT operator to share mint capability with proxy.
        self.account.storage.save(<- create PackNFTOperator(), to: self.OperatorStoragePath)
        self.account.capabilities.storage.issue<&{IPackNFT.IOperator}>(self.OperatorStoragePath)
    }
}
