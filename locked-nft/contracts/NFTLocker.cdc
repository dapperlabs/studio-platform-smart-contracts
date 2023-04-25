import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract NFTLocker {

    pub event ContractInitialized()
    pub event NFTLocked(
        id: UInt64,
        to: Address?,
        lockedAt: UInt64,
        lockedUntil: UInt64,
        duration: UInt64
    )
    pub event NFTUnlocked(
        id: UInt64,
        from: Address?
    )

    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath

    pub var totalLockedTokens:      UInt64
    access(self) let lockedTokens:  {String: LockedData}

    pub struct LockedData {
        pub let id: UInt64
        pub let owner: Address
        pub let lockedAt: UInt64
        pub let lockedUntil: UInt64
        pub let duration: UInt64
        pub let nftType: Type

        init (id: UInt64, owner: Address, duration: UInt64, nftType: Type) {
            let key = NFTLocker.getLockedTokenKey(id: id, nftType: nftType)
            if let lockedToken = NFTLocker.lockedTokens[key] {
                self.id = id
                self.owner = lockedToken.owner
                self.lockedAt = lockedToken.lockedAt
                self.lockedUntil = lockedToken.lockedUntil
                self.duration = lockedToken.duration
                self.nftType = lockedToken.nftType
            } else {
                self.id = id
                self.owner = owner
                self.lockedAt = UInt64(getCurrentBlock().timestamp)
                self.lockedUntil = self.lockedAt + duration
                self.duration = duration
                self.nftType = nftType
            }
        }
    }

    pub fun getNFTLockerDetails(id: UInt64, nftType: Type): NFTLocker.LockedData? {
        let key: String = NFTLocker.getLockedTokenKey(id: id, nftType: nftType)
        return NFTLocker.lockedTokens[key]
    }

    pub fun getLockedTokenKey(id: UInt64, nftType: Type): String {
        return nftType.identifier.concat(".").concat(id.toString())
    }

    pub fun canUnlockToken(id: UInt64, nftType: Type): Bool {
        let key: String = NFTLocker.getLockedTokenKey(id: id, nftType: nftType)
        if let lockedToken = NFTLocker.lockedTokens[key] {
            if lockedToken.lockedUntil < UInt64(getCurrentBlock().timestamp) {
                return true
            }
        }

        return false
    }

    pub resource interface LockedCollection {
        pub fun getIDs(): [String]
        pub fun borrowNFT(key: String): &NonFungibleToken.NFT
    }

    pub resource interface LockProvider {
        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64)
        pub fun unlock(id: UInt64, nftType: Type): @NonFungibleToken.NFT
    }

    pub resource Collection: LockedCollection {
        pub var lockedNFTs: @{String: NonFungibleToken.NFT}

        pub fun unlock(id: UInt64, nftType: Type): @NonFungibleToken.NFT {
            pre {
                NFTLocker.canUnlockToken(
                    id: id,
                    nftType: nftType
                ) == true : "locked duration has not been met"
            }

            let key: String = NFTLocker.getLockedTokenKey(id: id, nftType: nftType)
            let token <- self.lockedNFTs.remove(key: key) ?? panic("Missing NFT")
            NFTLocker.lockedTokens.remove(key: key)
            NFTLocker.totalLockedTokens = NFTLocker.totalLockedTokens - 1

            emit NFTUnlocked(
                id: token.id,
                from: self.owner?.address
            )

            return <-token
        }

        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64) {
            let id: UInt64 = token.id
            let nftType: Type = token.getType()
            let key: String = NFTLocker.getLockedTokenKey(id: id, nftType: nftType)
            let oldToken <- self.lockedNFTs[key] <- token
            let lockedData = NFTLocker.LockedData(
                id: id,
                owner: self.owner!.address,
                duration: duration,
                nftType: nftType
            )
            emit NFTLocked(
                id: id,
                to: self.owner?.address,
                lockedAt: lockedData.lockedAt,
                lockedUntil: lockedData.lockedUntil,
                duration: lockedData.duration
            )
            NFTLocker.lockedTokens[key] = lockedData
            NFTLocker.totalLockedTokens = NFTLocker.totalLockedTokens + 1

            destroy oldToken
        }

        pub fun getIDs(): [String] {
            return self.lockedNFTs.keys
        }

        pub fun borrowNFT(key: String): &NonFungibleToken.NFT {
            return (&self.lockedNFTs[key] as &NonFungibleToken.NFT?)!
        }

        destroy() {
            destroy self.lockedNFTs
        }

        init() {
            self.lockedNFTs <- {}
        }
    }

    pub fun createEmptyCollection(): @Collection {
        return <- create Collection()
    }

    init() {
        self.CollectionStoragePath = /storage/NFTLockerCollection
        self.CollectionPublicPath = /public/NFTLockerCollection

        self.totalLockedTokens = 0
        self.lockedTokens = {}

        emit ContractInitialized()
    }
}