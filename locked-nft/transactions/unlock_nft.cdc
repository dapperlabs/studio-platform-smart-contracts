import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import LockedNFT from "../contracts/LockedNFT.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(nftID: UInt64) {
    let unlockRef: &LockedNFT.Collection
    let signerAddress: Address

    prepare(signer: AuthAccount) {
        self.unlockRef = signer
            .borrow<&LockedNFT.Collection>(from: LockedNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
        self.signerAddress = signer.address
    }

    execute {
        let depositRef = getAccount(self.signerAddress)
            .getCapability(ExampleNFT.CollectionPublicPath).borrow<&{NonFungibleToken.CollectionPublic}>()!

        depositRef.deposit(token: <- self.unlockRef.unlock(id: nftID))
    }
}