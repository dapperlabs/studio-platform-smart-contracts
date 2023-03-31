import NonFungibleToken from "./NonFungibleToken.cdc"
import EnglishPremierLeague from "./EnglishPremierLeague.cdc"

transaction(recipientAddress: Address, editionID: UInt64) {

    let minter: &{EnglishPremierLeague.NFTMinter}
    let recipient: &{EnglishPremierLeague.MomentNFTCollectionPublic}

    prepare(signer: AuthAccount) {
        self.minter = signer.getCapability(EnglishPremierLeague.MinterPrivatePath)
            .borrow<&{EnglishPremierLeague.NFTMinter}>()
            ?? panic("Could not borrow a reference to the NFT minter")

        let recipientAccount = getAccount(recipientAddress)

        self.recipient = recipientAccount.getCapability(EnglishPremierLeague.CollectionPublicPath)
            .borrow<&{EnglishPremierLeague.MomentNFTCollectionPublic}>()
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        let nft <- self.minter.mintNFT(
            editionID: editionID,
            ext: nil
        )
        self.recipient.deposit(token: <- (nft as @NonFungibleToken.NFT))
    }
}