import NonFungibleToken from "NonFungibleToken"
import EditionNFT from "EditionNFT"

transaction(recipientAddress: Address, editionID: UInt64) {
    
    // local variable for storing the minter reference
    let minter: auth(EditionNFT.Mint) &EditionNFT.Admin
    let recipient: &EditionNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<auth(EditionNFT.Mint) &EditionNFT.Admin>(from: EditionNFT.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.capabilities.borrow<&EditionNFT.Collection>(EditionNFT.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // mint the NFT and deposit it to the recipient's collection
        let nft <- self.minter.mintNFT(editionID: editionID)
        self.recipient.deposit(token: <- (nft as @{NonFungibleToken.NFT}))
    }
}

