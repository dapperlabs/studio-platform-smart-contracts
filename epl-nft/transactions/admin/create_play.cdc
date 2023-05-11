import EnglishPremierLeague from "./EnglishPremierLeague.cdc"

transaction(classification: String, metadata: {String: String}, tagIds: [UInt64]) {
    let admin: &EnglishPremierLeague.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&EnglishPremierLeague.Admin>(from: EnglishPremierLeague.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EnglishPremierLeague Admin capability")
    }

    execute {
        let id = self.admin.createPlay(
            classification: classification,
            metadata: metadata,
            tagIds: tagIds
        )
    }
}