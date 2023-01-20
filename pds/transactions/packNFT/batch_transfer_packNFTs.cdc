import NonFungibleToken from 0x{{.NonFungibleTokenAddress}}
import {{.DapperSportContract}} from 0x{{.DapperSportAddress}}

// This transaction transfers NFL NFTs from one account to another.

transaction(recipientAddress: [Address], nftIDs: [UInt64]) {
    let collectionRef: &{{.DapperSportContract}}.Collection

    prepare(signer: AuthAccount) {
        // borrow a reference to the signer's NFT collection
        self.collectionRef = signer
            .borrow<&{{.DapperSportContract}}.Collection>(from: {{.DapperSportContract}}.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        var i = 0
        while i < nftIDs.length {
            // borrow a public reference to the receiver collection
            let depositRef = getAccount(recipientAddress[i])
                .getCapability({{.DapperSportContract}}.CollectionPublicPath).borrow<&{NonFungibleToken.CollectionPublic}>()!

            depositRef.deposit(token: <- self.collectionRef.withdraw(withdrawID: nftIDs[i]))
            i = i + 1
        }
    }
}
