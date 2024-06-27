/*
    Escrow Contract for managing NFTs in a Leaderboard Context.
    Holds NFTs in Escrow account awaiting transfer or burn.

    Authors:
        Corey Humeston: corey.humeston@dapperlabs.com
        Deewai Abdullahi: innocent.abdullahi@dapperlabs.com
*/

import NonFungibleToken from "NonFungibleToken"

access(all) contract Escrow {
    // Event emitted when a new leaderboard is created.
    access(all) event LeaderboardCreated(name: String, nftType: Type)

    // Event emitted when an NFT is deposited to a leaderboard.
    access(all) event EntryDeposited(leaderboardName: String, nftID: UInt64, owner: Address)

    // Event emitted when an NFT is returned to the original collection from a leaderboard.
    access(all) event EntryReturnedToCollection(leaderboardName: String, nftID: UInt64, owner: Address)

    // Event emitted when an NFT is burned from a leaderboard.
    access(all) event EntryBurned(leaderboardName: String, nftID: UInt64)

    // Named Paths
    access(all) let CollectionStoragePath: StoragePath
    access(all) let CollectionPublicPath: PublicPath

    // The resource representing an NFT leaderboard's info.
    access(all) struct LeaderboardInfo {
        access(all) let name: String
        access(all) let nftType: Type
        access(all) let entriesLength: Int

        // LeaderboardInfo struct initializer.
        init(name: String, nftType: Type, entriesLength: Int) {
            self.name = name
            self.nftType = nftType
            self.entriesLength = entriesLength
        }
    }

    // The resource representing a leaderboard.
    access(all) resource Leaderboard {
        access(all) var collection: @{NonFungibleToken.Collection}
        access(all) var entriesData: {UInt64: LeaderboardEntry}
        access(all) let name: String
        access(all) let nftType: Type
        access(all) var entriesLength: Int
        access(all) var metadata: {String: AnyStruct}

        // Adds an NFT entry to the leaderboard.
        access(contract) fun addEntryToLeaderboard(nft: @{NonFungibleToken.NFT}, ownerAddress: Address, metadata: {String: AnyStruct}) {
            pre {
                nft.isInstance(self.nftType): "This NFT cannot be used for leaderboard. NFT is not of the correct type."
            }

            let nftID = nft.id

            // Create the entry and add it to the entries map
            let entry = LeaderboardEntry(
                nftID: nftID,
                ownerAddress: ownerAddress,
                metadata: metadata,
            )

            // Store the entry data in the leaderboard.
            self.entriesData[nftID] = entry

            // Deposit the NFT into the leaderboard's NFT collection.
            self.collection.deposit(token: <-nft)

            // Increment entries length.
            self.entriesLength = self.entriesLength + 1

            emit EntryDeposited(leaderboardName: self.name, nftID: nftID, owner: ownerAddress)
        }

        // Withdraws an NFT entry from the leaderboard.
        access(contract) fun transferNftToCollection(nftID: UInt64, depositCap: Capability<&{NonFungibleToken.Collection}>) {
            pre {
                self.entriesData[nftID] != nil : "Entry does not exist with this NFT ID"
                depositCap.address == self.entriesData[nftID]!.ownerAddress : "Only the owner of the entry can withdraw it"
                depositCap.check() : "Deposit capability is not valid"
            }

            // Remove the NFT entry's data from the leaderboard.
            self.entriesData.remove(key: nftID)!

            // Transfer the NFT to the receiver's collection.
            let receiverCollection = depositCap.borrow()
                ?? panic("Could not borrow the NFT receiver from the capability")
            receiverCollection.deposit(token: <- self.collection.withdraw(withdrawID: nftID))
            emit EntryReturnedToCollection(leaderboardName: self.name, nftID: nftID, owner: depositCap.address)

            // Decrement entries length.
            self.entriesLength = self.entriesLength - 1
        }

        // Burns an NFT entry from the leaderboard.
        access(contract) fun burn(nftID: UInt64) {
            pre {
                self.entriesData[nftID] != nil : "Entry does not exist with this NFT ID"
            }

            // Remove the NFT entry's data from the leaderboard.
            self.entriesData.remove(key: nftID)!

            // Burn the NFT.
            destroy <- self.collection.withdraw(withdrawID: nftID)
            emit EntryBurned(leaderboardName: self.name, nftID: nftID)

            // Decrement entries length.
            self.entriesLength = self.entriesLength - 1
        }

        // Leaderboard resource initializer.
        init(name: String, nftType: Type, collection: @{NonFungibleToken.Collection}) {
            self.name = name
            self.nftType = nftType
            self.collection <- collection
            self.entriesLength = 0
            self.metadata = {}
            self.entriesData = {}
        }
    }

    // The resource representing an NFT entry in a leaderboard.
    access(all) struct LeaderboardEntry {
        access(all) let nftID: UInt64
        access(all) let ownerAddress: Address
        access(all) var metadata: {String: AnyStruct}

        // LeaderboardEntry struct initializer.
        init(nftID: UInt64, ownerAddress: Address, metadata: {String: AnyStruct}) {
            self.nftID = nftID
            self.ownerAddress = ownerAddress
            self.metadata = metadata
        }
    }

    // An interface containing the Collection function that gets leaderboards by name.
    access(all) resource interface ICollectionPublic {
        access(all) fun getLeaderboardInfo(name: String): LeaderboardInfo?
        access(all) fun addEntryToLeaderboard(nft: @{NonFungibleToken.NFT}, leaderboardName: String, ownerAddress: Address, metadata: {String: AnyStruct})
        access(Operate) fun createLeaderboard(name: String, nftType: Type, collection: @{NonFungibleToken.Collection})
        access(Operate) fun transferNftToCollection(leaderboardName: String, nftID: UInt64, depositCap: Capability<&{NonFungibleToken.Collection}>)
        access(Operate) fun burn(leaderboardName: String, nftID: UInt64)
    }

    // Entitlement that grants the ability to operate the Escrow Collection
    access(all) entitlement Operate

    // Deprecated in favor of Operate entitlement
    access(all) resource interface ICollectionPrivate: ICollectionPublic {}

    // The resource representing a collection.
    access(all) resource Collection: ICollectionPublic, ICollectionPrivate {
        // A dictionary holding leaderboards.
        access(self) var leaderboards: @{String: Leaderboard}

        // Creates a new leaderboard and stores it.
        access(Operate) fun createLeaderboard(name: String, nftType: Type, collection: @{NonFungibleToken.Collection}) {
            if self.leaderboards[name] != nil {
                panic("Leaderboard already exists with this name")
            }

            // Create and store leaderboard resource in the leaderboards dictionary.
            self.leaderboards[name] <-! create Leaderboard(name: name, nftType: nftType, collection: <-collection)

            // Emit the event.
            emit LeaderboardCreated(name: name, nftType: nftType)
        }

        // Returns leaderboard info with the given name.
        access(all) fun getLeaderboardInfo(name: String): LeaderboardInfo? {
            let leaderboard = &self.leaderboards[name] as &Leaderboard?
            if leaderboard == nil {
                return nil
            }

            return LeaderboardInfo(
                name: leaderboard!.name,
                nftType: leaderboard!.nftType,
                entriesLength: leaderboard!.entriesLength
            )
        }

        // Call addEntry.
        access(all) fun addEntryToLeaderboard(nft: @{NonFungibleToken.NFT}, leaderboardName: String, ownerAddress: Address, metadata: {String: AnyStruct}) {
            let leaderboard = &self.leaderboards[leaderboardName] as &Leaderboard?
            if leaderboard == nil {
                panic("Leaderboard does not exist with this name")
            }

            leaderboard!.addEntryToLeaderboard(nft: <-nft, ownerAddress: ownerAddress, metadata: metadata)
        }

        // Calls transferNftToCollection.
        access(Operate) fun transferNftToCollection(leaderboardName: String, nftID: UInt64, depositCap: Capability<&{NonFungibleToken.Collection}>) {
            let leaderboard = &self.leaderboards[leaderboardName] as &Leaderboard?
            if leaderboard == nil {
                panic("Leaderboard does not exist with this name")
            }

            leaderboard!.transferNftToCollection(nftID: nftID, depositCap: depositCap)
        }

        // Calls burn.
        access(Operate) fun burn(leaderboardName: String, nftID: UInt64) {
            let leaderboard = &self.leaderboards[leaderboardName] as &Leaderboard?
            if leaderboard == nil {
                panic("Leaderboard does not exist with this name")
            }

            leaderboard!.burn(nftID: nftID)
        }

        // Collection resource initializer.
        init() {
            self.leaderboards <- {}
        }
    }

    // Escrow contract initializer.
    init() {
        // Initialize paths.
        self.CollectionStoragePath = /storage/EscrowLeaderboardCollection
        self.CollectionPublicPath = /public/EscrowLeaderboardCollectionInfo

        // Create and save a Collection resource to account storage.
        self.account.storage.save(<- create Collection(), to: self.CollectionStoragePath)

        // Create a public capability to the Collection resource and publish it publicly.
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&Collection>(self.CollectionStoragePath),
            at: self.CollectionPublicPath
        )
    }
}
