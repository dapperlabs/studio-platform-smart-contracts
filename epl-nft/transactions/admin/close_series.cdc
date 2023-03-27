import EnglishPremierLeague from "../../contracts/EnglishPremierLeague.cdc"

transaction(seriesID: UInt64) {
    let admin: &EnglishPremierLeague.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&EnglishPremierLeague.Admin>(from: EnglishPremierLeague.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EnglishPremierLeague Admin capability")
    }

    execute {
        self.admin.closeSeries(
            id: seriesID
        )
    }
}