import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import EditionNFT from "../../contracts/EditionNFT.cdc"

// This transaction transfers a EditionNFT NFT from one account to another.

transaction(recipientAddress: Address, withdrawID: UInt64) {
    prepare(signer: AuthAccount) {
        
        // get the recipients public account object
        let recipient = getAccount(recipientAddress)

        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.borrow<&EditionNFT.Collection>(from: EditionNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // borrow a public reference to the receivers collection
        let depositRef = recipient.getCapability(EditionNFT.CollectionPublicPath).borrow<&{NonFungibleToken.CollectionPublic}>()!

        // withdraw the NFT from the owner's collection
        let nft <- collectionRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        depositRef.deposit(token: <-nft)
    }
}
