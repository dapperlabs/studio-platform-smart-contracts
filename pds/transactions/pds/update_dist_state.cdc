import PDS from "PDS"
import NonFungibleToken from "NonFungibleToken"

transaction (distId: UInt64, state: UInt8) {
    // state is an enum
    // - 0: Initialized
    // - 1: Invalid
    // - 2: Complete
    prepare(pds: auth(BorrowValue) &Account) {
        let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        cap.updateDistState(
            distId: distId,
            state: PDS.DistState(rawValue: state)!,
        )
    }
}
