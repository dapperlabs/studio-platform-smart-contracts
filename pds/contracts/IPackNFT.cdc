import Crypto
import NonFungibleToken from "NonFungibleToken"

/// Contract interface for PackNFT contracts.
///
access(all) contract interface IPackNFT{

    /// Entitlement to perform operations on the PackNFT
    ///
    access(all) entitlement Operatable

    /// StoragePath for Collection Resource
    ///
    access(all) let CollectionStoragePath: StoragePath

    /// PublicPath expected for deposit
    ///
    access(all) let CollectionPublicPath: PublicPath

    /// PublicPath for receiving PackNFT
    ///
    access(all) let CollectionIPackNFTPublicPath: PublicPath

    /// StoragePath for the PackNFT Operator Resource (issuer owns this)
    ///
    access(all) let OperatorStoragePath: StoragePath

    /// Request for Reveal
    ///
    access(all) event RevealRequest(id: UInt64, openRequest: Bool)

    /// Request for Open
    ///
    /// This is emitted when owner of a PackNFT request for the entitled NFT to be
    /// deposited to its account
    access(all) event OpenRequest(id: UInt64)

    /// Burned
    ///
    /// Emitted when a PackNFT has been burned
    access(all) event Burned(id: UInt64 )

    /// Opened
    ///
    /// Emitted when a packNFT has been opened
    access(all) event Opened(id: UInt64)

    // TODO: Clean up after enum handling/removal is clarified.
    // Enums cannot be declared anymore in interfaces in Cadence 1.0
    // access(all) enum Status: UInt8 {
    //     access(all) case Sealed
    //     access(all) case Revealed
    //     access(all) case Opened
    // }

    /// Struct interface for Collectible
    ///
    access(all) struct interface Collectible {
        access(all) let address: Address
        access(all) let contractName: String
        access(all) let id: UInt64
        access(all) fun hashString(): String
        init(address: Address, contractName: String, id: UInt64)
    }

    /// Resource interface for PackNFT
    ///
    access(all) resource interface IPack {
        access(all) let issuer: Address
        // access(all) var status: Status

        access(all) fun verify(nftString: String): Bool

        access(contract) fun reveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String)
        access(contract) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}])
        init(commitHash: String, issuer: Address)
    }

    /// Resource interface for IOperator
    ///
    access(all) resource interface IOperator {
        access(Operatable) fun mint(distId: UInt64, commitHash: String, issuer: Address): @{IPackNFT.NFT}
        access(Operatable) fun reveal(id: UInt64, nfts: [{Collectible}], salt: String)
        access(Operatable) fun open(id: UInt64, nfts: [{IPackNFT.Collectible}])
    }

    // Included for backwards compatibility
    access(all) resource interface PackNFTOperator: IOperator {}

    /// Resource interface for IPackNFTToken
    ///
    access(all) resource interface IPackNFTToken {
        access(all) let id: UInt64
        access(all) let issuer: Address
    }

    /// Resource interface for NFT
    ///
    access(all) resource interface NFT: NonFungibleToken.INFT, IPackNFTToken, IPackNFTOwnerOperator {
        access(all) let id: UInt64
        access(all) let issuer: Address
        access(NonFungibleToken.Update | NonFungibleToken.Owner) fun reveal(openRequest: Bool)
        access(NonFungibleToken.Update | NonFungibleToken.Owner) fun open()
    }

    // Included for backwards compatibility
    access(all) resource interface IPackNFTOwnerOperator{}
    access(all) resource interface IPackNFTCollectionPublic {}

    /// Emit a RevealRequest event to signal a Sealed Pack NFT should be revealed
    ///
    access(contract) fun revealRequest(id: UInt64, openRequest: Bool)

    /// Emit an OpenRequest event to signal a Revealed Pack NFT should be opened
    ///
    access(contract) fun openRequest(id: UInt64)

    /// Reveal a Sealed Pack NFT
    ///
    access(all) fun publicReveal(id: UInt64, nfts: [{IPackNFT.Collectible}], salt: String)
}