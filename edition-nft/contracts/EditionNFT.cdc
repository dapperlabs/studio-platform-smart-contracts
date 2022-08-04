/*
    Author: Jude Zhu jude.zhu@dapperlabs.com
*/


import NonFungibleToken from "./NonFungibleToken.cdc"

/*
    There are 2 levels of entity:
    1. Edition
    2. NFT
    
    An Edition is created with metadata. NFTs are minted out of Editions.
 */

// The EditionNFT contract
//
pub contract EditionNFT: NonFungibleToken {
    //------------------------------------------------------------
    // Events
    //------------------------------------------------------------

    // Contract Events
    //
    pub event ContractInitialized()

    // NFT Collection Events
    //
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event Burned(id: UInt64)
    pub event Minted(id: UInt64, editionID: UInt64)

    // Edition Events
    //
    // Emitted when a new edition has been created by an admin.
    pub event EditionCreated(
        id: UInt64, 
        metadata: {String: String}
    )
    // Emitted when an edition is closed by an admin.
    pub event EditionClosed(id: UInt64)

    //------------------------------------------------------------
    // Named values
    //------------------------------------------------------------

    // Named Paths
    //
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let AdminStoragePath:       StoragePath
    pub let MinterPrivatePath:      PrivatePath

    //------------------------------------------------------------
    // Publicly readable contract state
    //------------------------------------------------------------

    // totalSupply
    // The total number of NFTs that in circulation.
    //
    pub var totalSupply:        UInt64


    // totalMinted
    // The total number of NFTs that have been minted.
    //
    pub var totalMinted:        UInt64

    // totalEditions
    // The total number of editions that have been created.
    //
    pub var totalEditions: UInt64

    // nextEditionID
    // The editionID will be assigned to the next edition.
    //
    pub var nextEditionID:      UInt64


    //------------------------------------------------------------
    // Internal contract state
    //------------------------------------------------------------

    // Metadata Dictionaries
    //
    // This is so we can find Edition via ID.
    access(self) let editionByID:       @{UInt64: Edition}

    //------------------------------------------------------------
    // Edition
    //------------------------------------------------------------

    // A public struct to access Edition data
    //
    pub struct EditionData {
        pub let id: UInt64
        pub var numMinted: UInt64
        pub var active: Bool
        pub let metadata: {String: String}

        // initializer
        //
        init (id: UInt64) {
            if let edition = &EditionNFT.editionByID[id] as &EditionNFT.Edition? {
            self.id = id
            self.metadata = edition.metadata
            self.numMinted = edition.numMinted
            self.active = edition.active
            } else {
                panic("edition does not exist")
            }
        }
    }

    // A top level Edition with a unique ID
    //
    pub resource Edition {
        pub let id: UInt64
        // Contents writable if borrowed!
        // This is deliberate, as it allows admins to update the data.
        pub var numMinted: UInt64
        pub let metadata: {String: String}
        pub var active: Bool

        // Close this edition
        //
        pub fun close() {
            pre {
                self.active: "edtion is already closed"
            }

            self.active = false
            emit EditionClosed(id: self.id)
        }

        // Mint a Seasonal NFT in this edition, with the given minting mintingDate.
        // Note that this will panic if the max mint size has already been reached.
        //
        pub fun mint(): @EditionNFT.NFT {
            pre {
                self.active: "edition is already closed. minting is not allowed"
            }

            // Create thek NFT, filled out with our information
            let nft <- create NFT(
                id: EditionNFT.totalMinted + 1,
                editionID: self.id,
                serialNumber: self.numMinted + 1
            )
            EditionNFT.totalSupply = EditionNFT.totalSupply + 1
            EditionNFT.totalMinted = EditionNFT.totalMinted + 1

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
    pub fun getEditionData(id: UInt64): EditionNFT.EditionData {
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
    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        pub let editionID: UInt64

        // Destructor
        //
        destroy() {
            EditionNFT.totalSupply = EditionNFT.totalSupply - 1
            emit Burned(id: self.id)
        }

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
    }

    //------------------------------------------------------------
    // Collection
    //------------------------------------------------------------

    // A public collection interface that allows Moment NFTs to be borrowed
    //
    pub resource interface EditionNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowEditionNFT(id: UInt64): &EditionNFT.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id): 
                    "Cannot borrow Moment NFT reference: The ID of the returned reference is incorrect"
            }
        }
    }

    // An NFT Collection
    //
    pub resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        EditionNFTCollectionPublic
    {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an UInt64 ID field
        //
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        // withdraw removes an NFT from the collection and moves it to the caller
        //
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        //
        pub fun deposit(token: @NonFungibleToken.NFT) {
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
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection) {
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
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // borrowNFT gets a reference to an NFT in the collection
        //
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            pre {
                self.ownedNFTs[id] != nil: "Cannot borrow NFT, no such id"
            }

            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

        // borrowEditionNFT gets a reference to an EditionNFT in the collection
        //
        pub fun borrowEditionNFT(id: UInt64): &EditionNFT.NFT? {
            if self.ownedNFTs[id] != nil {
                if let ref = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT? {
                    return ref! as! &EditionNFT.NFT
                }
                return nil
            } else {
                return nil
            }
        }

        // Collection destructor
        //
        destroy() {
            destroy self.ownedNFTs
        }

        // Collection initializer
        //
        init() {
            self.ownedNFTs <- {}
        }
    }

    // public function that anyone can call to create a new empty collection
    //
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }

    //------------------------------------------------------------
    // Admin
    //------------------------------------------------------------

    // An interface containing the Admin function that allows minting NFTs
    //
    pub resource interface NFTMinter {
        // Mint a single NFT
        // The Edition for the given ID must already exist
        //
        pub fun mintNFT(editionID: UInt64): @EditionNFT.NFT
    }

    // A resource that allows managing metadata and minting NFTs
    //
    pub resource Admin: NFTMinter {

        // Borrow an Edition
        //
        pub fun borrowEdition(id: UInt64): &EditionNFT.Edition {
            pre {
                EditionNFT.editionByID[id] != nil: "Cannot borrow edition, no such id"
            }

            return (&EditionNFT.editionByID[id] as &EditionNFT.Edition?)!
        }

        // Create a Edition 
        //
        pub fun createEdition(metadata: {String: String}): UInt64 {
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
        pub fun closeEdition(id: UInt64): UInt64 {
            if let edition = &EditionNFT.editionByID[id] as &EditionNFT.Edition? {
                edition.close()
                return edition.id
            }
            panic("edition does not exist")
        }

        // Mint a single NFT
        // The Edition for the given ID must already exist
        //
        pub fun mintNFT(editionID: UInt64): @EditionNFT.NFT {
            pre {
                // Make sure the edition we are creating this NFT in exists
                EditionNFT.editionByID.containsKey(editionID): "No such EditionID"
            }
            return <- self.borrowEdition(id: editionID).mint()
        }
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
        self.MinterPrivatePath = /private/EditionNFTMinter

        // Initialize the entity counts        
        self.totalMinted = 0
        self.totalSupply = 0
        self.totalEditions = 0
        self.nextEditionID = 1

        // Initialize the metadata lookup dictionaries
        self.editionByID <- {}

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.save(<-admin, to: self.AdminStoragePath)
        // Link capabilites to the admin constrained to the Minter
        // and Metadata interfaces
        self.account.link<&EditionNFT.Admin{EditionNFT.NFTMinter}>(
            self.MinterPrivatePath,
            target: self.AdminStoragePath
        )

        // Let the world know we are here
        emit ContractInitialized()
    }
}
