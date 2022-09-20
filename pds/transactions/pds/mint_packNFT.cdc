import PDS from "../../contracts/PDS.cdc"
import {{.PackNFTName}} from 0x{{.PackNFTAddress}}
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"

transaction (distId: UInt64, commitHashes: [String], issuer: Address ) {
    prepare(pds: AuthAccount) {
        let recvAcct = getAccount(issuer)
        let recv = recvAcct.getCapability({{.PackNFTName}}.CollectionPublicPath).borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Unable to borrow Collection Public reference for recipient")
        let cap = pds.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath) ?? panic("pds does not have Dist manager")
        cap.mintPackNFT(distId: distId, commitHashes: commitHashes, issuer: issuer, recvCap: recv)
    }
}
