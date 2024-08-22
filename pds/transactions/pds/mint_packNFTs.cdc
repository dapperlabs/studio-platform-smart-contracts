import PDS from "PDS"
import PackNFT from "PackNFT"
import NonFungibleToken from "NonFungibleToken"
import DapperStorageRent from "DapperStorageAddress"

transaction (distIds: [UInt64], commitHashes: [String], issuer: Address, receiver: Address) {
    prepare(pds: auth(BorrowValue) &Account) {
        DapperStorageRent.tryRefill(receiver)
        for i, distId in distIds {
            let recvAcct = getAccount(receiver)
            let recv = recvAcct.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath)
               ?? panic("Unable to borrow Collection Public reference for recipient")
            let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
               ?? panic("pds does not have Dist manager")
            cap.mintPackNFT(distId: distId, commitHashes: commitHashes, issuer: issuer, recvCap: recv)
        }
    }
}