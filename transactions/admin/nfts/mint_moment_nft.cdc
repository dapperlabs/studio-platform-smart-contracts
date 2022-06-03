import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import Golazo from "../../../contracts/Golazo.cdc"

transaction(recipientAddress: Address, editionID: UInt64) {
    
    // local variable for storing the minter reference
    let minter: &{Golazo.NFTMinter}
    let recipient: &{Golazo.MomentNFTCollectionPublic}

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.getCapability(Golazo.MinterPrivatePath)
            .borrow<&{Golazo.NFTMinter}>()
            ?? panic("Could not borrow a reference to the NFT minter")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.getCapability(Golazo.CollectionPublicPath)
            .borrow<&{Golazo.MomentNFTCollectionPublic}>()
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // mint the NFT and deposit it to the recipient's collection
        let momentNFT <- self.minter.mintNFT(editionID: editionID)
        self.recipient.deposit(token: <- (momentNFT as @NonFungibleToken.NFT))
    }
}

