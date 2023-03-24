/*
    Contract for EPL sports moments & metadata
    Author: Jeremy Ahrens jer.ahrens@dapperlabs.com
*/

import FungibleToken from 0x{{.FungibleTokenAddress}}
import NonFungibleToken from 0x{{.NonFungibleTokenAddress}}
import MetadataViews from 0x{{.MetadataViewsAddress}}

/// The English Premier League NFT and metadata contract
///
pub contract EnglishPremierLeague: NonFungibleToken {

    /// Contract events
    ///
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event SeriesCreated(id: UInt64, name: String)
    pub event SeriesClosed(id: UInt64)
    pub event SetCreated(id: UInt64, name: String)
    pub event SetLocked(setID: UInt64)
    pub event PlayCreated(id: UInt64, metadata: {String: String}, tagIds: [UInt64])
    pub event TagCreated(id: UInt64, name: String)
    pub event EditionCreated(
        id: UInt64,
        seriesID: UInt64,
        setID: UInt64,
        playID: UInt64,
        maxMintSize: UInt64?,
        tier: String,
    )
    pub event EditionClosed(id: UInt64)
    pub event MomentNFTMinted(id: UInt64, editionID: UInt64, serialNumber: UInt64)
    pub event MomentNFTBurned(id: UInt64,  editionID: UInt64, serialNumber: UInt64)

    /// Named Paths
    ///
    pub let CollectionStoragePath:  StoragePath
    pub let CollectionPublicPath:   PublicPath
    pub let AdminStoragePath:       StoragePath
    pub let MinterPrivatePath:      PrivatePath

    /// Contract variables
    ///
    pub var royaltyAddress: Address
    pub var totalSupply:        UInt64
    pub var nextSeriesID:       UInt64
    pub var nextSetID:          UInt64
    pub var nextTagID:          UInt64
    pub var nextPlayID:         UInt64
    pub var nextEditionID:      UInt64

    /// Metadata Dictionaries
    ///
    access(self) let seriesIDByName:    {String: UInt64}
    access(self) let seriesByID:        {UInt64: Series}
    access(self) let setIDByName:       {String: UInt64}
    access(self) let setByID:           {UInt64: Set}
    access(self) let tagByID:           {UInt64: Tag}
    access(self) let playByID:          {UInt64: Play}
    access(self) let editionByID:       {UInt64: Edition}

    /// A top-level Series with a unique ID and name
    ///
    pub struct Series {
        pub let id: UInt64
        pub let name: String
        pub var active: Bool

        init (id: UInt64, name: String) {
            if let series = EnglishPremierLeague.seriesByID[id] {
                self.id = series.id
                self.name = series.name
                self.active = series.active
            } else {
                self.id = id
                self.name = name
                self.active = true
            }
        }

        pub fun close() { self.active = false }
    }

    /// Series getters
    ///
    pub fun getSeries(id: UInt64): EnglishPremierLeague.Series? {
        return EnglishPremierLeague.seriesByID[id]
    }

    pub fun getSeriesByName(name: String): EnglishPremierLeague.Series? {
        let id = EnglishPremierLeague.seriesIDByName[name]

        if id == nil{
            return nil
        }

        return EnglishPremierLeague.getSeries(id: id!)
    }

    pub fun getAllSeriesNames(): [String] {
        return EnglishPremierLeague.seriesIDByName.keys
    }

    pub fun getSeriesIDByName(name: String): UInt64? {
        return EnglishPremierLeague.seriesIDByName[name]
    }

    /// A top-level Set with ID and name
    ///
    pub struct Set {
        pub let id: UInt64
        pub let name: String
        pub var locked: Bool

        init (id: UInt64, name: String) {
            if let set = EnglishPremierLeague.setByID[id] {
                self.id = set.id
                self.name = set.name
                self.locked = set.locked
            } else {
                self.id = id
                self.name = name
                self.locked = false
            }
        }

        access(contract) fun lock() { self.locked = true }
    }

    /// Set getters
    ///
    pub fun getSet(id: UInt64): EnglishPremierLeague.Set? {
        return EnglishPremierLeague.setByID[id]
    }

    pub fun getSetByName(name: String): EnglishPremierLeague.Set? {
        let id = EnglishPremierLeague.setIDByName[name]!

        if id == nil {
            return nil
        }

        return EnglishPremierLeague.setByID[id!]
    }

    pub fun getAllSetNames(): [String] {
        return EnglishPremierLeague.setIDByName.keys
    }

    /// A top-level Tag with ID and name
    /// Examples include goal, match winner, assist, header, save, brace
    ///
    pub struct Tag {
        pub let id: UInt64
        pub let name: String

        init (id: UInt64, name: String) {
            if let tag = EnglishPremierLeague.tagByID[id] {
                self.id = tag.id
                self.name = tag.name
            } else {
                self.id = id
                self.name = name
            }
        }
    }

    /// Tag getters
    ///
    pub fun getTag(id: UInt64): EnglishPremierLeague.Tag? {
        return EnglishPremierLeague.tagByID[id]
    }

    /// A top level Play with a unique ID
    ///
    pub struct Play {
        pub let id: UInt64
        pub let metadata: {String: String}
        pub var tagIds: [UInt64]

        init (id: UInt64, metadata: {String: String}, tagIds: [UInt64]) {
            if let play = EnglishPremierLeague.playByID[id] {
                self.id = play.id
                self.metadata = play.metadata
                self.tagIds = play.tagIds
            } else {
                self.id = id
                self.metadata = metadata
                self.tagIds = tagIds
            }
        }
    }

    /// Play getters
    ///
    pub fun getPlay(id: UInt64): EnglishPremierLeague.Play? {
        return EnglishPremierLeague.playByID[id]
    }

    /// A public struct to access Edition data
    ///
    pub struct Edition {
        pub let id: UInt64
        pub let seriesID: UInt64
        pub let setID: UInt64
        pub let playID: UInt64
        pub var maxMintSize: UInt64?
        pub let tier: String
        pub var numMinted: UInt64

        init (id: UInt64, seriesID: UInt64, setID: UInt64, playID: UInt64, maxMintSize: UInt64?, tier: String) {
            if let edition = EnglishPremierLeague.editionByID[id] {
                self.id = edition.id
                self.seriesID = edition.seriesID
                self.setID = edition.setID
                self.playID = edition.playID
                self.maxMintSize = edition.maxMintSize
                self.tier = edition.tier
                self.numMinted = edition.numMinted
            } else {
                self.id = id
                self.seriesID = seriesID
                self.setID = setID
                self.playID = playID
                self.maxMintSize = maxMintSize

                // If an edition size is not set, it has unlimited minting potential
                if maxMintSize != nil && maxMintSize! == 0 {
                   self.maxMintSize = nil
                } else {
                  self.maxMintSize = maxMintSize
                }

                self.tier = tier
                self.numMinted = 0 as UInt64
            }
        }

        pub fun maxEditionMintSizeReached(): Bool {
            return self.numMinted == self.maxMintSize
        }

        access(contract) fun incrementNumMinted() {
            pre {
                self.numMinted != self.maxMintSize: "max number of minted moments has already been reached"
            }
            self.numMinted = self.numMinted + 1
        }

        access(contract) fun close() {
            pre {
                self.numMinted != self.maxMintSize: "max number of minted moments has already been reached"
            }

            self.maxMintSize = self.numMinted
        }
    }

    /// Edition getters
    ///
    pub fun getEdition(id: UInt64): Edition? {
        return EnglishPremierLeague.editionByID[id]
    }

    /// Validate tags exist
    ///
    pub fun validateTags(tagIds: [UInt64]): Bool {
        for tagId in tagIds {
            if EnglishPremierLeague.tagByID[tagId] == nil {
                return false
            }
        }
        return true
    }

    /// A Moment NFT
    ///
    pub resource NFT: NonFungibleToken.INFT, MetadataViews.Resolver {
        pub let id: UInt64
        pub let editionID: UInt64
        pub let serialNumber: UInt64
        pub let mintingDate: UFix64
        pub let ext: {String: AnyStruct}

        /// Destructor
        ///
        destroy() {
            emit MomentNFTBurned(id: self.id, editionID: self.editionID, serialNumber: self.serialNumber)
        }

        /// NFT initializer
        ///
        init(
            editionID: UInt64,
            serialNumber: UInt64,
            ext: {String: AnyStruct}
        ) {
            pre {
                EnglishPremierLeague.editionByID[editionID] != nil: "no such editionID"
                (&EnglishPremierLeague.editionByID[editionID] as &EnglishPremierLeague.Edition?)!
                    .maxEditionMintSizeReached() != true : "max edition size already reached"
            }

            self.id = self.uuid
            self.editionID = editionID
            self.serialNumber = serialNumber
            self.mintingDate = getCurrentBlock().timestamp
            self.ext = ext

            emit MomentNFTMinted(
                id: self.id,
                editionID: self.editionID,
                serialNumber: self.serialNumber
            )
        }

        pub fun assetPath(): String {
            let editionData = EnglishPremierLeague.getEdition(id: self.editionID)!
            let playDataID: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["PlayDataID"] ?? ""
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
            let editionData = EnglishPremierLeague.getEdition(id: self.editionID)!
            let playerKnownName: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["Player Known Name"] ?? ""
            let playerFirstName: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["Player First Name"] ?? ""
            let playerLastName: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["Player Last Name"] ?? ""
            let playType: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["Play Type"] ?? ""
            var playerName = playerKnownName
            if(playerName == ""){
                playerName = playerFirstName.concat(" ").concat(playerLastName)
            }
            return playType.concat(" by ").concat(playerName)
        }

        /// get the description of an nft
        ///
        pub fun description(): String {
            let editionData = EnglishPremierLeague.getEdition(id: self.editionID)!
            let metadata = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata
            let matchHomeTeam: String = metadata["Match Home Team"] ?? ""
            let matchAwayTeam: String = metadata["Match Away Team"] ?? ""
            let matchHomeScore: String = metadata["Match Home Score"] ?? ""
            let matchAwayScore: String = metadata["Match Away Score"] ?? ""
            let matchDay: String = metadata["Match Day"] ?? ""
            let matchSeason: String = metadata["Match Season"] ?? ""

            return "EnglishPremierLeague Moment from ".concat(matchHomeTeam)
            .concat(" x ").concat(matchAwayTeam).concat(" (").concat(matchHomeScore)
            .concat("-").concat(matchAwayScore).concat(") on Matchday ")
            .concat(matchDay).concat(" (").concat(matchSeason).concat(")")
        }

        /// get a thumbnail image that represents this nft
        ///
        pub fun thumbnail(): MetadataViews.HTTPFile {
            let editionData = EnglishPremierLeague.getEdition(id: self.editionID)!
            let playDataID: String = EnglishPremierLeague.getPlay(id: editionData.playID)!.metadata["PlayDataID"] ?? ""
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
                let editionData = EnglishPremierLeague.getEdition(id: self.editionID)!
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
                        storagePath: EnglishPremierLeague.CollectionStoragePath,
                        publicPath: EnglishPremierLeague.CollectionPublicPath,
                        providerPath: /private/dapperSportCollection,
                        publicCollection: Type<&EnglishPremierLeague.Collection{EnglishPremierLeague.MomentNFTCollectionPublic}>(),
                        publicLinkedType: Type<&EnglishPremierLeague.Collection{EnglishPremierLeague.MomentNFTCollectionPublic, NonFungibleToken.CollectionPublic, NonFungibleToken.Receiver, MetadataViews.ResolverCollection}>(),
                        providerLinkedType: Type<&EnglishPremierLeague.Collection{EnglishPremierLeague.MomentNFTCollectionPublic, NonFungibleToken.CollectionPublic, NonFungibleToken.Provider, MetadataViews.ResolverCollection}>(),
                        createEmptyCollectionFunction: (fun (): @NonFungibleToken.Collection {
                            return <-EnglishPremierLeague.createEmptyCollection()
                        })
                    )
                case Type<MetadataViews.Traits>():
                    let editiondata = EnglishPremierLeague.getEdition(id: self.editionID)!
                    let play = EnglishPremierLeague.getPlay(id: editiondata.playID)!
                    return MetadataViews.dictToTraits(dict: play.metadata, excludedNames: nil)

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
                        getAccount(EnglishPremierLeague.royaltyAddress).getCapability<&AnyResource{FungibleToken.Receiver}>(MetadataViews.getRoyaltyReceiverPublicPath())
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

    /// A public collection interface that allows Moment NFTs to be borrowed
    ///
    pub resource interface MomentNFTCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowNFTSafe(id: UInt64): &NonFungibleToken.NFT?
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
            let token <- token as! @EnglishPremierLeague.NFT
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
        ///
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

        /// borrowNFTSafe gets a reference to an NFT in the collection
        ///
        pub fun borrowNFTSafe(id: UInt64): &NonFungibleToken.NFT? {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)
        }

        pub fun borrowViewResolver(id: UInt64): &{MetadataViews.Resolver} {
            let nft = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
            let eplNFT = nft as! &EnglishPremierLeague.NFT
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

    /// An interface containing the Admin function that allows minting NFTs
    ///
    pub resource interface NFTMinter {
        pub fun mintNFT(editionID: UInt64, ext: {String: AnyStruct}): @EnglishPremierLeague.NFT
    }

    /// A resource that allows managing metadata and minting NFTs
    ///
    pub resource Admin: NFTMinter {

        /// Series
        ///
        pub fun createSeries(name: String): UInt64 {
            pre {
                !EnglishPremierLeague.seriesIDByName.containsKey(name): "A Series with that name already exists"
            }

            let series = EnglishPremierLeague.Series(
                id: EnglishPremierLeague.nextSeriesID,
                name: name,
            )
            EnglishPremierLeague.seriesByID[series.id] = series
            EnglishPremierLeague.seriesIDByName[name] = series.id
            emit SeriesCreated(id: series.id, name: series.name)
            EnglishPremierLeague.nextSeriesID = series.id + 1 as UInt64
            return series.id
        }

        pub fun closeSeries(id: UInt64): UInt64 {
            let series = (&EnglishPremierLeague.seriesByID[id] as &EnglishPremierLeague.Series?)!

            assert(
                series.active == true,
                message: "series is already inactive"
            )

            series.close()
            emit SeriesClosed(id: series.id)
            return series.id
        }

        /// Set
        ///
        pub fun createSet(name: String): UInt64 {
            let set = EnglishPremierLeague.Set(
                id: EnglishPremierLeague.nextSetID,
                name: name,
            )

            EnglishPremierLeague.setByID[set.id] = set
            emit SetCreated(id: set.id, name: set.name)
            EnglishPremierLeague.nextSetID = set.id + 1 as UInt64
            return set.id
        }

        pub fun lockSet(id: UInt64) {
            let set = (&EnglishPremierLeague.setByID[id] as &EnglishPremierLeague.Set?)!

            assert(
                set.locked == false,
                message: "set is already locked"
            )
            set.lock()
            emit SetLocked(setID: id)
        }

        /// Tag
        ///
        pub fun createTag(name: String): UInt64 {
            let tag = EnglishPremierLeague.Tag(
                id: EnglishPremierLeague.nextTagID,
                name: name,
            )

            EnglishPremierLeague.tagByID[tag.id] = tag
            emit TagCreated(id: tag.id, name: tag.name)
            EnglishPremierLeague.nextTagID = tag.id + 1 as UInt64
            return tag.id
        }

        /// Play
        ///
        pub fun createPlay(metadata: {String: String}, tagIds: [UInt64]): UInt64 {
            pre {
                EnglishPremierLeague.validateTags(
                    tagIds: tagIds
                ) == true : "Play contains tag that does not exist."
            }

            let play = EnglishPremierLeague.Play(
                id: EnglishPremierLeague.nextPlayID,
                metadata: metadata,
                tagIds: tagIds
            )

            EnglishPremierLeague.playByID[play.id] = play
            emit PlayCreated(id: play.id, metadata: play.metadata, tagIds: play.tagIds)
            EnglishPremierLeague.nextPlayID = play.id + 1 as UInt64
            return play.id
        }

        /// Edition
        ///
        pub fun createEdition(
            seriesID: UInt64,
            setID: UInt64,
            playID: UInt64,
            maxMintSize: UInt64?,
            tier: String): UInt64 {

            pre {
                maxMintSize != 0: "max mint size is zero, must either be null or greater than 0"
                EnglishPremierLeague.seriesByID.containsKey(seriesID): "seriesID does not exist"
                EnglishPremierLeague.setByID.containsKey(setID): "setID does not exist"
                EnglishPremierLeague.playByID.containsKey(playID): "playID does not exist"
                EnglishPremierLeague.getSeries(id: seriesID)!.active == true: "cannot create an Edition with a closed Series"
                EnglishPremierLeague.getSet(id: setID)!.locked == false: "cannot create an Edition with a locked Set"
            }

            let edition = EnglishPremierLeague.Edition(
                id: EnglishPremierLeague.nextEditionID,
                seriesID: seriesID,
                setID: setID,
                playID: playID,
                maxMintSize: maxMintSize,
                tier: tier
            )

            EnglishPremierLeague.editionByID[edition.id] = edition
            emit EditionCreated(
                id: edition.id,
                seriesID: edition.seriesID,
                setID: edition.setID,
                playID: edition.playID,
                maxMintSize: edition.maxMintSize,
                tier: edition.tier,
            )
            EnglishPremierLeague.nextEditionID = edition.id + 1 as UInt64
            return edition.id
        }


        pub fun closeEdition(id: UInt64): UInt64 {
            let edition = (&EnglishPremierLeague.editionByID[id] as &EnglishPremierLeague.Edition?)!
            edition.close()
            emit EditionClosed(id: edition.id)
            return edition.id
        }

        pub fun mintNFT(editionID: UInt64, ext: {String: AnyStruct}): @EnglishPremierLeague.NFT {
            pre {
                // Make sure the edition we are creating this NFT in exists
                EnglishPremierLeague.editionByID.containsKey(editionID): "No such EditionID"
            }

            let edition = (&EnglishPremierLeague.editionByID[editionID] as &EnglishPremierLeague.Edition?)!
            assert(
                !edition.maxEditionMintSizeReached(),
                message: "edition has reached capacity"
            )

            // Moments will not include serial numbers.
            let momentNFT <- create NFT(
                editionID: edition.id,
                serialNumber: 0,
                ext: ext
            )
            edition.incrementNumMinted()
            return <- momentNFT
        }

        /// Royalty Address
        ///
        pub fun setRoyaltyAddress(royaltyAddress: Address): Address {
            EnglishPremierLeague.royaltyAddress = royaltyAddress
            return EnglishPremierLeague.royaltyAddress
        }
    }

    /// EPL contract initializer
    ///
    init() {
        // Set the named paths
        self.CollectionStoragePath = /storage/EnglishPremierLeagueNFTCollection
        self.CollectionPublicPath = /public/EnglishPremierLeagueNFTCollection
        self.AdminStoragePath = /storage/EnglishPremierLeagueAdmin
        self.MinterPrivatePath = /private/EnglishPremierLeagueMinter

        // Initialize the entity counts
        self.totalSupply = 0
        self.nextSeriesID = 1
        self.nextSetID = 1
        self.nextTagID = 1
        self.nextPlayID = 1
        self.nextEditionID = 1

        // Initialize the metadata lookup dictionaries
        self.seriesByID = {}
        self.seriesIDByName = {}
        self.setIDByName = {}
        self.setByID = {}
        self.tagByID = {}
        self.playByID = {}
        self.editionByID = {}

        self.royaltyAddress = 0xf8d6e0586b0a20c7

        // Create an Admin resource and save it to storage
        let admin <- create Admin()
        self.account.save(<-admin, to: self.AdminStoragePath)
        // Link capabilites to the admin constrained to the Minter
        // and Metadata interfaces
        self.account.link<&EnglishPremierLeague.Admin{EnglishPremierLeague.NFTMinter}>(
            self.MinterPrivatePath,
            target: self.AdminStoragePath
        )

        // Let the world know we are here
        emit ContractInitialized()
    }
}