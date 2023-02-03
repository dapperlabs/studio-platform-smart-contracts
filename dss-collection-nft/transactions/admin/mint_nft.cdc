import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(
    recipientAddress: Address,
    collectionGroupID: UInt64,
    completionAddress: String,
    level: UInt8
) {
    let minter: &{DSSCollection.NFTMinter}
    let recipient: &{DSSCollection.DSSCollectionNFTCollectionPublic}

    prepare(signer: AuthAccount) {
        self.minter = signer.getCapability(DSSCollection.MinterPrivatePath)
            .borrow<&{DSSCollection.NFTMinter}>()
            ?? panic("Could not borrow a reference to the NFT minter")

        let recipientAccount = getAccount(recipientAddress)

        self.recipient = recipientAccount.getCapability(DSSCollection.CollectionPublicPath)
            .borrow<&{DSSCollection.DSSCollectionNFTCollectionPublic}>()
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        let nft <- self.minter.mintNFT(
            collectionGroupID: collectionGroupID,
            completionAddress: completionAddress,
            level: level
        )
        self.recipient.deposit(token: <- (nft as @NonFungibleToken.NFT))
    }
}