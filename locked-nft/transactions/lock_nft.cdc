import NonFungibleToken from 0xf8d6e0586b0a20c7
import LockedNFT from 0xf8d6e0586b0a20c7
import ExampleNFT from 0xf8d6e0586b0a20c7


transaction(nftID: UInt64) {
    let exampleCollectionRef: &ExampleNFT.Collection
    let signerAddress: Address

    prepare(signer: AuthAccount) {
        self.exampleCollectionRef = signer
            .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
        self.signerAddress = signer.address
    }

    execute {
            let lockRef = getAccount(self.signerAddress)
                .getCapability(LockedNFT.CollectionPublicPath).borrow<&{LockedNFT.LockedCollection}>()!

            lockRef.lock(token: <- self.exampleCollectionRef.withdraw(withdrawID: nftID), duration: 500)
    }
}