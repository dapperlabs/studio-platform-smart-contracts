import NFTLocker from "../contracts/NFTLocker.cdc"

transaction() {

    prepare(signer: AuthAccount) {
        NFTLocker.createAndSaveAdmin(acct: signer)
    }

    execute {
        log("NFTLocker Admin Account Created and Stored in the NFTLocker Admin Storage.")
    }
}