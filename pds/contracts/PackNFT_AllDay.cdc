import Crypto
import NonFungibleToken from "NonFungibleToken"
import FungibleToken from "FungibleToken"
import IPackNFT from "IPackNFT"
import MetadataViews from "MetadataViews"

access(all) contract PackNFT: NonFungibleToken, IPackNFT {

    access(all) var totalSupply: UInt64
    access(all) let version: String
    access(all) let CollectionStoragePath: StoragePath
    access(all) let CollectionPublicPath: PublicPath
    access(all) let CollectionIPackNFTPublicPath: PublicPath
    access(all) let OperatorStoragePath: StoragePath

    // representation of the NFT in this contract to keep track of states
    access(contract) let packs: @{UInt64: Pack}

    access(all) event RevealRequest(id: UInt64, openRequest: Bool)
    access(all) event OpenRequest(id: UInt64)
    access(all) event Revealed(id: UInt64, salt: String, nfts: String)
    access(all) event Opened(id: UInt64)
    access(all) event Mint(id: UInt64, commitHash: String, distId: UInt64)
    access(all) event ContractInitialized()
    access(all) event Withdraw(id: UInt64, from: Address?)
    access(all) event Deposit(id: UInt64, to: Address?)
    access(all) event Burned(id: UInt64)

    access(all) enum Status: UInt8 {
        access(all) case Sealed
        access(all) case Revealed
        access(all) case Opened
    }

    access(all) resource PackNFTOperator: IPackNFT.IOperator {

        access(IPackNFT.Operatable) fun mint(distId: UInt64, commitHash: String, issuer: Address): @{IPackNFT.NFT}{
            let nft <- create NFT(commitHash: commitHash, issuer: issuer)
            PackNFT.totalSupply = PackNFT.totalSupply + 1
            let p  <-create Pack(commitHash: commitHash, issuer: issuer)
            PackNFT.packs[nft.id] <-! p
            emit Mint(id: nft.id, commitHash: commitHash, distId: distId)
            return <- nft
        }

        access(IPackNFT.Operatable) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
            let p <- PackNFT.packs.remove(key: id) ?? panic("no such pack")
            p.reveal(id: id, nfts: nfts, salt: salt)
            PackNFT.packs[id] <-! p
        }

        access(IPackNFT.Operatable) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}]) {
            let p <- PackNFT.packs.remove(key: id) ?? panic("no such pack")
            p.open(id: id, nfts: nfts)
            PackNFT.packs[id] <-! p
        }

        init(){}
    }

    access(all) resource Pack {
        access(all) let commitHash: String
        access(all) let issuer: Address
        access(all) var status: PackNFT.Status
        access(all) var salt: String?

        access(all) fun verify(nftString: String): Bool {
            assert(self.status != PackNFT.Status.Sealed, message: "Pack not revealed yet")
            var hashString = self.salt!
            hashString = hashString.concat(",").concat(nftString)
            let hash = HashAlgorithm.SHA2_256.hash(hashString.utf8)
            assert(self.commitHash == String.encodeHex(hash), message: "CommitHash was not verified")
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
            assert(self.commitHash == String.encodeHex(hash), message: "CommitHash was not verified")
            return nftString
        }

        access(contract) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
            assert(self.status == PackNFT.Status.Sealed, message: "Pack status is not Sealed")
            let v = self._verify(nfts: nfts, salt: salt, commitHash: self.commitHash)
            self.salt = salt
            self.status = PackNFT.Status.Revealed
            emit Revealed(id: id, salt: salt, nfts: v)
        }

        access(contract) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}]) {
            assert(self.status == PackNFT.Status.Revealed, message: "Pack status is not Revealed")
            self._verify(nfts: nfts, salt: self.salt!, commitHash: self.commitHash)
            self.status = PackNFT.Status.Opened
            emit Opened(id: id)
        }

        init(commitHash: String, issuer: Address) {
            self.commitHash = commitHash
            self.issuer = issuer
            self.status = PackNFT.Status.Sealed
            self.salt = nil
        }
    }

    access(all) resource NFT: NonFungibleToken.INFT, IPackNFT.IPackNFTToken, IPackNFT.IPackNFTOwnerOperator, NonFungibleToken.NFT, IPackNFT.NFT {
        access(all) let id: UInt64
        access(all) let commitHash: String
        access(all) let issuer: Address

        access(NonFungibleToken.Update | NonFungibleToken.Owner) fun reveal(openRequest: Bool){
            PackNFT.revealRequest(id: self.id, openRequest: openRequest)
        }

        access(NonFungibleToken.Update | NonFungibleToken.Owner) fun open(){
            PackNFT.openRequest(id: self.id)
        }

        /// Event emitted when a NFT is destroyed (replaces Burned event before Cadence 1.0 update)
        ///
        access(all) event ResourceDestroyed(id: UInt64 = self.id)

        // When destroying a NFT resource, the corresponding PackNFT resource now has to be detroyed as well in a separate call.
        // This is because custome destroy() function is not allowed in Cadence 1.0.
        // destroy() {
        //     let p <- PackNFT.packs.remove(key: self.id) ?? panic("no such pack")
        //     PackNFT.totalSupply = PackNFT.totalSupply - (1 as UInt64)

        //     emit Burned(id: self.id)
        //     destroy p
        // }

        init(commitHash: String, issuer: Address) {
            self.id = self.uuid
            self.commitHash = commitHash
            self.issuer = issuer
        }

        /// Create an empty Collection for Pinnacle NFTs and return it to the caller
        ///
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- PackNFT.createEmptyCollection(nftType: Type<@PackNFT.NFT>())
        }

        // All supported metadata views for the Moment including the Core NFT Views
        //
        view access(all) fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.Medias>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Royalties>(),
                Type<MetadataViews.Serial>()
            ]
        }

        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: "NFL All Day Pack",
                        description: "Reveals official NFL All Day Moments when opened",
                        thumbnail: MetadataViews.HTTPFile(url: self.getImage(imageType: "image", format: "jpeg", width: 256))
                    )
                case Type<MetadataViews.ExternalURL>():
                    return MetadataViews.ExternalURL("https://nflallday.com/packnfts/".concat(self.id.toString())) // might have to make a URL that redirects to packs page based on packNFT id -> distribution id
                case Type<MetadataViews.Medias>():
                    return MetadataViews.Medias(
                        [
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getImage(imageType: "image", format: "jpeg", width: 512)),
                                mediaType: "image/jpeg"
                            )
                        ]
                    )
                case Type<MetadataViews.NFTCollectionData>():
                    return MetadataViews.NFTCollectionData(
                        storagePath: PackNFT.CollectionStoragePath,
                        publicPath: PackNFT.CollectionPublicPath,
                        publicCollection: Type<&PackNFT.Collection>(),
                        publicLinkedType: Type<&PackNFT.Collection>(),
                        createEmptyCollectionFunction: (fun (): @{NonFungibleToken.Collection} {
                            return <-PackNFT.createEmptyCollection(nftType: Type<@PackNFT.NFT>())
                        })
                    )
                case Type<MetadataViews.NFTCollectionDisplay>():
                    let bannerImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://assets.nflallday.com/flow/catalogue/NFLAD_BANNER.png"
                        ),
                        mediaType: "image/png"
                    )
                    let squareImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://assets.nflallday.com/flow/catalogue/NFLAD_SQUARE.png"
                        ),
                        mediaType: "image/png"
                    )
                    return MetadataViews.NFTCollectionDisplay(
                        name: "NFL All Day Packs",
                        description: "Officially Licensed Digital Collectibles Featuring the NFL’s Best Highlights. Buy, Sell and Collect Your Favorite NFL Moments",
                        externalURL: MetadataViews.ExternalURL("https://nflallday.com/"),
                        squareImage: squareImage,
                        bannerImage: bannerImage,
                        socials: {
                            "instagram": MetadataViews.ExternalURL("https://www.instagram.com/nflallday/"),
                            "twitter": MetadataViews.ExternalURL("https://twitter.com/NFLAllDay"),
                            "discord": MetadataViews.ExternalURL("https://discord.com/invite/5K6qyTzj2k")
                        }
                    )
                case Type<MetadataViews.Royalties>():
                    let royaltyReceiver: Capability<&{FungibleToken.Receiver}> =
                        getAccount({{.RoyaltyAddress}}).capabilities.get<&{FungibleToken.Receiver}>(MetadataViews.getRoyaltyReceiverPublicPath())!
                    return MetadataViews.Royalties(
                        [
                            MetadataViews.Royalty(
                                receiver: royaltyReceiver,
                                cut: 0.05,
                                description: "NFL All Day marketplace royalty"
                            )
                        ]
                    )
                case Type<MetadataViews.Serial>():
                    return MetadataViews.Serial(self.id)
            }
            return nil
        }

        access(all) fun assetPath(): String {
            return "https://media.nflallday.com/packnfts/".concat(self.id.toString()).concat("/media/")
        }

        access(all) fun getImage(imageType: String, format: String, width: Int): String {
            return self.assetPath().concat(imageType).concat("?format=").concat(format).concat("&width=").concat(width.toString())
        }
    }

    access(all) resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        IPackNFT.IPackNFTCollectionPublic,
        // MetadataViews.ResolverCollection
        NonFungibleToken.Collection {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}

        init () {
            self.ownedNFTs <- {}
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        access(NonFungibleToken.Withdraw | NonFungibleToken.Owner) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")
            emit Withdraw(id: token.id, from: self.owner?.address)
            return <- token
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @PackNFT.NFT

            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token
            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        // getIDs returns an array of the IDs that are in the collection
        view access(all) fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        /// Return the amount of NFTs stored in the collection
        ///
        view access(all) fun getLength(): Int {
            return self.ownedNFTs.keys.length
        }

        /// Return a list of NFT types that this receiver accepts
        ///
        view access(all) fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@PackNFT.NFT>()] = true
            return supportedTypes
        }

        /// Return whether or not the given type is accepted by the collection
        /// A collection that can accept any type should just return true by default
        ///
        view access(all) fun isSupportedNFTType(type: Type): Bool {
            if type == Type<@PackNFT.NFT>() {
                return true
            }
            return false
        }

        // borrowNFT gets a reference to an NFT in the collection
        // so that the caller can read its metadata and call its methods
        view access(all) fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        view access(all) fun borrowPackNFT(id: UInt64): &{IPackNFT.NFT}? {
            return self.borrowNFT(id) as! &{IPackNFT.NFT}?
        }

         /// createEmptyCollection creates an empty Collection of the same type
        /// and returns it to the caller
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <-PackNFT.createEmptyCollection(nftType: Type<@PackNFT.NFT>())
        }
    }

    access(contract) fun revealRequest(id: UInt64, openRequest: Bool ) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        assert(p.status.rawValue == PackNFT.Status.Sealed.rawValue, message: "Pack status must be Sealed for reveal request")
        emit RevealRequest(id: id, openRequest: openRequest)
    }

    access(contract) fun openRequest(id: UInt64) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        assert(p.status.rawValue == PackNFT.Status.Revealed.rawValue, message: "Pack status must be Revealed for open request")
        emit OpenRequest(id: id)
    }

    access(all) fun publicReveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String) {
        let p = PackNFT.borrowPackRepresentation(id: id) ?? panic ("No such pack")
        p.reveal(id: id, nfts: nfts, salt: salt)
    }

    access(all) fun borrowPackRepresentation(id: UInt64):  &Pack? {
        return (&self.packs[id] as &Pack?)!
    }

    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        if nftType != Type<@PackNFT.NFT>() {
            panic("NFT type is not supported")
        }
        return <- create Collection()
    }

    /// Function that returns all the Metadata Views implemented by a Non Fungible Token
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

    /// Function that resolves a metadata view for this contract.
    ///
    /// @param view: The Type of the desired view.
    /// @return A structure representing the requested view.
    ///
    access(all) fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                let collectionData = MetadataViews.NFTCollectionData(
                    storagePath: /storage/cadenceExampleNFTCollection,
                    publicPath: /public/cadenceExampleNFTCollection,
                    publicCollection: Type<&PackNFT.Collection>(),
                    publicLinkedType: Type<&PackNFT.Collection>(),
                    createEmptyCollectionFunction: (fun(): @{NonFungibleToken.Collection} {
                        return <-PackNFT.createEmptyCollection(nftType: Type<@PackNFT.NFT>())
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

    init(
        CollectionStoragePath: StoragePath,
        CollectionPublicPath: PublicPath,
        CollectionIPackNFTPublicPath: PublicPath,
        OperatorStoragePath: StoragePath,
        version: String
    ){
        self.totalSupply = 0
        self.packs <- {}
        self.CollectionStoragePath = CollectionStoragePath
        self.CollectionPublicPath = CollectionPublicPath
        self.CollectionIPackNFTPublicPath = CollectionIPackNFTPublicPath
        self.OperatorStoragePath = OperatorStoragePath
        self.version = version

        // Create a collection to receive Pack NFTs
        let collection <- create Collection()
        self.account.storage.save(<-collection, to: self.CollectionStoragePath)
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{NonFungibleToken.CollectionPublic}>(self.CollectionStoragePath),
            at: self.CollectionPublicPath
        )
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{IPackNFT.IPackNFTCollectionPublic}>(self.CollectionStoragePath),
            at: self.CollectionIPackNFTPublicPath
        )

        // Create a operator to share mint capability with proxy
        let operator <- create PackNFTOperator()
        self.account.storage.save(<-operator, to: self.OperatorStoragePath)
        self.account.capabilities.storage.issue<&{IPackNFT.IOperator}>(self.OperatorStoragePath)
    }

}