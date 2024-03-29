import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS

transaction(nftID: UInt64, to: Address) {
  /// Reference to the withdrawer's collection
  let withdrawRef: &ExampleNFT.Collection

  /// Reference of the collection to deposit the NFT to
  let depositRef: &{NonFungibleToken.CollectionPublic}

  prepare(signer: AuthAccount) {
    // borrow a reference to the signer's NFT collection
    self.withdrawRef = signer
        .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
        ?? panic("Account does not store an object at the specified path")

    // get the recipients public account object
    let recipient = getAccount(to)

    // borrow a public reference to the receivers collection
    self.depositRef = recipient
        .getCapability(ExampleNFT.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow a reference to the receiver's collection")

  }

  execute {
    // withdraw the NFT from the owner's collection
    let nft <- self.withdrawRef.withdraw(withdrawID: nftID)

    // Deposit the NFT in the recipient's collection
    self.depositRef.deposit(token: <-nft)
  }
}