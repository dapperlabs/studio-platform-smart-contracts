import NonFungibleToken from "NonFungibleToken"
import EditionNFT from "EditionNFT"

// This transaction transfers a EditionNFT NFT from one account to another.

transaction(recipientAddress: Address, withdrawID: UInt64) {
    prepare(signer: auth(BorrowValue) &Account) {
        
        // get the recipients public account object
        let recipient = getAccount(recipientAddress)

        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw)
            &{NonFungibleToken.Collection}>(from: EditionNFT.CollectionStoragePath)
                ?? panic("Could not borrow a reference to the owner's collection")

        // borrow a public reference to the receivers collection
        let depositRef = recipient.capabilities.borrow<&EditionNFT.Collection>(EditionNFT.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the collection receiver")

        // withdraw the NFT from the owner's collection
        let nft: @{NonFungibleToken.NFT} <- collectionRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        depositRef.deposit(token: <-nft)
    }
}
