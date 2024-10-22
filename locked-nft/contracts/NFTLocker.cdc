import NonFungibleToken from "NonFungibleToken"


/// A contract to lock NFT for a given duration
/// Locked NFT are stored in a user owned collection
/// The collection owner can unlock the NFT after duration has been exceeded
///
access(all) contract NFTLocker {

    /// Contract events
    ///
    access(all) event Withdraw(id: UInt64, from: Address?)
    access(all) event Deposit(id: UInt64, to: Address?)
    access(all) event NFTLocked(
        id: UInt64,
        to: Address?,
        lockedAt: UInt64,
        lockedUntil: UInt64,
        duration: UInt64,
        nftType: Type
    )
    access(all) event NFTUnlocked(
        id: UInt64,
        from: Address?,
        nftType: Type,
        receiverName: String?, // receiver name if unlocked with authorized deposit, nil otherwise
        lockedUntilBeforeEarlyUnlock: UInt64? // lockedUntil if unlocked with authorized deposit, nil otherwise
    )
    access(all) event ReceiverAdded(name: String, eligibleNFTTypes: {Type: Bool})
    access(all) event ReceiverRemoved(name: String, eligibleNFTTypes: {Type: Bool})

    /// Named Paths
    ///
    access(all) let CollectionStoragePath:  StoragePath
    access(all) let CollectionPublicPath:   PublicPath

    /// Contract variables
    ///
    access(all) var totalLockedTokens:      UInt64

    /// Metadata Dictionaries
    ///
    access(self) let lockedTokens:  {Type: {UInt64: LockedData}}

    /// Entitlement that grants the ability to operate authorized functions
    ///
    access(all) entitlement Operate

    /// Data describing characteristics of the locked NFT
    ///
    access(all) struct LockedData {
        access(all) let id: UInt64
        access(all) let owner: Address
        access(all) let lockedAt: UInt64
        access(all) let lockedUntil: UInt64
        access(all) let duration: UInt64
        access(all) let nftType: Type
        access(all) let extension: {String: AnyStruct}

        view init (id: UInt64, owner: Address, duration: UInt64, nftType: Type) {
            if let lockedToken = (NFTLocker.lockedTokens[nftType]!)[id] {
                self.id = id
                self.owner = lockedToken.owner
                self.lockedAt = lockedToken.lockedAt
                self.lockedUntil = lockedToken.lockedUntil
                self.duration = lockedToken.duration
                self.nftType = lockedToken.nftType
                self.extension = lockedToken.extension
            } else {
                self.id = id
                self.owner = owner
                self.lockedAt = UInt64(getCurrentBlock().timestamp)
                self.lockedUntil = self.lockedAt + duration
                self.duration = duration
                self.nftType = nftType
                self.extension = {}
            }
        }
    }

    /// Get the details of a locked NFT
    ///
    access(all) view fun getNFTLockerDetails(id: UInt64, nftType: Type): NFTLocker.LockedData? {
        return (NFTLocker.lockedTokens[nftType]!)[id]
    }

    /// Determine if NFT can be unlocked
    ///
    access(all) view fun canUnlockToken(id: UInt64, nftType: Type): Bool {
        if let lockedTokens = &NFTLocker.lockedTokens[nftType] as &{UInt64: NFTLocker.LockedData}? {
            if let lockedToken = lockedTokens[id] {
                if lockedToken.lockedUntil <= UInt64(getCurrentBlock().timestamp) {
                    return true
                }
            }
        }
        return false
    }

    /// The path to the Admin resource belonging to the account where this contract is deployed
    ///
    access(all) view fun GetAdminStoragePath(): StoragePath {
        return /storage/NFTLockerAdmin
    }

    /// The path to the ReceiverCollector resource belonging to the account where this contract is deployed
    ///
    access(all) view fun getReceiverCollectorStoragePath(): StoragePath {
        return /storage/NFTLockerAdminReceiverCollector
    }

    /// Return an unauthorized reference to the admin's ReceiverCollector resource if it exists
    ///
    access(all) view fun borrowAdminReceiverCollectorPublic(): &ReceiverCollector? {
        return self.account.storage.borrow<&ReceiverCollector>(from: NFTLocker.getReceiverCollectorStoragePath())
    }

    /// Interface for depositing NFTs to authorized receivers
    ///
    access(all) struct interface IAuthorizedDepositHandler {
        access(all) fun deposit(nft: @{NonFungibleToken.NFT}, ownerAddress: Address, passThruParams: {String: AnyStruct})
    }

    /// Struct that defines a Receiver
    ///
    /// Receivers are entities that can receive locked NFTs and deposit them using a specific deposit method
    ///
    access(all) struct Receiver {
        /// Handler for depositing NFTs for the receiver
        ///
        access(all) var authorizedDepositHandler: {IAuthorizedDepositHandler}

        /// The eligible NFT types for the receiver
        ///
        access(all) let eligibleNFTTypes: {Type: Bool}

        /// Extension map for additional data
        ///
        access(all) let metadata: {String: AnyStruct}

        /// Initialize Receiver struct
        ///
        view init(
            authorizedDepositHandler: {IAuthorizedDepositHandler},
            eligibleNFTTypes: {Type: Bool}
        ) {
            self.authorizedDepositHandler = authorizedDepositHandler
            self.eligibleNFTTypes = eligibleNFTTypes
            self.metadata = {}
        }
    }

    /// ReceiverCollector resource
    ///
    /// Note: This resource is used to store receivers and corresponding authorized deposit handlers; currently,
    /// only the admin account can add or remove receivers - in the future, a ReceiverProvider resource could
    /// be added to provide this capability to separate authorized accounts.
    ///
    access(all) resource ReceiverCollector  {
        /// Map of receivers by name
        ///
        access(self) let receiversByName: {String: Receiver}

        /// Map of receiver names by NFT type for lookup
        ///
        access(self) let receiverNamesByNFTType: {Type: {String: Bool}}

        /// Extension map for additional data
        ///
        access(self) let metadata: {String: AnyStruct}

        /// Add a deposit handler for given NFT types
        ///
        access(Operate) fun addReceiver(
            name: String,
            authorizedDepositHandler: {IAuthorizedDepositHandler},
            eligibleNFTTypes: {Type: Bool}
        ) {
            pre {
                !self.receiversByName.containsKey(name): "Receiver with the same name already exists"
            }

            // Add the receiver
            self.receiversByName[name] = Receiver(
                authorizedDepositHandler: authorizedDepositHandler,
                eligibleNFTTypes: eligibleNFTTypes
            )

            // Add the receiver to the lookup map
            for nftType in eligibleNFTTypes.keys {
                if let namesMap = self.receiverNamesByNFTType[nftType] {
                    namesMap[name] = true
                    self.receiverNamesByNFTType[nftType] = namesMap
                } else {
                    self.receiverNamesByNFTType[nftType] = {name: true}
                }
            }

            // Emit event
            emit ReceiverAdded(name: name, eligibleNFTTypes: eligibleNFTTypes)
        }

        /// Remove a deposit method for a given NFT type
        ///
        access(Operate) fun removeReceiver(name: String) {
            // Get the receiver
            let receiver = self.receiversByName[name]
                ?? panic("Receiver with the given name does not exist")

            // Remove the receiver from the lookup map
            for nftType in receiver.eligibleNFTTypes.keys {
                if self.receiverNamesByNFTType.containsKey(nftType) {
                    self.receiverNamesByNFTType[nftType]!.remove(key: name)
                }
            }

            // Remove the receiver
            self.receiversByName.remove(key: name)

            // Emit event
            emit ReceiverRemoved(name: name, eligibleNFTTypes: receiver.eligibleNFTTypes)
        }

        /// Get the receiver for the given name if it exists
        ///
        access(all) view fun getReceiver(name: String): Receiver? {
            return self.receiversByName[name]
        }

        /// Get the receiver names for the given NFT type if it exists
        ///
        access(all) view fun getReceiverNamesByNFTType(nftType: Type): {String: Bool}? {
            return self.receiverNamesByNFTType[nftType]
        }

        /// Initialize ReceiverCollector struct
        ///
        view init() {
            self.receiversByName = {}
            self.receiverNamesByNFTType = {}
            self.metadata = {}
        }
    }

    /// Admin resource
    ///
    access(all) resource Admin {
        /// Expire lock
        ///
        access(all) fun expireLock(id: UInt64, nftType: Type) {
            NFTLocker.expireLock(id: id, nftType: nftType)
        }

        /// Create and return a ReceiverCollector resource
        ///
        access(all) fun createReceiverCollector(): @ReceiverCollector {
            return <- create ReceiverCollector()
        }
    }

    /// Expire lock
    ///
    /// This can be called either by the admin or by the user unlockWithAuthorizedDeposit, if the locked NFT
    /// type is eligible.
    ///
    access(contract) fun expireLock(id: UInt64, nftType: Type) {
        if let locker = &NFTLocker.lockedTokens[nftType] as auth(Mutate) &{UInt64: NFTLocker.LockedData}?{
            if locker[id] != nil {
                // Update locked data's duration to 0
                if let oldLockedData = locker.remove(key: id){
                    locker.insert(
                        key: id,
                        LockedData(
                            id: id,
                            owner: oldLockedData.owner,
                            duration: 0,
                            nftType: nftType
                        )
                    )
                }
            }
        }
    }


    /// A public collection interface that requires the ability to lock and unlock NFTs and return the ids
    /// of NFTs locked for a given type
    ///
    access(all) resource interface LockedCollection {
        access(all) view fun getIDs(nftType: Type): [UInt64]?
        access(Operate) fun lock(token: @{NonFungibleToken.NFT}, duration: UInt64)
        access(Operate) fun unlock(id: UInt64, nftType: Type): @{NonFungibleToken.NFT}
        access(Operate) fun unlockWithAuthorizedDeposit(
            id: UInt64,
            nftType: Type,
            receiverName: String,
            passThruParams: {String: AnyStruct}
        )
    }

    /// Deprecated in favor of Operate entitlement
    ///
    access(all) resource interface LockProvider: LockedCollection {}

    /// An NFT Collection
    ///
    access(all) resource Collection: LockedCollection, LockProvider {
        /// This collection's locked NFTs
        ///
        access(all) var lockedNFTs: @{Type: {UInt64: {NonFungibleToken.NFT}}}

        /// Unlock an NFT of a given type
        ///
        access(Operate) fun unlock(id: UInt64, nftType: Type): @{NonFungibleToken.NFT} {
            pre {
                NFTLocker.canUnlockToken(id: id, nftType: nftType): "locked duration has not been met"
            }

            return <- self.withdrawFromLockedNFTs(id: id, nftType: nftType, receiverName: nil, lockedUntilBeforeEarlyUnlock: nil)
        }

        /// Force unlock the NFT with the given id and type, and deposit it using the receiver's deposit method;
        /// additional function parameters may be required by the receiver's deposit method and are passed in the
        /// passThruParams map.
        ///
        access(Operate) fun unlockWithAuthorizedDeposit(
            id: UInt64,
            nftType: Type,
            receiverName: String,
            passThruParams: {String: AnyStruct}
        ) {
            pre {
                !NFTLocker.canUnlockToken(id: id, nftType: nftType): "locked duration has been met, use unlock() instead"
            }

            // Get the locked token details, panic if it doesn't exist
            let lockedTokenDetails = NFTLocker.getNFTLockerDetails(id: id, nftType: nftType)
                ?? panic("No locked token found for the given id and NFT type")

            // Get a public reference to the admin's receiver collector, panic if it doesn't exist
            let receiverCollector = NFTLocker.borrowAdminReceiverCollectorPublic()
                ?? panic("No receiver collector found")

            // Get the receiver names for the given NFT type, panic if there is no record
            let nftTypeReceivers = receiverCollector.getReceiverNamesByNFTType(nftType: nftType)
                ?? panic("No authorized receiver for the given NFT type")

            // Verify that the receiver with the given name is authorized
            assert(
                nftTypeReceivers[receiverName] == true,
                message: "Provided receiver does not exist or is not authorized for the given NFT type"
            )

            // Expire the NFT's lock
            NFTLocker.expireLock(id: id, nftType: nftType)

            // Unlock and deposit the NFT using the receiver's deposit method
            receiverCollector.getReceiver(name: receiverName)!.authorizedDepositHandler.deposit(
                nft: <- self.withdrawFromLockedNFTs(
                    id: id,
                    nftType: nftType,
                    receiverName: receiverName,
                    lockedUntilBeforeEarlyUnlock: lockedTokenDetails.lockedUntil
                ),
                ownerAddress: lockedTokenDetails.owner,
                passThruParams: passThruParams,
            )
        }

        /// Withdraw the NFT with the given id and type, used in the unlock and unlockWithAuthorizedDeposit functions
        ///
        access(self) fun withdrawFromLockedNFTs(id: UInt64, nftType: Type, receiverName: String?, lockedUntilBeforeEarlyUnlock: UInt64?): @{NonFungibleToken.NFT} {
            // Remove the token's locked data
            if let lockedTokens = &NFTLocker.lockedTokens[nftType] as auth(Remove) &{UInt64: NFTLocker.LockedData}? {
                lockedTokens.remove(key: id)
            }

            // Decrement the locked tokens count
            NFTLocker.totalLockedTokens = NFTLocker.totalLockedTokens - 1

            // Emit events
            emit NFTUnlocked(
                id: id,
                from: self.owner?.address,
                nftType: nftType,
                receiverName: receiverName,
                lockedUntilBeforeEarlyUnlock: lockedUntilBeforeEarlyUnlock
            )
            emit Withdraw(id: id, from: self.owner?.address)

            return <- self.lockedNFTs[nftType]?.remove(key: id)!!
        }

        /// Lock the given NFT for the specified duration
        ///
        access(Operate) fun lock(token: @{NonFungibleToken.NFT}, duration: UInt64) {
            // Get the NFT's id and type
            let nftId: UInt64 = token.id
            let nftType: Type = token.getType()

            // Initialize the collection's locked NFTs for the given type if it doesn't exist
            if self.lockedNFTs[nftType] == nil {
                self.lockedNFTs[nftType] <-! {}
            }

            // Initialize the contract's locked tokens data for the given type if it doesn't exist
            if NFTLocker.lockedTokens[nftType] == nil {
                NFTLocker.lockedTokens[nftType] = {}
            }

            // Get a reference to this collection's locked NFTs map
            let collectionLockedNFTsRef = &self.lockedNFTs[nftType] as auth(Insert) &{UInt64: {NonFungibleToken.NFT}}?

            // Deposit the provided NFT in this collection's locked NFTs map - Cadence design requires destroying the resource-typed return value
            destroy <- collectionLockedNFTsRef!.insert(key: nftId, <- token)

            // Get a reference to the contract's nested map containing locked tokens data
            let lockedTokensDataRef = &NFTLocker.lockedTokens[nftType] as auth(Insert) &{UInt64: NFTLocker.LockedData}?
                ?? panic("Could not get a reference to the locked tokens data")

            // Create locked data
            let lockedData = NFTLocker.LockedData(
                id: nftId,
                owner: self.owner!.address,
                duration: duration,
                nftType: nftType
            )

            // Insert the locked data
            lockedTokensDataRef.insert(key: nftId, lockedData)

            // Increment the total locked tokens
            NFTLocker.totalLockedTokens = NFTLocker.totalLockedTokens + 1

            // Emit events
            emit NFTLocked(
                id: nftId,
                to: self.owner?.address,
                lockedAt: lockedData.lockedAt,
                lockedUntil: lockedData.lockedUntil,
                duration: lockedData.duration,
                nftType: nftType
            )

            emit Deposit(id: nftId, to: self.owner?.address)
        }

        /// Get the ids of NFTs locked for a given type
        ///
        access(all) view fun getIDs(nftType: Type): [UInt64]? {
            return self.lockedNFTs[nftType]?.keys
        }

        /// Initialize Collection resource
        ///
        view init() {
            self.lockedNFTs <- {}
        }
    }

    /// Create and return an empty collection
    ///
    access(all) fun createEmptyCollection(): @Collection {
        return <- create Collection()
    }

    /// Contract initializer
    ///
    init() {
        // Set paths
        self.CollectionStoragePath = /storage/NFTLockerCollection
        self.CollectionPublicPath = /public/NFTLockerCollection

        // Create and save the admin resource
        self.account.storage.save(<- create Admin(), to: NFTLocker.GetAdminStoragePath())

        // Set contract variables
        self.totalLockedTokens = 0
        self.lockedTokens = {}
    }
}