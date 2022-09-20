import PackNFT from "../../contracts/PackNFT.cdc"
import PDS from "../../contracts/PDS.cdc"
import IPackNFT from "../../contracts/IPackNFT.cdc"
import ExampleNFT from "../../contracts/ExampleNFT.cdc"

transaction (packId: UInt64, nftContractAddrs: [Address], nftContractName: [String], nftIds: [UInt64], salt: String) {
    prepare(pds: AuthAccount) {
        assert(
            nftContractAddrs.length == nftContractName.length && 
            nftContractName.length == nftIds.length, 
            message: "NFTs must be fully described"
        )
        let arr: [{IPackNFT.Collectible}] = []
        var i = 0
        while i < nftContractAddrs.length {
            let s = PDS.Collectible(address: nftContractAddrs[i], contractName: nftContractName[i], id: nftIds[i])
            arr.append(s)
            i = i + 1
        }

        PackNFT.publicReveal(
            id: packId, 
            nfts: arr, 
            salt: salt)
    }
}

