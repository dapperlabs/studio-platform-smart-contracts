import EnglishPremierLeague from 0xf8d6e0586b0a20c7

transaction(name: String) {
    let admin: &EnglishPremierLeague.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&EnglishPremierLeague.Admin>(from: EnglishPremierLeague.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EnglishPremierLeague Admin capability")
    }

    execute {
        let id = self.admin.createTag(
            name: name,
        )
    }
}