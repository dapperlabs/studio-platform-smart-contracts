import Escrow from "Escrow"

// This script returns the leaderboard info for the given leaderboard name.
access(all) fun main(leaderboardName: String): Escrow.LeaderboardInfo? {
    let account = getAccount(0xf8d6e0586b0a20c7)

    let collectionPublic = account
        .capabilities.borrow<&Escrow.Collection>(Escrow.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the public leaderboard collection")

    return collectionPublic.getLeaderboardInfo(name: leaderboardName)
}
