import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract LockedNFT {

    pub event ContractInitialized()
    pub event NFTLocked(id: UInt64, to: Address?)
    pub event NFTUnlocked(id: UInt64, from: Address?)

    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath

    pub var totalLockedTokens:      UInt64
    access(self) let lockedTokens:  {UInt64: LockedData}

    pub struct LockedData {
        pub let owner: Address
        pub let lockedAt: UInt64
        pub let lockedUntil: UInt64
        pub let nftType: Type

        init (id: UInt64, owner: Address, lockedAt: UInt64, duration: UInt64, nftType: Type) {
            if let lockedToken = LockedNFT.lockedTokens[id] {
                self.owner = lockedToken.owner
                self.lockedAt = lockedToken.lockedAt
                self.lockedUntil = lockedToken.lockedUntil
                self.nftType = lockedToken.nftType
            } else {
                self.owner = owner
                self.lockedAt = UInt64(getCurrentBlock().timestamp)
                self.lockedUntil = self.lockedAt + duration
                self.nftType = nftType
            }
        }
    }

    pub fun canUnlockToken(id: UInt64): Bool {
        if let lockedToken = LockedNFT.lockedTokens[id] {
            if lockedToken.lockedUntil < UInt64(getCurrentBlock().timestamp) {
                return true
            }
        }

        return false
    }

    pub resource interface LockedCollection {
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64)
        pub fun unlock(id: UInt64): @NonFungibleToken.NFT
    }

    pub resource Collection: LockedCollection {
        pub var lockedNFTs: @{UInt64: NonFungibleToken.NFT}

        pub fun unlock(id: UInt64): @NonFungibleToken.NFT {
            pre {
                LockedNFT.canUnlockToken(
                    id: id
                ) == true : "locked duration has not been met"
            }

            let token <- self.lockedNFTs.remove(key: id) ?? panic("Missing NFT")
            LockedNFT.lockedTokens.remove(key: id)
            LockedNFT.totalLockedTokens = LockedNFT.totalLockedTokens - 1
            emit NFTUnlocked(id: token.id, from: self.owner?.address)

            return <-token
        }

        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64) {
            let id: UInt64 = token.id
            let oldToken <- self.lockedNFTs[id] <- token
            emit NFTLocked(id: id, to: self.owner?.address)
            let lockedData = LockedNFT.LockedData(
                id: id,
                owner: self.owner!.address,
                lockedAt: UInt64(getCurrentBlock().timestamp),
                lockedUntil: UInt64(getCurrentBlock().timestamp) + duration,
                nftType: oldToken.getType()
            )
            LockedNFT.lockedTokens[id] = lockedData
            LockedNFT.totalLockedTokens = LockedNFT.totalLockedTokens + 1

            destroy oldToken
        }

        pub fun getIDs(): [UInt64] {
            return self.lockedNFTs.keys
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.lockedNFTs[id] as &NonFungibleToken.NFT?)!
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
        self.CollectionStoragePath = /storage/LockedNFTCollection
        self.CollectionPublicPath = /public/LockedNFTCollection

        self.totalLockedTokens = 0
        self.lockedTokens = {}

        emit ContractInitialized()
    }
}