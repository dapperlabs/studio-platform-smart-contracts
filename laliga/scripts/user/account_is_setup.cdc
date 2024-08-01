import Golazos from "Golazos"

// Check to see if an account looks like it has been set up to hold Golazos NFTs.

access(all) fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        != nil
}

