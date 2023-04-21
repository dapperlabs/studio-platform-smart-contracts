import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import LockedNFT from "../contracts/LockedNFT.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(nftID: UInt64, duration: UInt64) {
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

            lockRef.lock(token: <- self.exampleCollectionRef.withdraw(withdrawID: nftID), duration: duration)
    }
}