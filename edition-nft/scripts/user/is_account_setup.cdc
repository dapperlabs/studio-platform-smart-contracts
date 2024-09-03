import NonFungibleToken from "NonFungibleToken"
import EditionNFT from "EditionNFT"

// Check to see if an account looks like it has been set up to hold EditionNFTs.

access(all) fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.capabilities.borrow<&EditionNFT.Collection>(EditionNFT.CollectionPublicPath)
        != nil
}

