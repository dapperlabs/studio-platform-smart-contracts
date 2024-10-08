/*
    Author: Jude Zhu jude.zhu@dapperlabs.com
*/


// import FungibleToken from "FungibleToken"
import NonFungibleToken from "NonFungibleToken"
import MetadataViews from "MetadataViews"
import ViewResolver from "ViewResolver"

/*
    There are 2 levels of entity:
    1. Edition
    2. NFT
    
    An Edition is created with metadata. NFTs are minted out of Editions.
 */

// The EditionNFT contract
//
access(all) contract EditionNFT: NonFungibleToken {
    //------------------------------------------------------------
    // Events
    //------------------------------------------------------------

    // Contract Events
    //
    access(all) event ContractInitialized()

    // NFT Collection Events
    //
    access(all) event Withdraw(id: UInt64, from: Address?)
    access(all) event Deposit(id: UInt64, to: Address?)
    access(all) event Burned(id: UInt64)
    access(all) event Minted(id: UInt64, editionID: UInt64)

    // Edition Events
    //
    // Emitted when a new edition has been created by an admin
    access(all) event EditionCreated(
        id: UInt64, 
        metadata: {String: String}
    )
    // Emitted when an edition is either closed by an admin, or the max amount of moments have been minted
    access(all) event EditionClosed(id: UInt64)

    //------------------------------------------------------------
    // Named values
    //------------------------------------------------------------

    // Named Paths
    //
    access(all) let CollectionStoragePath:  StoragePath
    access(all) let CollectionPublicPath:   PublicPath
    access(all) let AdminStoragePath:       StoragePath

    //------------------------------------------------------------
    // Publicly readable contract state
    //------------------------------------------------------------

    // totalSupply
    // The total number of {{ contractName }} NFTs that have been minted.
    //
    access(all) var totalSupply:        UInt64

    // totalEditions
    // The total number of {{ contractName }} editions that have been created.
    //
    access(all) var totalEditions: UInt64

    access(all) var nextEditionID:      UInt64


    //------------------------------------------------------------
    // Internal contract state
    //------------------------------------------------------------

    // Metadata Dictionaries
    //
    // This is so we can find Series by their names (via seriesByID)
    access(self) let editionByID:       @{UInt64: Edition}

    //------------------------------------------------------------
    // Edition
    //------------------------------------------------------------

    // A public struct to access Edition data
    //
    access(all) struct EditionData {
        access(all) let id: UInt64
        access(all) var numMinted: UInt64
        access(all) var active: Bool
        access(all) let metadata: {String: String}

        // initializer
        //
        view init (id: UInt64) {
            let edition = &EditionNFT.editionByID[id] as &EditionNFT.Edition?
                ?? panic("edition does not exist")
            self.id = id
            self.metadata = edition.getMetadata()
            self.numMinted = edition.numMinted
            self.active = edition.active
        }
    }

    // A top level Edition with a unique ID
    //
    access(all) resource Edition {
        access(all) let id: UInt64
        // Contents writable if borrowed!
        // This is deliberate, as it allows admins to update the data.
        access(all) var numMinted: UInt64
        access(all) let metadata: {String: String}
        access(all) var active: Bool

        // Close this edition
        //
        access(all) fun close() {
            pre {
                self.active == true: "edtion is already closed"
            }

            self.active = false
            emit EditionClosed(id: self.id)
        }

        /// returns the metadata set for this edition
        access(all) view fun getMetadata(): {String:String} {
            return self.metadata
        }

        // Mint a Seasonal NFT in this edition, with the given minting mintingDate.
        // Note that this will panic if the max mint size has already been reached.
        //
        access(contract) fun mint(): @EditionNFT.NFT {
            pre {
                self.active: "edition closed, cannot mint"
            }

            // Create thek NFT, filled out with our information
            let nft <- create NFT(
                id: EditionNFT.totalSupply + 1,
                editionID: self.id,
                serialNumber: self.numMinted + 1
            )
            EditionNFT.totalSupply = EditionNFT.totalSupply + 1
            // Keep a running total (you'll notice we used this as the serial number)
            self.numMinted = self.numMinted + 1 as UInt64

            return <- nft
        }

        // initializer
        //
        init (metadata: {String: String}) {
            self.id = EditionNFT.nextEditionID
            self.metadata = metadata
            self.numMinted = 0 as UInt64
            self.active = true

            EditionNFT.nextEditionID = self.id + 1 as UInt64
            emit EditionCreated(id: self.id, metadata: self.metadata)
        }
    }

    // Get the publicly available data for a Edition
    //
    access(all) view fun getEditionData(id: UInt64): EditionNFT.EditionData {
        pre {
            EditionNFT.editionByID[id] != nil: "Cannot borrow edition, no such id"
        }

        return EditionNFT.EditionData(id: id)
    }

    //------------------------------------------------------------
    // NFT
    //------------------------------------------------------------

    // A Seasonal NFT
    //
    access(all) resource NFT: NonFungibleToken.NFT {
        access(all) let id: UInt64
        access(all) let editionID: UInt64

        access(all) event ResourceDestroyed(
            id: UInt64 = self.id,
            editionID: UInt64 = self.editionID,
            serialNumber: UInt64 = 0,
            mintingDate: UFix64 = 0.0
        )

        // NFT initializer
        //
        init(
            id: UInt64,
            editionID: UInt64,
            serialNumber: UInt64
        ) {
            pre {
                EditionNFT.editionByID[editionID] != nil: "no such editionID"
                EditionData(id: editionID).active == true: "edition already closed"
            }

            self.id = id
            self.editionID = editionID

            emit Minted(id: self.id, editionID: self.editionID)
        }
        access(all) fun getTraits() : {String: AnyStruct} {
            let edition: EditionData = EditionNFT.getEditionData(id: self.editionID)

            let traitDictionary: {String: AnyStruct} = {}

            for name in edition.metadata.keys {
                let value = edition.metadata[name] ?? ""
                if value != "" {
                    traitDictionary.insert(key: name, value)
                }
            }
            return traitDictionary
        }

        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- EditionNFT.createEmptyCollection(nftType: Type<@EditionNFT.NFT>())
        }

        /// get the metadata view types available for this nft
        ///
        access(all) view fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.Editions>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.Traits>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Royalties>()
            ]
        }

        /// resolve a metadata view type returning the properties of the view type
        ///
        access(all) fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: "Edition NFT",
                        description: "An edition NFT",
                        thumbnail: MetadataViews.HTTPFile(url:"https://edition.com")
                    )

                case Type<MetadataViews.Editions>():
                let editionData = EditionNFT.getEditionData(id: self.editionID)
                    let editionInfo = MetadataViews.Edition(
                        name: nil,
                        number: 0,
                        max: editionData.numMinted
                    )
                    let editionList: [MetadataViews.Edition] = [editionInfo]
                    return MetadataViews.Editions(
                        editionList
                    )

                
                case Type<MetadataViews.NFTCollectionData>():
                    return EditionNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>())

                case Type<MetadataViews.Traits>():
                    return MetadataViews.dictToTraits(dict: self.getTraits(), excludedNames: nil)
                case Type<MetadataViews.NFTCollectionDisplay>():
                    return EditionNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionDisplay>())

                case Type<MetadataViews.Royalties>():
                    return EditionNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.Royalties>())
            }

            return nil
        }
    }

    //------------------------------------------------------------
    // Collection
    //------------------------------------------------------------

    // Deprecated: This is no longer used for defining access control anymore.
    //
    access(all) resource interface EditionNFTCollectionPublic {}

    // An NFT Collection
    //
    access(all) resource Collection:
        NonFungibleToken.Collection,
        EditionNFTCollectionPublic
    {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an UInt64 ID field
        //
        access(all) var ownedNFTs: @{UInt64: {NonFungibleToken.NFT}}

        // Return a list of NFT types that this receiver accepts
        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@EditionNFT.NFT>()] = true
            return supportedTypes
        }

        // Return whether or not the given type is accepted by the collection
        // A collection that can accept any type should just return true by default
        access(all) view fun isSupportedNFTType(type: Type): Bool {
            if type == Type<@EditionNFT.NFT>() {
                return true
            }
            return false
        }

        // Return the amount of NFTs stored in the collection
        access(all) view fun getLength(): Int {
            return self.ownedNFTs.keys.length
        }

        // Create an empty Collection for Golazos NFTs and return it to the caller
        access(all) fun createEmptyCollection(): @{NonFungibleToken.Collection} {
            return <- EditionNFT.createEmptyCollection(nftType: Type<@EditionNFT.NFT>())
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        //
        access(NonFungibleToken.Withdraw) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        //
        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {
            let token <- token as! @EditionNFT.NFT
            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        // batchDeposit takes a Collection object as an argument
        // and deposits each contained NFT into this Collection
        //
        access(all) fun batchDeposit(tokens: @{NonFungibleToken.Collection}) {
            // Get an array of the IDs to be deposited
            let keys = tokens.getIDs()

            // Iterate through the keys in the collection and deposit each one
            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

            // Destroy the empty Collection
            destroy tokens
        }

        // getIDs returns an array of the IDs that are in the collection
        //
        access(all) view fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // borrowNFT gets a reference to an NFT in the collection
        //
        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {
            return &self.ownedNFTs[id]
        }

        // borrowEditionNFT gets a reference to an EditionNFT in the collection
        //
        access(all) fun borrowEditionNFT(id: UInt64): &EditionNFT.NFT? {
            return self.borrowNFT(id) as! &EditionNFT.NFT?

        }

        // Collection initializer
        //
        init() {
            self.ownedNFTs <- {}
        }
    }

    // public function that anyone can call to create a new empty collection
    //
    access(all) fun createEmptyCollection(nftType: Type): @{NonFungibleToken.Collection} {
        if nftType != Type<@EditionNFT.NFT>() {
            panic("NFT type is not supported")
        }
        return <- create Collection()
    }

    //------------------------------------------------------------
    // Admin
    //------------------------------------------------------------

    /// Entitlement that grants the ability to mint Golazos NFTs
    access(all) entitlement Mint

    /// Entitlement that grants the ability to operate admin functions
    access(all) entitlement Operate

    // An interface containing the Admin function that allows minting NFTs
    //
    // This is no longer used for defining access control anymore.
    // Keeping this because removing it is not a valid change for contract update
    access(all) resource interface NFTMinter {}

    // A resource that allows managing metadata and minting NFTs
    //
    access(all) resource Admin: NFTMinter {

        // Borrow an Edition
        //
        access(self) view fun borrowEdition(id: UInt64): &EditionNFT.Edition {
            pre {
                EditionNFT.editionByID[id] != nil: "Cannot borrow edition, no such id"
            }

            return (&EditionNFT.editionByID[id])!
        }

        // Create a Edition 
        //
        access(Operate) fun createEdition(metadata: {String: String}): UInt64 {
            // Create and store the new edition
            let edition <- create EditionNFT.Edition(
                metadata: metadata,
            )
            let editionID = edition.id
            EditionNFT.editionByID[edition.id] <-! edition

            // Return the new ID for convenience
            return editionID
        }


        // Close an Edition
        //
        access(Operate) fun closeEdition(id: UInt64): UInt64 {
            if let edition = &EditionNFT.editionByID[id] as &EditionNFT.Edition? {
                edition.close()
                return edition.id
            }
            panic("edition does not exist")
        }

        // Mint a single NFT
        // The Edition for the given ID must already exist
        //
        access(Mint) fun mintNFT(editionID: UInt64): @EditionNFT.NFT {
            pre {
                // Make sure the edition we are creating this NFT in exists
                EditionNFT.editionByID.containsKey(editionID): "No such EditionID"
            }
            return <- self.borrowEdition(id: editionID).mint()
        }
    }

    /// Return the metadata view types available for this contract
    ///
    access(all) view fun getContractViews(resourceType: Type?): [Type] {
        return [Type<MetadataViews.NFTCollectionData>(), Type<MetadataViews.NFTCollectionDisplay>(), Type<MetadataViews.Royalties>()]
    }

    /// Resolve this contract's metadata views
    ///
    access(all) view fun resolveContractView(resourceType: Type?, viewType: Type): AnyStruct? {
        post {
            result == nil || result!.getType() == viewType: "The returned view must be of the given type or nil"
        }
        switch viewType {
            case Type<MetadataViews.NFTCollectionData>():
                return MetadataViews.NFTCollectionData(
                    storagePath: self.CollectionStoragePath,
                    publicPath: self.CollectionPublicPath,
                    publicCollection: Type<&EditionNFT.Collection>(),
                    publicLinkedType: Type<&EditionNFT.Collection>(),
                    createEmptyCollectionFunction: (fun (): @{NonFungibleToken.Collection} {
                        return <-EditionNFT.createEmptyCollection(nftType: Type<@EditionNFT.NFT>())
                    })
                )
            case Type<MetadataViews.NFTCollectionDisplay>():
                let bannerImage = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://assets.edition.com"
                    ),
                    mediaType: "image/png"
                )
                let squareImage = MetadataViews.Media(
                    file: MetadataViews.HTTPFile(
                        url: "https://assets.edition.com"
                    ),
                    mediaType: "image/png"
                )
                return MetadataViews.NFTCollectionDisplay(
                    name: "Edition NFT",
                    description: "An edition nft",
                    externalURL: MetadataViews.ExternalURL("https://edition.com/"),
                    squareImage: squareImage,
                    bannerImage: bannerImage,
                    socials: {}
                )
            case Type<MetadataViews.Royalties>():
                return nil
        }
        return nil
    }

    //------------------------------------------------------------
    // Contract lifecycle
    //------------------------------------------------------------

    // EditionNFT contract initializer
    //
    init() {
        // Set the named paths
        self.CollectionStoragePath = /storage/EditionNFTCollection
        self.CollectionPublicPath = /public/EditionNFTCollection
        self.AdminStoragePath = /storage/EditionNFTAdmin

        // Initialize the entity counts        
        self.totalSupply = 0
        self.totalEditions = 0
        self.nextEditionID = 1

        // Initialize the metadata lookup dictionaries
        self.editionByID <- {}

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.storage.save(<-admin, to: self.AdminStoragePath)
       

        // Let the world know we are here
        emit ContractInitialized()
    }
}
