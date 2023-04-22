import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract LockedNFT {

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
        pub let nftType: String

        init (id: UInt64, owner: Address, duration: UInt64, nftType: String) {
            let key = LockedNFT.getLockedTokenKey(id: id, nftType: nftType)
            if let lockedToken = LockedNFT.lockedTokens[key] {
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

    pub fun getLockedToken(key: String): LockedNFT.LockedData? {
        return LockedNFT.lockedTokens[key]
    }

    pub fun getLockedTokenKey(id: UInt64, nftType: String): String {
        return id.toString().concat("-").concat(nftType)
    }

    pub fun canUnlockToken(key: String): Bool {
        if let lockedToken = LockedNFT.lockedTokens[key] {
            if lockedToken.lockedUntil < UInt64(getCurrentBlock().timestamp) {
                return true
            }
        }

        return false
    }

    pub resource interface LockedCollection {
        pub fun getIDs(): [String]
        pub fun borrowNFT(key: String): &NonFungibleToken.NFT
        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64)
        pub fun unlock(key: String): @NonFungibleToken.NFT
    }

    pub resource Collection: LockedCollection {
        pub var lockedNFTs: @{String: NonFungibleToken.NFT}

        pub fun unlock(key: String): @NonFungibleToken.NFT {
            pre {
                LockedNFT.canUnlockToken(
                    key: key
                ) == true : "locked duration has not been met"
            }

            let lockedData = LockedNFT.getLockedToken(key: key)
            let token <- self.lockedNFTs.remove(key: key) ?? panic("Missing NFT")
            LockedNFT.lockedTokens.remove(key: key)
            LockedNFT.totalLockedTokens = LockedNFT.totalLockedTokens - 1

            emit NFTUnlocked(
                id: token.id,
                from: self.owner?.address
            )

            return <-token
        }

        pub fun lock(token: @NonFungibleToken.NFT, duration: UInt64) {
            let id: UInt64 = token.id
            let nftType: String = token.getType().identifier
            let key: String = LockedNFT.getLockedTokenKey(id: id, nftType: nftType)
            let oldToken <- self.lockedNFTs[key] <- token
            let lockedData = LockedNFT.LockedData(
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
            LockedNFT.lockedTokens[key] = lockedData
            LockedNFT.totalLockedTokens = LockedNFT.totalLockedTokens + 1

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
        self.CollectionStoragePath = /storage/LockedNFTCollection
        self.CollectionPublicPath = /public/LockedNFTCollection

        self.totalLockedTokens = 0
        self.lockedTokens = {}

        emit ContractInitialized()
    }
}