import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS
import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(address: Address, collectionGroupID: UInt64): Bool {
    let completedNFTs: [DSSCollection.CollectionCompletedWith]? = DSSCollection.getCompletedCollectionIDs(address: address)

    if completedNFTs == nil {
        return false
    }
    let account = getAccount(address)

    let userNFTs = account
        .getCapability(ExampleNFT.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")

    // Iterate over the completed NFTs and check if the user still
    // owns them and they are of type ExampleNFT.NFT
    for completedCollectionWith in completedNFTs! {
        if completedCollectionWith.collectionGroupID == collectionGroupID {
            for nftID in completedCollectionWith.nftIDs {
                // Use getIDs to check if the nftID exists in the collection
                let ownedNFTs = userNFTs.getIDs()

                if !(ownedNFTs.contains(nftID)) {
                    return false
                }
            }
        }
    }

    return true
}
