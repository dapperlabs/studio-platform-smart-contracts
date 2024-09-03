import Golazos from "Golazos"

transaction(recipientAddress: Address, editionIDs: [UInt64], counts: [UInt64]) {
    
    // local variable for storing the minter reference
    let minter: auth(Golazos.Mint) &Golazos.Admin
    let recipient: &Golazos.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")

        // get the recipients public account object
        let recipientAccount = getAccount(recipientAddress)

        // borrow a public reference to the receivers collection
        self.recipient = recipientAccount.capabilities.borrow<auth(Golazos.Mint) &Golazos.Collection>(Golazos.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    pre {
        editionIDs.length == counts.length: "must pass arrays of same length"
    }

    execute {
        var i = 0
        while i < editionIDs.length {
            var remaining = counts[i]
            while remaining > 0 {
                // mint the NFT and deposit it to the recipient's collection
                self.recipient.deposit(token: <- self.minter.mintNFT(editionID: editionIDs[i]))
                remaining = remaining - 1
            }
            i = i + 1
        }
    }
}

