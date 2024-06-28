import NFTLocker from "NFTLocker"

transaction() {

    prepare(signer: &Account) {
        NFTLocker.createAndSaveAdmin(acct: signer)
    }

    execute {
        log("NFTLocker Admin Account Created and Stored in the NFTLocker Admin Storage.")
    }
}