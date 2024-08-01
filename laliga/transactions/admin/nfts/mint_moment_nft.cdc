import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"

transaction(recipientAddress: Address, editionID: UInt64) {
    
    // local variable for storing the minter reference
    let minter: &Golazos.Admin
    let recipient: &Golazos.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // mint the NFT and deposit it to the recipient's collection
        let momentNFT <- self.minter.mintNFT(editionID: editionID)
        self.recipient.deposit(token: <- (momentNFT))
    }
}

