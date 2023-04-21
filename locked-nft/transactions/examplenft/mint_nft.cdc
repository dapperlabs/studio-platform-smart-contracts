import NonFungibleToken from 0xf8d6e0586b0a20c7
import ExampleNFT from 0xf8d6e0586b0a20c7

transaction(targetAcct: Address) {
    prepare(signer: AuthAccount) {
		let acct = getAccount(targetAcct)
        let minter = signer.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)
            ?? panic("Account does not store an object at the specified path")
		//let collectionRef = acct.getCapability(ExampleNFT.CollectionPublicPath)!.borrow<&{ExampleNFT.ExampleNFTCollectionPublic}>()!
        let collectionRef = acct
            .getCapability(ExampleNFT.CollectionPublicPath)
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")
		minter.mintNFT(recipient: collectionRef, name: "derp", description: "derp", thumbnail: "derp", royalties: [])
    }
}
