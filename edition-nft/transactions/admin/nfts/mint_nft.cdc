import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import EditionNFT from "../../../contracts/EditionNFT.cdc"

transaction(recipientAddress: Address, editionID: UInt64) {
    
    // local variable for storing the minter reference
    let minter: &{EditionNFT.NFTMinter}
    let recipient: &{EditionNFT.EditionNFTCollectionPublic}

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.getCapability(EditionNFT.MinterPrivatePath)
            .borrow<&{EditionNFT.NFTMinter}>()
            ?? panic("Could not borrow a reference to the NFT minter")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.getCapability(EditionNFT.CollectionPublicPath)
            .borrow<&{EditionNFT.EditionNFTCollectionPublic}>()
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // mint the NFT and deposit it to the recipient's collection
        let nft <- self.minter.mintNFT(editionID: editionID)
        self.recipient.deposit(token: <- (nft as @NonFungibleToken.NFT))
    }
}

