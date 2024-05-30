import PDS from "PDS"
import "PackNFTName" from "PackNFT"
import NonFungibleToken from "NonFungibleToken"

transaction (distId: UInt64, commitHashes: [String], issuer: Address ) {
    prepare(pds: auth(BorrowValue) &Account) {
        let recvAcct = getAccount(issuer)
        let recv = recvAcct.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>("PackNFTName".CollectionPublicPath)
            ?? panic("Unable to borrow Collection Public reference for recipient")
        let cap = pds.storage.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        cap.mintPackNFT(distId: distId, commitHashes: commitHashes, issuer: issuer, recvCap: recv)
    }
}
