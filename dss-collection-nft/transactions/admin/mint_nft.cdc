import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(recipientAddress: Address, collectionGroupID: UInt64, completedBy: String) {
    
    // local variable for storing the minter reference
    let minter: &{DSSCollection.NFTMinter}
    let recipient: &{DSSCollection.DSSCollectionNFTCollectionPublic}

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.getCapability(DSSCollection.MinterPrivatePath)
            .borrow<&{DSSCollection.NFTMinter}>()
            ?? panic("Could not borrow a reference to the NFT minter")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.getCapability(DSSCollection.CollectionPublicPath)
            .borrow<&{DSSCollection.DSSCollectionNFTCollectionPublic}>()
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // mint the NFT and deposit it to the recipient's collection
        let nft <- self.minter.mintNFT(collectionGroupID: collectionGroupID, completedBy: completedBy)
        self.recipient.deposit(token: <- (nft as @NonFungibleToken.NFT))
    }
}