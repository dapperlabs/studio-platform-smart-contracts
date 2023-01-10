/*
    Adapted from: AllDay.cdc
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/


import NonFungibleToken from 0xf8d6e0586b0a20c7

/*
    DSSCollection contains collection group & completion functionality. 
    It is designed for use from all Dapper Sports.
*/

// The DSSCollection contract
//
pub contract DSSCollection: NonFungibleToken {
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

    // CollectionGroup Events
    //
    pub event CollectionGroupCreated(id: UInt64, name: String, productPublicPath: PublicPath, startTime: UFix64?, endTime: UFix64?, timeBound: Bool)
    pub event CollectionGroupClosed(id: UInt64)
    pub event NFTAddedToCollectionGroup(nftID: UInt64, collectionGroupID: UInt64)

    // NFT Events
    //
    pub event DSSCollectionNFTMinted(id: UInt64, collectionGroupID: UInt64, serialNumber: UInt64, completedBy: String, completionDate: UFix64)
    pub event DSSCollectionNFTBurned(id: UInt64)


    // Named Paths
    //
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let AdminStoragePath:       StoragePath
    pub let MinterPrivatePath:      PrivatePath

    //------------------------------------------------------------
    // Publicly readable contract state
    //------------------------------------------------------------

    // Entity Counts
    //
    pub var totalSupply:        UInt64
    pub var nextCollectionGroupID:       UInt64


    // Metadata Dictionaries
    //
    // This is so we can find CollectionGroup by their names (via collectionGroupByID)
    access(self) let collectionGroupByID:        @{UInt64: CollectionGroup}

    //------------------------------------------------------------
    // CollectionGroup
    //------------------------------------------------------------

    // A public struct to access CollectionGroup data
    //
    pub struct CollectionGroupData {
        pub let id: UInt64
        pub let name: String
        pub let productPublicPath: PublicPath
        pub let open: Bool
        pub let startTime: UFix64?
        pub let endTime: UFix64?
        pub let timeBound: Bool
        pub var nftIDInCollectionGroup: {UInt64: Bool}

        pub fun nftIDExistsInCollectionGroup(collectionGroupID: UInt64): Bool {
           return self.nftIDInCollectionGroup.containsKey(collectionGroupID)
        }

        init (id: UInt64) {
            if let collectionGroup = &DSSCollection.collectionGroupByID[id] as &DSSCollection.CollectionGroup? {
                self.id = collectionGroup.id
                self.name = collectionGroup.name
                self.productPublicPath = collectionGroup.productPublicPath
                self.open = collectionGroup.open
                self.startTime = collectionGroup.startTime
                self.endTime = collectionGroup.endTime
                self.timeBound = collectionGroup.timeBound
                self.nftIDInCollectionGroup = collectionGroup.nftIDInCollectionGroup
            } else {
                panic("collectionGroup does not exist")
            }
        }
    }

    // A top-level CollectionGroup with a unique ID and name
    //
    pub resource CollectionGroup {
        pub let id: UInt64
        pub let name: String
        pub let productPublicPath: PublicPath
        pub var open: Bool
        pub let startTime: UFix64?
        pub let endTime: UFix64?
        pub let timeBound: Bool
        pub var numMinted: UInt64
        pub var nftIDInCollectionGroup: {UInt64: Bool}

        // Close this collection group
        //
        access(contract) fun close() {
            pre {
                self.open == true: "not open"
            }

            self.open = false

            emit CollectionGroupClosed(id: self.id)
        }

        // Add nftID to collection group
        //
        access(contract) fun addNFTToCollectionGroup(nftID: UInt64) {
            pre {
                self.open == true: "not open"
            }

            self.nftIDInCollectionGroup[nftID] = true

            emit NFTAddedToCollectionGroup(nftID: nftID, collectionGroupID: self.id)
        }

        // Mint a DSSCollection NFT in this group
        //
        pub fun mint(completedBy: String): @DSSCollection.NFT {
            pre {
                self.open != true: "cannot mint an open collection group"
                DSSCollection.validateTimeRange(timeBound: self.timeBound, startTime: self.startTime, endTime: self.endTime) == true : "cannot mint a collection group outside of time bounds"
            }

            // Create the DSSCollection NFT, filled out with our information
            let dssCollectionNFT <- create NFT(
                id: DSSCollection.totalSupply + 1,
                collectionGroupID: self.id,
                serialNumber: self.numMinted + 1,
                completedBy: completedBy
            )
            DSSCollection.totalSupply = DSSCollection.totalSupply + 1
            
            // increment serial number
            self.numMinted = self.numMinted + 1 as UInt64

            return <- dssCollectionNFT
        }

        init (
            name: String, 
            productPublicPath: PublicPath, 
            startTime: UFix64?, 
            endTime: UFix64?, 
            timeBound: Bool
        ) {
            self.id = DSSCollection.nextCollectionGroupID
            self.name = name
            self.productPublicPath = productPublicPath
            self.open = true
            self.startTime = startTime
            self.endTime = endTime
            self.timeBound = timeBound
            self.numMinted = 0 as UInt64
            self.nftIDInCollectionGroup = {}

            // Increment for the nextCollectionGroupID
            DSSCollection.nextCollectionGroupID = self.id + 1 as UInt64

            emit CollectionGroupCreated(id: self.id, name: self.name, productPublicPath: self.productPublicPath, startTime: self.startTime, endTime: self.endTime, timeBound: self.timeBound)
        }
    }

    // Get the publicly available data for a CollectionGroup by id
    //
    pub fun getCollectionGroupData(id: UInt64): DSSCollection.CollectionGroupData {
        pre {
            DSSCollection.collectionGroupByID[id] != nil: "Cannot borrow collection group, no such id"
        }

        return DSSCollection.CollectionGroupData(id: id)
    }

    // Get the publicly available data for a CollectionGroup by id
    //
    pub fun validateTimeRange(timeBound: Bool, startTime: UFix64?, endTime: UFix64?): Bool {
        if timeBound {
            return true
        }

        if startTime! <= getCurrentBlock().timestamp && endTime! >= getCurrentBlock().timestamp {
            return true
        } else {
            return false
        }
    }

    //------------------------------------------------------------
    // NFT
    //------------------------------------------------------------

    // A DSSCollection NFT
    //
    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        pub let collectionGroupID: UInt64
        pub let serialNumber: UInt64
        pub let completionDate: UFix64
        pub let completedBy: String

        destroy() {
            emit DSSCollectionNFTBurned(id: self.id)
        }

        init(
            id: UInt64,
            collectionGroupID: UInt64,
            serialNumber: UInt64,
            completedBy: String
        ) {
            pre {
                DSSCollection.collectionGroupByID[collectionGroupID] != nil: "no such collectionGroupID"
            }

            self.id = id
            self.collectionGroupID = collectionGroupID
            self.serialNumber = serialNumber
            self.completionDate = getCurrentBlock().timestamp
            self.completedBy = completedBy

            emit DSSCollectionNFTMinted(id: self.id, collectionGroupID: self.collectionGroupID, serialNumber: self.serialNumber, completedBy: self.completedBy, completionDate: self.completionDate)
        }
    }

    //------------------------------------------------------------
    // Collection
    //------------------------------------------------------------

    // A public collection interface that allows DSSCollection NFTs to be borrowed
    //
    pub resource interface DSSCollectionNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowDSSCollectionNFT(id: UInt64): &DSSCollection.NFT? {
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
        DSSCollectionNFTCollectionPublic
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
            let token <- token as! @DSSCollection.NFT
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

        // borrowDSSCollectionNFT gets a reference to an NFT in the collection
        //
        pub fun borrowDSSCollectionNFT(id: UInt64): &DSSCollection.NFT? {
            if self.ownedNFTs[id] != nil {
                if let ref = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT? {
                    return ref! as! &DSSCollection.NFT
                }
                return nil
            } else {
                return nil
            }
        }

        destroy() {
            destroy self.ownedNFTs
        }

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
        // The collectionGroupID for the given ID must already exist
        //
        pub fun mintNFT(collectionGroupID: UInt64, completedBy: String): @DSSCollection.NFT
    }

    // A resource that allows managing metadata and minting NFTs
    //
    pub resource Admin: NFTMinter {

        // Borrow a Collection Group
        //
        pub fun borrowCollectionGroup(id: UInt64): &DSSCollection.CollectionGroup {
            pre {
                DSSCollection.collectionGroupByID[id] != nil: "Cannot borrow collection group, no such id"
            }

            return (&DSSCollection.collectionGroupByID[id] as &DSSCollection.CollectionGroup?)!
        }

        // Create a Collection Group
        //
        pub fun createCollectionGroup(name: String, productPublicPath: PublicPath, startTime: UFix64?, endTime: UFix64?, timeBound: Bool): UInt64 {
            // Create and store the new collection group
            let collectionGroup <- create DSSCollection.CollectionGroup(
                name: name,
                productPublicPath: productPublicPath,
                startTime: startTime,
                endTime: endTime,
                timBound: timeBound
            )
            let collectionGroupID = collectionGroup.id
            DSSCollection.collectionGroupByID[collectionGroup.id] <-! collectionGroup

            // Return the new ID for convenience
            return collectionGroupID
        }

        // Close a Collection Group
        //
        pub fun closeCollectionGroup(id: UInt64): UInt64 {
            if let collectionGroup = &DSSCollection.collectionGroupByID[id] as &DSSCollection.CollectionGroup? {
                collectionGroup.close()
                return collectionGroup.id
            }
            panic("collection group does not exist")
        }

        // Add NFT to Collection Group
        //
        pub fun addNFTToCollectionGroup(nftID: UInt64, collectionGroupID: UInt64) {
            if let collectionGroup = &DSSCollection.collectionGroupByID[collectionGroupID] as &DSSCollection.CollectionGroup? {
                collectionGroup.addNFTToCollectionGroup(nftID: nftID)
                return
            }
            panic("collection group does not exist")
        }


        // Mint a single NFT
        // The CollectionGroup for the given ID must already exist
        //
        pub fun mintNFT(collectionGroupID: UInt64, completedBy: String): @DSSCollection.NFT {
            pre {
                // Make sure the collection group exists
                DSSCollection.collectionGroupByID.containsKey(collectionGroupID): "No such CollectionGroupID"
            }
            return <- self.borrowCollectionGroup(id: collectionGroupID).mint(completedBy: completedBy)
        }
    }

    //------------------------------------------------------------
    // Contract lifecycle
    //------------------------------------------------------------

    // DSSCollection contract initializer
    //
    init() {
        // Set the named paths
        self.CollectionStoragePath = /storage/DSSCollectionNFTCollection
        self.CollectionPublicPath = /public/DSSCollectionNFTCollection
        self.AdminStoragePath = /storage/CollectionGroupAdmin
        self.MinterPrivatePath = /private/CollectionGroupMinter

        // Initialize the entity counts
        self.totalSupply = 0
        self.nextCollectionGroupID = 1

        // Initialize the metadata lookup dictionaries
        self.collectionGroupByID <- {}

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.save(<-admin, to: self.AdminStoragePath)
        // Link capabilites to the admin constrained to the Minter
        // and Metadata interfaces
        self.account.link<&DSSCollection.Admin{DSSCollection.NFTMinter}>(
            self.MinterPrivatePath,
            target: self.AdminStoragePath
        )

        // Let the world know we are here
        emit ContractInitialized()
    }
}
