import EnglishPremierLeague from "../../contracts/EnglishPremierLeague.cdc"

transaction(seriesID: UInt64, setID: UInt64, playID: UInt64, maxMintSize: UInt64?, tier: String) {
    let admin: &EnglishPremierLeague.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&EnglishPremierLeague.Admin>(from: EnglishPremierLeague.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EnglishPremierLeague Admin capability")
    }

    execute {
        let id = self.admin.createEdition(
            seriesID: seriesID,
            setID: setID,
            playID: playID,
            maxMintSize: maxMintSize,
            tier: tier
        )
    }
}