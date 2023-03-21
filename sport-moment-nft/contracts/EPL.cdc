/*
    Contract for EPL sports moments & metadata
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import FungibleToken from 0x{{.FungibleTokenAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenAddress}}
import MetadataViews from 0x{{.MetadataViewsAddress}}

/// The EPL NFT and metadata contract
//
pub contract EPL: NonFungibleToken {
    // -----------------------------------------------------------------------
    // EPL deployment variables
    // -----------------------------------------------------------------------

    pub fun RoyaltyAddress() : Address { return 0xf8d6e0586b0a20c7 }
    //------------------------------------------------------------
    // Events
    //------------------------------------------------------------

    // Contract Events
    //
    pub event ContractInitialized()

    // NFT Collection Events
    //
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)

    // Series Events
    //
    /// Emitted when a new series has been created by an admin
    pub event SeriesCreated(id: UInt64, name: String)
    /// Emitted when a series is closed by an admin
    pub event SeriesClosed(id: UInt64)

    // Set Events
    //
    /// Emitted when a new set has been created by an admin
    pub event SetCreated(id: UInt64, name: String)

    /// Emitted when a Set is locked, meaning Editions cannot be created with the set
    pub event SetLocked(setID: UInt64)

    // Play Events
    //
    /// Emitted when a new play has been created by an admin
    pub event PlayCreated(id: UInt64, metadata: {String: String}, tagIds: [UInt64])

    // Tag Events
    //
    /// Emitted when a new tag has been created by an admin
    pub event TagCreated(id: UInt64, name: String)

    // Edition Events
    //
    /// Emitted when a new edition has been created by an admin
    pub event EditionCreated(
        id: UInt64,
        seriesID: UInt64,
        setID: UInt64,
        playID: UInt64,
        maxMintSize: UInt64?,
        tier: String,
    )
    /// Emitted when an edition is either closed by an admin, or the max amount of moments have been minted
    pub event EditionClosed(id: UInt64)

    // NFT Events
    //
    /// Emitted when a moment nft is minted
    pub event MomentNFTMinted(id: UInt64, editionID: UInt64, serialNumber: UInt64)
    /// Emitted when a moment nft resource is destroyed
    pub event MomentNFTBurned(id: UInt64,  editionID: UInt64, serialNumber: UInt64)

    //------------------------------------------------------------
    // Named values
    //------------------------------------------------------------

    /// Named Paths
    ///
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let AdminStoragePath:       StoragePath
    pub let MinterPrivatePath:      PrivatePath

    //------------------------------------------------------------
    // Publicly readable contract state
    //------------------------------------------------------------

    /// Entity Counts
    ///
    pub var totalSupply:        UInt64
    pub var nextSeriesID:       UInt64
    pub var nextSetID:          UInt64
    pub var nextTagID:          UInt64
    pub var nextPlayID:         UInt64
    pub var nextEditionID:      UInt64

    //------------------------------------------------------------
    // Internal contract state
    //------------------------------------------------------------

    /// Metadata Dictionaries
    ///
    /// This is so we can find Series by their names (via seriesByID)
    access(self) let seriesIDByName:    {String: UInt64}
    access(self) let seriesByID:        @{UInt64: Series}
    access(self) let setIDByName:       {String: UInt64}
    access(self) let setByID:           @{UInt64: Set}
    access(self) let tagByID:           @{UInt64: Tag}
    access(self) let playByID:          @{UInt64: Play}
    access(self) let editionByID:       @{UInt64: Edition}

    //------------------------------------------------------------
    // Series
    //------------------------------------------------------------

    /// A public struct to access Series data
    ///
    pub struct SeriesData {
        pub let id: UInt64
        pub let name: String
        pub let active: Bool

        /// initializer
        //
        init (id: UInt64) {
            let series = (&EPL.seriesByID[id] as! &EPL.Series?)!
            self.id = series.id
            self.name = series.name
            self.active = series.active
        }
    }

    /// A top-level Series with a unique ID and name
    ///
    pub resource Series {
        pub let id: UInt64
        pub let name: String
        pub var active: Bool

        /// Close this series
        ///
        pub fun close() {
            pre {
                self.active == true: "series is not active"
            }

            self.active = false

            emit SeriesClosed(id: self.id)
        }

        /// initializer
        ///
        init (name: String) {
            pre {
                !EPL.seriesIDByName.containsKey(name): "A Series with that name already exists"
            }
            self.id = EPL.nextSeriesID
            self.name = name
            self.active = true

            // Cache the new series's name => ID
            EPL.seriesIDByName[name] = self.id
            // Increment for the nextSeriesID
            EPL.nextSeriesID = self.id + 1 as UInt64

            emit SeriesCreated(id: self.id, name: self.name)
        }
    }

    /// Get the publicly available data for a Series by id
    ///
    pub fun getSeriesData(id: UInt64): EPL.SeriesData {
        pre {
            EPL.seriesByID[id] != nil: "Cannot borrow series, no such id"
        }

        return EPL.SeriesData(id: id)
    }

    /// Get the publicly available data for a Series by name
    ///
    pub fun getSeriesDataByName(name: String): EPL.SeriesData? {
        let id = EPL.seriesIDByName[name]

        if id == nil{
            return nil
        }

        return EPL.SeriesData(id: id!)
    }

    /// Get all series names (this will be *long*)
    ///
    pub fun getAllSeriesNames(): [String] {
        return EPL.seriesIDByName.keys
    }

    /// Get series id by name
    ///
    pub fun getSeriesIDByName(name: String): UInt64? {
        return EPL.seriesIDByName[name]
    }

    //------------------------------------------------------------
    // Set
    //------------------------------------------------------------

    /// A public struct to access Set data
    ///
    pub struct SetData {
        pub let id: UInt64
        pub let name: String
        pub let locked: Bool

        /// initializer
        ///
        init (id: UInt64) {
            let set = (&EPL.setByID[id] as! &EPL.Set?)!
            self.id = id
            self.name = set.name
            self.locked = set.locked
        }
    }

    /// A top level Set with a unique ID and a name
    ///
    pub resource Set {
        pub let id: UInt64
        pub let name: String

        // Indicates if the Set is currently locked.
        // When a Set is created, it is unlocked
        // and Editions can be created with it.
        // When a Set is locked, new Editions cannot be created with the Set.
        // A Set can never be changed from locked to unlocked,
        // the decision to lock a Set is final.
        // If a Set is locked, Moments can still be minted from the
        // Editions already created from the Set.
        pub var locked: Bool

        /// initializer
        ///
        init (name: String) {
            pre {
                !EPL.setIDByName.containsKey(name): "A Set with that name already exists"
            }
            self.id = EPL.nextSetID
            self.name = name
            self.locked = false

            // Cache the new set's name => ID
            EPL.setIDByName[name] = self.id
            // Increment for the nextSeriesID
            EPL.nextSetID = self.id + 1 as UInt64

            emit SetCreated(id: self.id, name: self.name)
        }

        // lock() locks the Set so that no more Plays can be added to it
        //
        // Pre-Conditions:
        // The Set should not be locked
        pub fun lock() {
            if !self.locked {
                self.locked = true
                emit SetLocked(setID: self.id)
            }
        }
    }

    /// Get the publicly available data for a Set
    ///
    pub fun getSetData(id: UInt64): EPL.SetData? {
        if EPL.setByID[id] == nil {
            return nil
        }
        return EPL.SetData(id: id!)
    }

    /// Get the publicly available data for a Set by name
    ///
    pub fun getSetDataByName(name: String): EPL.SetData? {
        let id = EPL.setIDByName[name]

        if id == nil {
            return nil
        }
        return EPL.SetData(id: id!)
    }

    /// Get all set names (this will be *long*)
    ///
    pub fun getAllSetNames(): [String] {
        return EPL.setIDByName.keys
    }

    //------------------------------------------------------------
    // Tag
    //------------------------------------------------------------
    pub struct TagData {
        pub let id: UInt64
        pub let name: String

        /// initializer
        ///
        init (id: UInt64) {
            let tag = (&EPL.tagByID[id] as! &EPL.Tag?)!
            self.id = id
            self.name = tag.name
        }
    }

    /// A top level Tag with a unique ID
    //
    pub resource Tag {
        pub let id: UInt64
        pub let name: String

        /// initializer
        ///
        init (name: String) {
            self.id = EPL.nextTagID
            self.name = name

            EPL.nextTagID = self.id + 1 as UInt64

            emit TagCreated(id: self.id, name: self.name)
        }
    }

    /// Get the publicly available data for a Tag
    ///
    pub fun getTagData(id: UInt64): EPL.TagData? {
        if EPL.tagByID[id] == nil {
            return nil
        }

        return EPL.TagData(id: id!)
    }

    //------------------------------------------------------------
    // Play
    //------------------------------------------------------------

    /// A public struct to access Play data
    ///
    pub struct PlayData {
        pub let id: UInt64
        pub let metadata: {String: String}
        pub var tagIds: [UInt64]

        /// initializer
        ///
        init (id: UInt64) {
            let play = (&EPL.playByID[id] as! &EPL.Play?)!
            self.id = id
            self.metadata = play.getMetadata()
            self.tagIds = play.getTagIds()
        }
    }

    /// A top level Play with a unique ID
    //
    pub resource Play {
        pub let id: UInt64
        access(self) let metadata: {String: String}
        pub var tagIds: [UInt64]

        /// returns the metadata set for this play
        pub fun getMetadata(): {String:String} {
            return self.metadata
        }

        /// returns the tagIds for this play
        pub fun getTagIds(): [UInt64] {
            return self.tagIds
        }

        /// initializer
        ///
        init (metadata: {String: String}, tagIds: [UInt64]) {
            pre {
                EPL.validateTags(
                    tagIds: tagIds
                ) == true : "Play contains tag that does not exist."
            }
            self.id = EPL.nextPlayID
            self.metadata = metadata
            self.tagIds = tagIds

            EPL.nextPlayID = self.id + 1 as UInt64

            emit PlayCreated(id: self.id, metadata: self.metadata, tagIds: self.tagIds)
        }
    }

    /// Get the publicly available data for a Play
    ///
    pub fun getPlayData(id: UInt64): EPL.PlayData? {
        if EPL.playByID[id] == nil {
            return nil
        }

        return EPL.PlayData(id: id!)
    }

    //------------------------------------------------------------
    // Edition
    //------------------------------------------------------------

    /// A public struct to access Edition data
    ///
    pub struct EditionData {
        pub let id: UInt64
        pub let seriesID: UInt64
        pub let setID: UInt64
        pub let playID: UInt64
        pub var maxMintSize: UInt64?
        pub let tier: String
        pub var numMinted: UInt64

       /// member function to check if max edition size has been reached
       pub fun maxEditionMintSizeReached(): Bool {
            return self.numMinted == self.maxMintSize
        }

        /// initializer
        ///
        init (id: UInt64) {
            let edition = (&EPL.editionByID[id] as! &EPL.Edition?)!
            self.id = id
            self.seriesID = edition.seriesID
            self.playID = edition.playID
            self.setID = edition.setID
            self.maxMintSize = edition.maxMintSize
            self.tier = edition.tier
            self.numMinted = edition.numMinted
        }
    }

    /// A top level Edition that contains a Series, Set, and Play
    ///
    pub resource Edition {
        pub let id: UInt64
        pub let seriesID: UInt64
        pub let setID: UInt64
        pub let playID: UInt64
        pub let tier: String
        /// Null value indicates that there is unlimited minting potential for the Edition
        pub var maxMintSize: UInt64?
        /// Updates each time we mint a new moment for the Edition to keep a running total
        pub var numMinted: UInt64

        /// Close this edition so that no more Moment NFTs can be minted in it
        ///
        access(contract) fun close() {
            pre {
                self.numMinted != self.maxMintSize: "max number of minted moments has already been reached"
            }

            self.maxMintSize = self.numMinted

            emit EditionClosed(id: self.id)
        }

        /// Mint a Moment NFT in this edition, with the given minting mintingDate.
        /// Note that this will panic if the max mint size has already been reached.
        ///
        pub fun mint(): @EPL.NFT {
            pre {
                self.numMinted != self.maxMintSize: "max number of minted moments has been reached"
            }

            // Create the Moment NFT, filled out with our information
            // Base set moments will not include serial numbers.
            // Future sets may include serial numbers, so leaving data structure in-tact
            let momentNFT <- create NFT(
                editionID: self.id,
                serialNumber: 0
            )
            EPL.totalSupply = EPL.totalSupply + 1
            // Keep a running total (you'll notice we used this as the serial number)
            self.numMinted = self.numMinted + 1 as UInt64

            return <- momentNFT
        }

        /// initializer
        ///
        init (
            seriesID: UInt64,
            setID: UInt64,
            playID: UInt64,
            maxMintSize: UInt64?,
            tier: String,
        ) {
            pre {
                maxMintSize != 0: "max mint size is zero, must either be null or greater than 0"
                EPL.seriesByID.containsKey(seriesID): "seriesID does not exist"
                EPL.setByID.containsKey(setID): "setID does not exist"
                EPL.playByID.containsKey(playID): "playID does not exist"
                EPL.getSeriesData(id: seriesID)!.active == true: "cannot create an Edition with a closed Series"
                EPL.getSetData(id: setID)!.locked == false: "cannot create an Edition with a locked Set"
            }

            self.id = EPL.nextEditionID
            self.seriesID = seriesID
            self.setID = setID
            self.playID = playID

            // If an edition size is not set, it has unlimited minting potential
            if maxMintSize == 0 {
                self.maxMintSize = nil
            } else {
                self.maxMintSize = maxMintSize
            }

            self.tier = tier
            self.numMinted = 0 as UInt64

            EPL.nextEditionID = EPL.nextEditionID + 1 as UInt64

            emit EditionCreated(
                id: self.id,
                seriesID: self.seriesID,
                setID: self.setID,
                playID: self.playID,
                maxMintSize: self.maxMintSize,
                tier: self.tier,
            )
        }
    }

    /// Get the publicly available data for an Edition
    ///
    pub fun getEditionData(id: UInt64): EditionData? {
        if EPL.editionByID[id] == nil{
            return nil
        }

        return EPL.EditionData(id: id)
    }

    // Validate tags exist
    //
    pub fun validateTags(tagIds: [UInt64]): Bool {
        for tagId in tagIds {
            if EPL.tagByID[tagId] == nil {
                return false
            }
        }
        return true
    }

    //------------------------------------------------------------
    // NFT
    //------------------------------------------------------------

    /// A Moment NFT
    ///
    pub resource NFT: NonFungibleToken.INFT, MetadataViews.Resolver {
        pub let id: UInt64
        pub let editionID: UInt64
        pub let serialNumber: UInt64
        pub let mintingDate: UFix64

        /// Destructor
        ///
        destroy() {
            emit MomentNFTBurned(id: self.id, editionID: self.editionID, serialNumber: self.serialNumber)
        }

        /// NFT initializer
        ///
        init(
            editionID: UInt64,
            serialNumber: UInt64
        ) {
            pre {
                EPL.editionByID[editionID] != nil: "no such editionID"
                EditionData(id: editionID).maxEditionMintSizeReached() != true: "max edition size already reached"
            }

            self.id = self.uuid
            self.editionID = editionID
            self.serialNumber = serialNumber
            self.mintingDate = getCurrentBlock().timestamp

            emit MomentNFTMinted(id: self.id, editionID: self.editionID, serialNumber: self.serialNumber)
        }

        pub fun assetPath(): String {
            let editionData = EPL.getEditionData(id: self.editionID)!
            let playDataID: String = EPL.PlayData(id: editionData.playID).metadata["PlayDataID"] ?? ""
            return "https://assets.eplonflow.com/editions/".concat(playDataID).concat("/play_").concat(playDataID)
        }

        pub fun getImage(imageType: String): String {
            return self.assetPath().concat("__").concat(imageType).concat("_2880_2880_").concat(".png")
        }

        pub fun getVideo(videoType: String): String {
            return self.assetPath().concat("__").concat(videoType).concat("_1080_1080_").concat(".mp4")
        }

        /// get the name of an nft
        ///
        pub fun name(): String {
            let editionData = EPL.getEditionData(id: self.editionID)!
            let playerKnownName: String = EPL.PlayData(id: editionData.playID).metadata["PlayerKnownName"] ?? ""
            let playerFirstName: String = EPL.PlayData(id: editionData.playID).metadata["PlayerFirstName"] ?? ""
            let playerLastName: String = EPL.PlayData(id: editionData.playID).metadata["PlayerLastName"] ?? ""
            let playType: String = EPL.PlayData(id: editionData.playID).metadata["PlayType"] ?? ""
            var playerName = playerKnownName
            if(playerName == ""){
                playerName = playerFirstName.concat(" ").concat(playerLastName)
            }
            return playType.concat(" by ").concat(playerName)
        }

        /// get the description of an nft
        ///
        pub fun description(): String {
            let editionData = EPL.getEditionData(id: self.editionID)!
            let metadata = EPL.PlayData(id: editionData.playID).metadata
            let matchHomeTeam: String = metadata["MatchHomeTeam"] ?? ""
            let matchAwayTeam: String = metadata["MatchAwayTeam"] ?? ""
            let matchHomeScore: String = metadata["MatchHomeScore"] ?? ""
            let matchAwayScore: String = metadata["MatchAwayScore"] ?? ""
            let matchDay: String = metadata["MatchDay"] ?? ""
            let matchSeason: String = metadata["MatchSeason"] ?? ""

            return "EPL Moment from ".concat(matchHomeTeam)
            .concat(" x ").concat(matchAwayTeam).concat(" (").concat(matchHomeScore)
            .concat("-").concat(matchAwayScore).concat(") on Matchday ")
            .concat(matchDay).concat(" (").concat(matchSeason).concat(")")
        }

        /// get a thumbnail image that represents this nft
        ///
        pub fun thumbnail(): MetadataViews.HTTPFile {
            let editionData = EPL.getEditionData(id: self.editionID)!
            let playDataID: String = EPL.PlayData(id: editionData.playID).metadata["PlayDataID"] ?? ""
            if playDataID == "" {
                return MetadataViews.HTTPFile(url:"https://ipfs.dapperlabs.com/ipfs/QmPvr5zTwji1UGpun57cbj719MUBsB5syjgikbwCMPmruQ")
            }
            return MetadataViews.HTTPFile(url: self.getImage(imageType: "capture_Hero_Black"))
        }

        /// get the metadata view types available for this nft
        ///
        pub fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.Editions>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.Traits>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.Medias>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Royalties>()
            ]
        }

        /// resolve a metadata view type returning the properties of the view type
        ///
        pub fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name(),
                        description: self.description(),
                        thumbnail: self.thumbnail()
                    )

                case Type<MetadataViews.Editions>():
                let editionData = EPL.getEditionData(id: self.editionID)!
                    let editionInfo = MetadataViews.Edition(
                        name: nil,
                        number: self.serialNumber,
                        max: editionData.maxMintSize
                    )
                    let editionList: [MetadataViews.Edition] = [editionInfo]
                    return MetadataViews.Editions(
                        editionList
                    )

                case Type<MetadataViews.NFTCollectionData>():
                    return MetadataViews.NFTCollectionData(
                        storagePath: EPL.CollectionStoragePath,
                        publicPath: EPL.CollectionPublicPath,
                        providerPath: /private/dapperSportCollection,
                        publicCollection: Type<&EPL.Collection{EPL.MomentNFTCollectionPublic}>(),
                        publicLinkedType: Type<&EPL.Collection{EPL.MomentNFTCollectionPublic, NonFungibleToken.CollectionPublic, NonFungibleToken.Receiver, MetadataViews.ResolverCollection}>(),
                        providerLinkedType: Type<&EPL.Collection{EPL.MomentNFTCollectionPublic, NonFungibleToken.CollectionPublic, NonFungibleToken.Provider, MetadataViews.ResolverCollection}>(),
                        createEmptyCollectionFunction: (fun (): @NonFungibleToken.Collection {
                            return <-EPL.createEmptyCollection()
                        })
                    )
                case Type<MetadataViews.Traits>():
                    let editiondata = EPL.getEditionData(id: self.editionID)!
                    let playdata = EPL.getPlayData(id: editiondata.playID)!
                    return MetadataViews.dictToTraits(dict: playdata.metadata, excludedNames: nil)

                case Type<MetadataViews.ExternalURL>():
                    return MetadataViews.ExternalURL("https://eplonflow.com/moments/".concat(self.id.toString()))

                case Type<MetadataViews.Medias>():
                    return MetadataViews.Medias(
                        items: [
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getImage(imageType: "capture_Hero_Black")),
                                mediaType: "image/png"
                            ),
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getImage(imageType: "capture_Front_Black")),
                                mediaType: "image/png"
                            ),
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getImage(imageType: "capture_Legal_Black")),
                                mediaType: "image/png"
                            ),
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getImage(imageType: "capture_Details_Black")),
                                mediaType: "image/png"
                            ),
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getVideo(videoType: "capture_Animated_Video_Popout_Black")),
                                mediaType: "video/mp4"
                            ),
                             MetadataViews.Media(
                                file: MetadataViews.HTTPFile(url: self.getVideo(videoType: "capture_Animated_Video_Idle_Black")),
                                mediaType: "video/mp4"
                            )
                        ]
                    )
                case Type<MetadataViews.NFTCollectionDisplay>():
                    let bannerImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://assets.eplonflow.com/static/epl-logos/EPL_Logo_Horizontal_B.png"
                        ),
                        mediaType: "image/png"
                    )
                    let squareImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://assets.eplonflow.com/static/epl-logos/EPL_Logo_Primary_B.png"
                        ),
                        mediaType: "image/png"
                    )
                    return MetadataViews.NFTCollectionDisplay(
                        name: "EPL On Flow",
                        description: "Collect EPL's biggest Moments and get closer to the game than ever before",
                        externalURL: MetadataViews.ExternalURL("https://eplonflow.com/"),
                        squareImage: squareImage,
                        bannerImage: bannerImage,
                        socials: {
                            "instagram": MetadataViews.ExternalURL(" https://instagram.com/eplonflow"),
                            "twitter": MetadataViews.ExternalURL("https://twitter.com/EPLOnFlow"),
                            "discord": MetadataViews.ExternalURL("https://discord.gg/EPLOnFlow"),
                            "facebook": MetadataViews.ExternalURL("https://www.facebook.com/EPLOnFlow/")
                        }
                    )
                case Type<MetadataViews.Royalties>():
                    let royaltyReceiver: Capability<&{FungibleToken.Receiver}> =
                        getAccount(EPL.RoyaltyAddress()).getCapability<&AnyResource{FungibleToken.Receiver}>(MetadataViews.getRoyaltyReceiverPublicPath())
                    return MetadataViews.Royalties(
                        royalties: [
                            MetadataViews.Royalty(
                                receiver: royaltyReceiver,
                                cut: 0.05,
                                description: "EPL marketplace royalty"
                            )
                        ]
                    )
            }

            return nil
        }
    }

    //------------------------------------------------------------
    // Collection
    //------------------------------------------------------------

    /// A public collection interface that allows Moment NFTs to be borrowed
    ///
    pub resource interface MomentNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowMomentNFT(id: UInt64): &EPL.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Moment NFT reference: The ID of the returned reference is incorrect"
            }
        }
    }

    /// An NFT Collection
    ///
    pub resource Collection:
        NonFungibleToken.Provider,
        NonFungibleToken.Receiver,
        NonFungibleToken.CollectionPublic,
        MomentNFTCollectionPublic,
        MetadataViews.ResolverCollection
    {
        /// dictionary of NFT conforming tokens
        /// NFT is a resource type with an UInt64 ID field
        ///
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        /// withdraw removes an NFT from the collection and moves it to the caller
        ///
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("Could not find a moment with the given ID in the EPL collection")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        /// deposit takes a NFT and adds it to the collections dictionary
        /// and adds the ID to the id array
        ///
        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @EPL.NFT
            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        /// batchDeposit takes a Collection object as an argument
        /// and deposits each contained NFT into this Collection
        ///
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection) {
            // Get an array of the IDs to be deposited
            let keys = tokens.getIDs()

            // Iterate through the keys in the collection and deposit each one
            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

            // Destroy the empty Collection
            destroy tokens
        }

        /// getIDs returns an array of the IDs that are in the collection
        ///
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        /// borrowNFT gets a reference to an NFT in the collection
        //
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

        /// borrowMomentNFT gets a reference to an NFT in the collection
        ///
        pub fun borrowMomentNFT(id: UInt64): &EPL.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &EPL.NFT
            } else {
                return nil
            }
        }

        pub fun borrowViewResolver(id: UInt64): &{MetadataViews.Resolver} {
            let nft = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
            let eplNFT = nft as! &EPL.NFT
            return eplNFT as &AnyResource{MetadataViews.Resolver}
        }

        /// Collection destructor
        ///
        destroy() {
            destroy self.ownedNFTs
        }

        /// Collection initializer
        ///
        init() {
            self.ownedNFTs <- {}
        }
    }

    /// public function that anyone can call to create a new empty collection
    ///
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }

    //------------------------------------------------------------
    // Admin
    //------------------------------------------------------------

    /// An interface containing the Admin function that allows minting NFTs
    ///
    pub resource interface NFTMinter {
        // Mint a single NFT
        // The Edition for the given ID must already exist
        //
        pub fun mintNFT(editionID: UInt64): @EPL.NFT
    }

    /// A resource that allows managing metadata and minting NFTs
    ///
    pub resource Admin: NFTMinter {
        /// Borrow a Series
        ///
        pub fun borrowSeries(id: UInt64): &EPL.Series {
            pre {
                EPL.seriesByID[id] != nil: "Cannot borrow series, no such id"
            }

            return (&EPL.seriesByID[id] as &EPL.Series?)!
        }

        /// Borrow a Set
        ///
        pub fun borrowSet(id: UInt64): &EPL.Set {
            pre {
                EPL.setByID[id] != nil: "Cannot borrow Set, no such id"
            }

            return (&EPL.setByID[id] as &EPL.Set?)!
        }

        /// Borrow a Play
        ///
        pub fun borrowPlay(id: UInt64): &EPL.Play {
            pre {
                EPL.playByID[id] != nil: "Cannot borrow Play, no such id"
            }

            return (&EPL.playByID[id] as &EPL.Play?)!
        }

        /// Borrow an Edition
        ///
        pub fun borrowEdition(id: UInt64): &EPL.Edition {
            pre {
                EPL.editionByID[id] != nil: "Cannot borrow edition, no such id"
            }

            return (&EPL.editionByID[id] as &EPL.Edition?)!
        }

        /// Create a Series
        ///
        pub fun createSeries(name: String): UInt64 {
            // Create and store the new series
            let series <- create EPL.Series(
                name: name,
            )
            let seriesID = series.id
            EPL.seriesByID[series.id] <-! series

            // Return the new ID for convenience
            return seriesID
        }

        /// Close a Series
        ///
        pub fun closeSeries(id: UInt64): UInt64 {
            let series = (&EPL.seriesByID[id] as &EPL.Series?)!
            series.close()
            return series.id
        }

        /// Create a Set
        ///
        pub fun createSet(name: String): UInt64 {
            // Create and store the new set
            let set <- create EPL.Set(
                name: name,
            )
            let setID = set.id
            EPL.setByID[set.id] <-! set

            // Return the new ID for convenience
            return setID
        }

        /// Locks a Set
        ///
        pub fun lockSet(id: UInt64): UInt64 {
            let set = (&EPL.setByID[id] as &EPL.Set?)!
            set.lock()
            return set.id
        }

        /// Create a Tag
        ///
        pub fun createTag(name: String): UInt64 {
            // Create and store the new tag
            let tag <- create EPL.Tag(
                name: name
            )
            let tagID = tag.id
            EPL.tagByID[tag.id] <-! tag

            // Return the new ID for convenience
            return tagID
        }

        /// Create a Play
        ///
        pub fun createPlay(metadata: {String: String}, tagIds: [UInt64]): UInt64 {
            // Create and store the new play
            let play <- create EPL.Play(
                metadata: metadata,
                tagIds: tagIds
            )
            let playID = play.id
            EPL.playByID[play.id] <-! play

            // Return the new ID for convenience
            return playID
        }

        /// Create an Edition
        ///
        pub fun createEdition(
            seriesID: UInt64,
            setID: UInt64,
            playID: UInt64,
            maxMintSize: UInt64?,
            tier: String): UInt64 {
            let edition <- create Edition(
                seriesID: seriesID,
                setID: setID,
                playID: playID,
                maxMintSize: maxMintSize,
                tier: tier,
            )
            let editionID = edition.id
            EPL.editionByID[edition.id] <-! edition

            return editionID
        }

        /// Close an Edition
        ///
        pub fun closeEdition(id: UInt64): UInt64 {
            let edition = (&EPL.editionByID[id] as &EPL.Edition?)!
            edition.close()
            return edition.id
        }

        /// Mint a single NFT
        /// The Edition for the given ID must already exist
        ///
        pub fun mintNFT(editionID: UInt64): @EPL.NFT {
            pre {
                // Make sure the edition we are creating this NFT in exists
                EPL.editionByID.containsKey(editionID): "No such EditionID"
            }

            return <- self.borrowEdition(id: editionID).mint()
        }
    }

    //------------------------------------------------------------
    // Contract lifecycle
    //------------------------------------------------------------

    /// EPL contract initializer
    ///
    init() {
        // Set the named paths
        self.CollectionStoragePath = /storage/EPLNFTCollection
        self.CollectionPublicPath = /public/EPLNFTCollection
        self.AdminStoragePath = /storage/EPLAdmin
        self.MinterPrivatePath = /private/EPLMinter

        // Initialize the entity counts
        self.totalSupply = 0
        self.nextSeriesID = 1
        self.nextSetID = 1
        self.nextTagID = 1
        self.nextPlayID = 1
        self.nextEditionID = 1

        // Initialize the metadata lookup dictionaries
        self.seriesByID <- {}
        self.seriesIDByName = {}
        self.setIDByName = {}
        self.setByID <- {}
        self.tagByID <- {}
        self.playByID <- {}
        self.editionByID <- {}

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.save(<-admin, to: self.AdminStoragePath)
        // Link capabilites to the admin constrained to the Minter
        // and Metadata interfaces
        self.account.link<&EPL.Admin{EPL.NFTMinter}>(
            self.MinterPrivatePath,
            target: self.AdminStoragePath
        )

        // Let the world know we are here
        emit ContractInitialized()
    }
}