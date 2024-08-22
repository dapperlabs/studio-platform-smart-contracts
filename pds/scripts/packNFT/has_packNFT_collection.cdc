import PackNFT from "PackNFT"

/// Check if an account has been set up to hold Pinnacle NFTs.
///
access(all) fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.capabilities.borrow<
        &PackNFT.Collection>(PackNFT.CollectionPublicPath) != nil
}
