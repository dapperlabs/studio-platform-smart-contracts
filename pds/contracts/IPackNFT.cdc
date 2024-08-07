import Crypto
import NonFungibleToken from 0x{{.NonFungibleToken}}

pub contract interface IPackNFT{
    /// StoragePath for Collection Resource
    ///
    pub let CollectionStoragePath: StoragePath
    /// PublicPath expected for deposit
    ///
    pub let CollectionPublicPath: PublicPath
    /// PublicPath for receiving PackNFT
    ///
    pub let CollectionIPackNFTPublicPath: PublicPath
    /// StoragePath for the PackNFT Operator Resource (issuer owns this)
    ///
    pub let OperatorStoragePath: StoragePath
    /// PrivatePath to share IOperator interfaces with Operator (typically with PDS account)
    ///
    pub let OperatorPrivPath: PrivatePath
    /// Request for Reveal
    ///
    pub event RevealRequest(id: UInt64, openRequest: Bool)
    /// Request for Open
    ///
    /// This is emitted when owner of a PackNFT request for the entitled NFT to be
    /// deposited to its account
    pub event OpenRequest(id: UInt64)
    /// Burned
    ///
    /// Emitted when a PackNFT has been burned
    pub event Burned(id: UInt64 )
    /// Opened
    ///
    /// Emitted when a packNFT has been opened
    pub event Opened(id: UInt64)

    pub enum Status: UInt8 {
        pub case Sealed
        pub case Revealed
        pub case Opened
    }

    pub struct interface Collectible {
        pub let address: Address
        pub let contractName: String
        pub let id: UInt64
        pub fun hashString(): String
        init(address: Address, contractName: String, id: UInt64)
    }

    pub resource interface IPack {
        pub let issuer: Address
        pub var status: Status

        pub fun verify(nftString: String): Bool

        access(contract) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String)
        access(contract) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}])
        init(commitHash: String, issuer: Address)
    }

    pub resource interface IOperator {
        pub fun mint(distId: UInt64, commitHash: String, issuer: Address): @NFT
        pub fun reveal(id: UInt64, nfts: [{Collectible}], salt: String)
        pub fun open(id: UInt64, nfts: [{IPackNFT.Collectible}])
    }
    pub resource PackNFTOperator: IOperator {
        pub fun mint(distId: UInt64, commitHash: String, issuer: Address): @NFT
        pub fun reveal(id: UInt64, nfts: [{Collectible}], salt: String)
        pub fun open(id: UInt64, nfts: [{IPackNFT.Collectible}])
    }

    pub resource interface IPackNFTToken {
        pub let id: UInt64
        pub let issuer: Address
    }

    pub resource NFT: NonFungibleToken.INFT, IPackNFTToken, IPackNFTOwnerOperator{
        pub let id: UInt64
        pub let issuer: Address
        pub fun reveal(openRequest: Bool)
        pub fun open()
    }

    pub resource interface IPackNFTOwnerOperator{
        pub fun reveal(openRequest: Bool)
        pub fun open()
    }

    pub resource interface IPackNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
    }

    access(contract) fun revealRequest(id: UInt64, openRequest: Bool)
    access(contract) fun openRequest(id: UInt64)
    pub fun publicReveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String)
}