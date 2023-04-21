import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

transaction(recipient: Address) {
    let minter: &ExampleNFT.NFTMinter
    let recipientCollectionRef: &{NonFungibleToken.CollectionPublic}
    let mintingIDBefore: UInt64

    prepare(signer: AuthAccount) {
        self.mintingIDBefore = ExampleNFT.totalSupply

        self.minter = signer.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)
            ?? panic("Account does not store an object at the specified path")

        self.recipientCollectionRef = getAccount(recipient)
            .getCapability(ExampleNFT.CollectionPublicPath)
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")
    }

    execute {
        // Stub data for other parameters
        let name: String = "Example NFT"
        let description: String = "This is an example NFT."
        let thumbnail: String = "nft.jpg"

        self.minter.mintNFT(
            recipient: self.recipientCollectionRef,
            name: name,
            description: description,
            thumbnail: thumbnail,
            royalties: []
        )
    }

    post {
        self.recipientCollectionRef.getIDs().contains(self.mintingIDBefore): "The next NFT ID should have been minted and delivered"
        ExampleNFT.totalSupply == self.mintingIDBefore + 1: "The total supply should have been increased by 1"
    }
}