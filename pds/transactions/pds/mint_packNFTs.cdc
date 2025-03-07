import PDS from "PDS"
import PackNFT from "PackNFT"
import NonFungibleToken from "NonFungibleToken"
import DapperStorageRent from "DapperStorageRent"

transaction (distIds: [UInt64], commitHashes: [String], issuer: Address, receiver: Address) {
    prepare(pds: auth(BorrowValue) &Account) {
        DapperStorageRent.tryRefill(receiver)
        
        if distIds.length != commitHashes.length {
            panic("Number of distribution IDs must match number of commit hashes")
        }
        
        let recvAcct = getAccount(receiver)
        let recv = recvAcct.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath)
            ?? panic("Unable to borrow Collection Public reference for recipient")
            
        let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
            
        for i, distId in distIds {
            cap.mintPackNFT(distId: distId, commitHashes: [commitHashes[i]], issuer: issuer, recvCap: recv)
        }
    }
}