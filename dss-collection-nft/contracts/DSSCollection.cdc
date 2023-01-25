/*
    DSSCollection contains collection group & completion functionality for DSS.
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import NonFungibleToken from "./NonFungibleToken.cdc"
import MetadataViews from 0xMETADATAVIEWSADDRESS

// The DSSCollection contract
//
pub contract DSSCollection: NonFungibleToken {

    // Contract Events
    //
    pub event ContractInitialized()

    // NFT Collection Events
    //
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)

    // Events
    //
    pub event CollectionGroupCreated(
        id: UInt64,
        name: String,
        description: String,
        typeName: String,
        startTime: UFix64?,
        endTime: UFix64?,
        timeBound: Bool
    )
    pub event CollectionGroupClosed(id: UInt64)
    pub event ItemCreated(
        id: UInt64,
        itemID: UInt64,
        points: UInt64,
        itemType: String
    )
    pub event SlotCreated(
        id: UInt64,
        collectionGroupID: UInt64,
        logicalOperator: String,
        typeName: String
    )
    pub event ItemAddedToSlot(
        id: UInt64,
        slotID: UInt64,
        collectionGroupID: UInt64
    )
    pub event CollectionNFTMinted(
        id: UInt64,
        collectionGroupID: UInt64,
        serialNumber: UInt64,
        completedBy: String,
        completionDate: UFix64,
        level: UInt8
    )
    pub event CollectionNFTBurned(id: UInt64)

    // Named Paths
    //
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let AdminStoragePath:       StoragePath
    pub let MinterPrivatePath:      PrivatePath

    // Entity Counts
    //
    pub var totalSupply:    UInt64
    pub var nextItemID:     UInt64

    // Lists in contract
    //
    access(self) let collectionGroupByID: @{UInt64: CollectionGroup}
    access(self) let slotByID: @{UInt64: Slot}
    access(self) let itemByID: {UInt64: Item}

    // A public struct to access Item data
    //
    pub struct Item {
        pub let id: UInt64 // unique uuid of resource
        pub let itemID: UInt64 // the id of the edition, tier, play
        pub let points: UInt64 // points for item
        pub let itemType: String // (edition.id, edition.tier, play.id)

        init (
            id: UInt64,
            itemID: UInt64,
            points: UInt64,
            itemType: String
        ) {
            self.id = id
            self.itemID = itemID
            self.points = points
            self.itemType = itemType

            DSSCollection.nextItemID = self.id + 1 as UInt64

            emit ItemCreated(
                id: self.id,
                itemID: self.itemID,
                points: self.points,
                itemType: self.itemType
            )
        }
    }

    // A public struct to access Slot data
    //
    pub struct SlotData {
        pub let id: UInt64
        pub let collectionGroupID: UInt64
        pub let logicalOperator: String // (AND / OR)
        pub let typeName: String // (A.contractAddress.NFT...)
        pub var items: [Item]

        init (id: UInt64) {
            if let slot = &DSSCollection.slotByID[id] as &DSSCollection.Slot? {
                self.id = slot.id
                self.collectionGroupID = slot.collectionGroupID
                self.logicalOperator = slot.logicalOperator
                self.typeName = slot.typeName
                self.items = slot.items
            } else {
                panic("slot does not exist")
            }
        }
    }

    // A top-level Slot with a unique ID
    //
    pub resource Slot {
        pub let id: UInt64
        pub let collectionGroupID: UInt64
        pub let logicalOperator: String // (AND / OR)
        pub let typeName: String // (A.contractAddress.NFT...)
        pub var items: [Item]

        // Add item to slot
        //
        access(contract) fun addItemToSlot(id: UInt64) {
            pre {
                DSSCollection.CollectionGroupData(
                    id: self.collectionGroupID
                ).open == true: "collection group not open"
            }

            let item = DSSCollection.getItem(id: id)
            self.items.append(item)

            emit ItemAddedToSlot(
                 id: item.id,
                 slotID: self.id,
                 collectionGroupID: self.collectionGroupID
             )
        }

        init (
            collectionGroupID: UInt64,
            logicalOperator: String,
            typeName: String
        ) {
            pre {
                DSSCollection.CollectionGroupData(
                    id: collectionGroupID
                ).open == true: "collection group not open"
            }

            self.id = self.uuid
            self.collectionGroupID = collectionGroupID
            self.logicalOperator = logicalOperator
            self.typeName = typeName
            self.items = []

            emit SlotCreated(
                id: self.id,
                collectionGroupID: self.collectionGroupID,
                logicalOperator: self.logicalOperator,
                typeName: self.typeName
            )
        }
    }

    // A public struct to access CollectionGroup data
    //
    pub struct CollectionGroupData {
        pub let id: UInt64
        pub let name: String
        pub let description: String
        pub let typeName: String
        pub let open: Bool
        pub let startTime: UFix64?
        pub let endTime: UFix64?
        pub let timeBound: Bool

        init (id: UInt64) {
            if let collectionGroup = &DSSCollection.collectionGroupByID[id] as &DSSCollection.CollectionGroup? {
                self.id = collectionGroup.id
                self.name = collectionGroup.name
                self.description = collectionGroup.description
                self.typeName = collectionGroup.typeName
                self.open = collectionGroup.open
                self.startTime = collectionGroup.startTime
                self.endTime = collectionGroup.endTime
                self.timeBound = collectionGroup.timeBound
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
        pub let description: String
        pub let typeName: String
        pub var open: Bool
        pub let startTime: UFix64?
        pub let endTime: UFix64?
        pub let timeBound: Bool
        pub var numMinted: UInt64

        // Close this collection group
        //
        access(contract) fun close() {
            pre {
                self.open :  "Already closed"
            }

            self.open = false

            emit CollectionGroupClosed(id: self.id)
        }

        // Mint a DSSCollection NFT in this group
        //
        pub fun mint(completedBy: String, level: UInt8): @DSSCollection.NFT {
            pre {
                !self.open : "Cannot mint an open collection group"
                DSSCollection.validateTimeRange(
                    timeBound: self.timeBound,
                    startTime: self.startTime,
                    endTime: self.endTime
                ) == true : "cannot mint a collection group outside of time bounds"
                level <= 10: "token level must be less than 10"
            }

            // Create the DSSCollection NFT, filled out with our information
            //
            let dssCollectionNFT <- create NFT(
                collectionGroupID: self.id,
                serialNumber: self.numMinted + 1,
                completedBy: completedBy,
                level: level
            )
            DSSCollection.totalSupply = DSSCollection.totalSupply + 1
            self.numMinted = self.numMinted + 1 as UInt64

            return <- dssCollectionNFT
        }

        init (
            name: String,
            description: String,
            typeName: String,
            startTime: UFix64?,
            endTime: UFix64?,
            timeBound: Bool
        ) {
            self.id = self.uuid
            self.name = name
            self.description = description
            self.typeName = typeName
            self.open = true
            self.startTime = startTime
            self.endTime = endTime
            self.timeBound = timeBound
            self.numMinted = 0 as UInt64

            emit CollectionGroupCreated(
                id: self.id,
                name: self.name,
                description: self.description,
                typeName: self.typeName,
                startTime: self.startTime,
                endTime: self.endTime,
                timeBound: self.timeBound
            )
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

    // Get the publicly available data for a Slot by id
    //
    pub fun getSlotData(id: UInt64): DSSCollection.SlotData {
        pre {
            DSSCollection.slotByID[id] != nil: "Cannot borrow slot, no such id"
        }

        return DSSCollection.SlotData(id: id)
    }

    // Get the publicly available data for a Slot by id
    //
    pub fun getItem(id: UInt64): DSSCollection.Item {
        return DSSCollection.itemByID[id]!
    }

    // Validate time range of collection group
    //
    pub fun validateTimeRange(timeBound: Bool, startTime: UFix64?, endTime: UFix64?): Bool {
        if !timeBound {
            return true
        }

        if startTime! <= getCurrentBlock().timestamp && endTime! >= getCurrentBlock().timestamp {
            return true
        } else {
            return false
        }
    }

    // A DSSCollection NFT
    //
    pub resource NFT: NonFungibleToken.INFT, MetadataViews.Resolver {
        pub let id: UInt64
        pub let collectionGroupID: UInt64
        pub let serialNumber: UInt64
        pub let completionDate: UFix64
        pub let completedBy: String
        pub let level: UInt8

        pub fun name(): String {
            let collectionGroupData: DSSCollection.CollectionGroupData
                = DSSCollection.getCollectionGroupData(id: self.collectionGroupID)
            let level: String = self.level.toString()
            return collectionGroupData.name
                .concat(" Level ")
                .concat(level)
                .concat(" Completion Token")
        }

        pub fun description(): String {
            let serialNumber: String = self.serialNumber.toString()
            let completionDate: String = self.completionDate.toString()
            return "Completed by "
                .concat(self.completedBy)
                .concat(" on ")
                .concat(completionDate)
                .concat(" with serial number ")
                .concat(serialNumber)
        }

        destroy() {
            emit CollectionNFTBurned(id: self.id)
        }

        pub fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>()
            ]
        }

        pub fun resolveView(_ view: Type): AnyStruct? {
            return MetadataViews.Display(
                name: self.name(),
                description: self.description(),
                thumbnail: MetadataViews.HTTPFile(
                    url:"https://storage.googleapis.com/dl-nfl-assets-prod/static/images/collection-group/token-placeholder.png"
                )
            )
        }

        init(
            collectionGroupID: UInt64,
            serialNumber: UInt64,
            completedBy: String,
            level: UInt8
        ) {
            pre {
                DSSCollection.collectionGroupByID[collectionGroupID] != nil: "no such collectionGroupID"
            }

            self.id = self.uuid
            self.collectionGroupID = collectionGroupID
            self.serialNumber = serialNumber
            self.completionDate = getCurrentBlock().timestamp
            self.completedBy = completedBy
            self.level = level

            emit CollectionNFTMinted(
                id: self.id,
                collectionGroupID: self.collectionGroupID,
                serialNumber: self.serialNumber,
                completedBy: self.completedBy,
                completionDate: self.completionDate,
                level: self.level,
            )
        }
    }

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
        DSSCollectionNFTCollectionPublic,
        MetadataViews.ResolverCollection
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
            let keys = tokens.getIDs()

            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

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

        pub fun borrowViewResolver(id: UInt64): &AnyResource{MetadataViews.Resolver} {
            let nft = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
            let dssNFT = nft as! &DSSCollection.NFT
            return dssNFT as &AnyResource{MetadataViews.Resolver}
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

    // An interface containing the Admin function that allows minting NFTs
    //
    pub resource interface NFTMinter {
        pub fun mintNFT(collectionGroupID: UInt64, completedBy: String, level: UInt8): @DSSCollection.NFT
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

        // Borrow a Slot
        //
        pub fun borrowSlot(id: UInt64): &DSSCollection.Slot {
            pre {
                DSSCollection.slotByID[id] != nil: "Cannot borrow slot, no such id"
            }

            return (&DSSCollection.slotByID[id] as &DSSCollection.Slot?)!
        }

        // Borrow a Item
        //
        pub fun borrowItem(id: UInt64): &DSSCollection.Item {
            pre {
                DSSCollection.itemByID[id] != nil: "Cannot borrow item, no such id"
            }

            return (&DSSCollection.itemByID[id] as &DSSCollection.Item?)!
        }

        // Create a Collection Group
        //
        pub fun createCollectionGroup(
            name: String,
            description: String,
            typeName: String,
            startTime: UFix64?,
            endTime: UFix64?,
            timeBound: Bool
        ): UInt64 {
            let collectionGroup <- create DSSCollection.CollectionGroup(
                name: name,
                description: description,
                typeName: typeName,
                startTime: startTime,
                endTime: endTime,
                timeBound: timeBound
            )
            let collectionGroupID = collectionGroup.id
            DSSCollection.collectionGroupByID[collectionGroup.id] <-! collectionGroup

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

        // Create a Slot
        //
        pub fun createSlot(
            collectionGroupID: UInt64,
            logicalOperator: String,
            typeName: String
        ): UInt64 {
            let slot <- create DSSCollection.Slot(
                collectionGroupID: collectionGroupID,
                logicalOperator: logicalOperator,
                typeName: typeName
            )
            let slotID = slot.id
            DSSCollection.slotByID[slot.id] <-! slot
            return slotID
        }

        // Create an Item
        //
        pub fun createItem(
            itemID: UInt64,
            points: UInt64,
            itemType: String
        ): UInt64 {
            let item = DSSCollection.Item(
                id: DSSCollection.nextItemID,
                itemID: itemID,
                points: points,
                itemType: itemType
            )

            DSSCollection.itemByID[item.id] = item
            return item.id
        }

        // Add Item to Slot
        //
        pub fun addItemToSlot(slotID: UInt64, id: UInt64) {
            if let slot = &DSSCollection.slotByID[slotID] as &DSSCollection.Slot? {
                slot.addItemToSlot(id: id)
                return
            }
            panic("slot does not exist")
        }


        // Mint a single NFT
        // The CollectionGroup for the given ID must already exist
        //
        pub fun mintNFT(collectionGroupID: UInt64, completedBy: String, level: UInt8): @DSSCollection.NFT {
            pre {
                // Make sure the collection group exists
                DSSCollection.collectionGroupByID.containsKey(collectionGroupID): "No such CollectionGroupID"
            }
            return <- self.borrowCollectionGroup(id: collectionGroupID).mint(completedBy: completedBy, level: level)
        }
    }

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
        self.nextItemID = 1

        // Initialize the metadata lookup dictionaries
        self.collectionGroupByID <- {}
        self.slotByID <- {}
        self.itemByID = {}

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.save(<-admin, to: self.AdminStoragePath)
        // Link capabilites to the admin constrained to the Minter
        // and Metadata interfaces
        self.account.link<&DSSCollection.Admin{DSSCollection.NFTMinter}>(
            self.MinterPrivatePath,
            target: self.AdminStoragePath
        )

        emit ContractInitialized()
    }
}
