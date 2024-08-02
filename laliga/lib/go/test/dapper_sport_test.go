package test

import (
	"fmt"
	"testing"

	emulator "github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"
)

const (
	metadatakey1         = "name"
	metadatakey2         = "PlayerJerseyName"
	metadatakey3         = "PlayType"
	playerKnownNameKey   = "PlayerKnownName"
	playerKnownNameValue = "Deewai"
	matchHomeTeamKey     = "MatchHomeTeam"
	matchHomeTeamValue   = "Manchester United"
	matchAwayTeamKey     = "MatchAwayTeam"
	matchAwayTeamValue   = "Barcelona"
	matchHomeScoreKey    = "MatchHomeScore"
	matchHomeScoreValue  = "5"
	matchAwayScoreKey    = "MatchAwayScore"
	matchAwayScoreValue  = "1"
	matchDayKey          = "MatchDay"
	matchDayValue        = "15"
	matchSeasonKey       = "MatchSeason"
	matchSeasonValue     = "2022/2023"
	playerJerseyName     = "Deewai"
	playType             = "Goal"
	play1Name            = "TEST PLAY METADATA A"
	play2Name            = "TEST PLAY METADATA B"
)

var (
	editions  = make(map[uint64]EditionData)
	sets      = make(map[uint64]SetData)
	series    = make(map[uint64]SeriesData)
	playNames = []string{
		play1Name,
		play2Name,
	}
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func TestGolazosDeployContracts(t *testing.T) {
	b := newEmulator()
	GolazosDeployContracts(t, b)
}

func TestGolazosSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupGolazos(t, b, userAddress, userSigner, contracts)

	t.Run("Account should be set up", func(t *testing.T) {
		accountIsSetUp := accountIsSetup(
			t,
			b,
			contracts,
			userAddress,
		)
		assert.Equal(t, true, accountIsSetUp)
	})
}

// ------------------------------------------------------------
// Series
// ------------------------------------------------------------
func TestSeries(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	createTestSeries(t, b, contracts)
}

func createTestSeries(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new series", func(t *testing.T) {
		testCreateSeries(
			t,
			b,
			contracts,
			"Series One",
			1,
			false,
		)
	})

	t.Run("Should be able to create a new series with an incrementing ID from the first", func(t *testing.T) {
		testCreateSeries(
			t,
			b,
			contracts,
			"Series Two",
			2,
			false,
		)
	})

	t.Run("Should be able to close a series", func(t *testing.T) {
		testCloseSeries(
			t,
			b,
			contracts,
			2,
			false,
		)
	})
}

func testCreateSeries(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	seriesName string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createSeries(
		t,
		b,
		contracts,
		seriesName,
		false,
	)

	if !shouldRevert {
		seriesData := getSeriesData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, seriesData.ID)
		assert.Equal(t, seriesName, seriesData.Name)
		assert.Equal(t, true, seriesData.Active)
		series[seriesData.ID] = seriesData
	}
}

func testCloseSeries(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	seriesID uint64,
	shouldRevert bool,
) {
	wasActive := getSeriesData(t, b, contracts, seriesID).Active
	closeSeries(
		t,
		b,
		contracts,
		seriesID,
		shouldRevert,
	)

	series := getSeriesData(t, b, contracts, seriesID)
	assert.Equal(t, seriesID, series.ID)
	if !shouldRevert {
		assert.Equal(t, false, series.Active)
	} else {
		assert.Equal(t, wasActive, series.Active)
	}
}

// ------------------------------------------------------------
// Sets
// ------------------------------------------------------------
func TestSets(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	createTestSets(t, b, contracts)

}

func createTestSets(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new set", func(t *testing.T) {
		testCreateSet(
			t,
			b,
			contracts,
			"Set One",
			1,
			false,
		)
	})

	t.Run("Should be able to create a new set with an incrementing ID from the first", func(t *testing.T) {
		testCreateSet(
			t,
			b,
			contracts,
			"Set Two",
			2,
			false,
		)
	})

	t.Run("Should create a new set to test set locking", func(t *testing.T) {
		testCreateSet(
			t,
			b,
			contracts,
			"Set Three",
			3,
			false,
		)
	})

	t.Run("Should be able to lock a set", func(t *testing.T) {
		testLockSet(
			t,
			b,
			contracts,
			3,
			false,
		)
	})
}

func testCreateSet(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	setName string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createSet(
		t,
		b,
		contracts,
		setName,
		false,
	)

	if !shouldRevert {
		set := getSetData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, set.ID)
		assert.Equal(t, setName, set.Name)
		sets[set.ID] = set
	}
}

func testLockSet(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	setID uint64,
	shouldRevert bool,
) {
	wasLocked := getSetData(t, b, contracts, setID).Locked
	lockSet(
		t,
		b,
		contracts,
		setID,
		shouldRevert,
	)

	set := getSetData(t, b, contracts, setID)
	assert.Equal(t, setID, set.ID)
	if !shouldRevert {
		assert.Equal(t, true, set.Locked)
	} else {
		assert.Equal(t, wasLocked, set.Locked)
	}
}

// ------------------------------------------------------------
// Plays
// ------------------------------------------------------------
func TestPlays(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	createTestPlays(t, b, contracts)
}

func createTestPlays(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new play", func(t *testing.T) {
		testCreatePlay(
			t,
			b,
			contracts,
			"TEST_CLASSIFICATION",
			map[string]string{metadatakey1: play1Name, metadatakey2: playerJerseyName,
				metadatakey3: playType, playerKnownNameKey: playerKnownNameValue,
				matchHomeTeamKey: matchHomeTeamValue, matchAwayTeamKey: matchAwayTeamValue,
				matchHomeScoreKey: matchHomeScoreValue, matchAwayScoreKey: matchAwayScoreValue,
				matchDayKey: matchDayValue, matchSeasonKey: matchSeasonValue},
			1,
			false,
		)
	})

	t.Run("Should be able to create a new play with an incrementing ID from the first", func(t *testing.T) {
		testCreatePlay(
			t,
			b,
			contracts,
			"TEST_CLASSIFICATION",
			map[string]string{metadatakey1: play2Name, metadatakey2: playerJerseyName,
				metadatakey3: playType, playerKnownNameKey: playerKnownNameValue,
				matchHomeTeamKey: matchHomeTeamValue, matchAwayTeamKey: matchAwayTeamValue,
				matchHomeScoreKey: matchHomeScoreValue, matchAwayScoreKey: matchAwayScoreValue,
				matchDayKey: matchDayValue, matchSeasonKey: matchSeasonValue},
			2,
			false,
		)
	})
}

func testCreatePlay(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	classification string,
	metadata map[string]string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createPlay(
		t,
		b,
		contracts,
		classification,
		metadata,
		false,
	)

	if !shouldRevert {
		play := getPlayData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, play.ID)
		assert.Equal(t, classification, play.Classification)
		assert.Equal(t, metadata, play.Metadata)
	}
}

// ------------------------------------------------------------
// Editions
// ------------------------------------------------------------
func TestEditions(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	createTestEditions(t, b, contracts)
}

func testCreateEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	seriesID uint64,
	setID uint64,
	playID uint64,
	maxMintSize *uint64,
	tier string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createEdition(
		t,
		b,
		contracts,
		seriesID,
		setID,
		playID,
		maxMintSize,
		tier,
		shouldRevert,
	)

	if !shouldRevert {
		edition := getEditionData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, edition.ID)
		assert.Equal(t, seriesID, edition.SeriesID)
		assert.Equal(t, setID, edition.SetID)
		assert.Equal(t, playID, edition.PlayID)
		assert.Equal(t, tier, edition.Tier)
		if maxMintSize != nil {
			assert.Equal(t, &maxMintSize, &edition.MaxMintSize)
		}
		editions[edition.ID] = edition
	}
}

func testCloseEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	editionID uint64,
	shouldBeID uint64,
	shouldRevert bool,
) {
	closeEdition(
		t,
		b,
		contracts,
		editionID,
		false,
	)

	if !shouldRevert {
		edition := getEditionData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, edition.ID)
	}
}

func createTestEditions(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	var maxMintSize uint64 = 2
	createTestSeries(t, b, contracts)
	createTestSets(t, b, contracts)
	createTestPlays(t, b, contracts)

	t.Run("Should be able to create a new edition with series/set/play IDs and a max mint size of 100", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			1,
			1,
			1,
			&maxMintSize,
			"COMMON",
			1,
			false,
		)
	})

	t.Run("Should be able to create another new edition with series/set/play IDs and no max mint size", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			1,
			2,
			1,
			nil,
			"COMMON",
			2,
			false,
		)
	})

	t.Run("Should be able to create a new edition with series/set/play IDs and no max mint size", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			1,
			1,
			2,
			nil,
			"COMMON",
			3,
			false,
		)
	})

	t.Run("Should not be able to create a new edition with a closed series", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			2,
			1,
			1,
			nil,
			"COMMON",
			4,
			true,
		)
	})

	t.Run("Should not be able to create an Edition with a Set/Play combination that already exists", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			1,
			1,
			2,
			nil,
			"COMMON",
			5,
			true,
		)
	})

	t.Run("Should not be able to create a new edition with a locked set", func(t *testing.T) {
		testCreateEdition(
			t,
			b,
			contracts,
			1,
			3,
			1,
			nil,
			"COMMON",
			4,
			true,
		)
	})

	t.Run("Should be able to close and edition that has no max mint size", func(t *testing.T) {
		testCloseEdition(
			t,
			b,
			contracts,
			3,
			3,
			false,
		)
	})

}

// ------------------------------------------------------------
// MomentNFTs
// ------------------------------------------------------------
func TestMomentNFTs(t *testing.T) {
	b := newEmulator()
	contracts := GolazosDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupGolazos(t, b, userAddress, userSigner, contracts)

	createTestEditions(t, b, contracts)

	t.Run("Should be able to mint a new MomentNFT from an edition that has a maxMintSize", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(1),
			false,
		)
	})

	t.Run("Should be able to mint a second new MomentNFT from an edition that has a maxmintSize", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(2),
			false,
		)
	})

	t.Run("Should be able to mint a new MomentNFT from an edition with no max mint size", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(2),
			userAddress,
			uint64(1),
			false,
		)
	})

	t.Run("Should be able to mint a second new MomentNFT from an edition with no max mint size", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(2),
			userAddress,
			uint64(2),
			false,
		)
	})

	t.Run("Should not be able to mint an edition that has reached max minting size", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(3),
			true,
		)
	})

	t.Run("Should not be able to mint an edition that is already closed", func(t *testing.T) {
		testMintMomentNFT(
			t,
			b,
			contracts,
			uint64(3),
			userAddress,
			uint64(1),
			true,
		)
	})
}

func testMintMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	editionID uint64,
	userAddress flow.Address,
	shouldBeSerialNumber uint64,
	shouldRevert bool,
) {
	// Make sure the total supply of NFTs is tracked correctly
	previousSupply := getMomentNFTSupply(t, b, contracts)

	nftID := mintMomentNFT(
		t,
		b,
		contracts,
		userAddress,
		editionID,
		shouldRevert,
	)

	newSupply := getMomentNFTSupply(t, b, contracts)
	if !shouldRevert {
		assert.Equal(t, previousSupply+uint64(1), newSupply)

		nftProperties := getMomentNFTProperties(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Equal(t, editionID, nftProperties.EditionID)
		assert.Equal(t, shouldBeSerialNumber, nftProperties.SerialNumber)
		assert.Equal(t, shouldBeSerialNumber, nftProperties.SerialNumber)
		//FIXME: query the block time and check equality.
		//       Here we just make sure it's nonzero.
		assert.Less(t, uint64(0), nftProperties.MintingDate)
		displayView := getMomentNFTDisplayMetadataView(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Equal(t, playType+" by "+playerKnownNameValue, displayView.Name)
		assert.Equal(t, "LaLiga Golazos Moment from "+matchHomeTeamValue+" x "+matchAwayTeamValue+" ("+matchHomeScoreValue+"-"+matchAwayScoreValue+") on Matchday "+matchDayValue+" ("+matchSeasonValue+")", displayView.Description)
		//TODO: check the image reurned based on tier
		assert.Equal(t, "https://ipfs.dapperlabs.com/ipfs/QmPvr5zTwji1UGpun57cbj719MUBsB5syjgikbwCMPmruQ", displayView.ImageURL)

		editionView := getMomentNFTEditionMetadataView(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Equal(t, nftProperties.SerialNumber, editionView.Number)
		assert.Equal(t, *editions[editionID].MaxMintSize, editionView.Max)

		serialView := getMomentNFTSerialMetadataView(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Equal(t, shouldBeSerialNumber, serialView)

		nftCollectionDataView := getMomentNFTCollectionDataMetadataView(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Equal(t, "GolazosNFTCollection", nftCollectionDataView.StoragePath)
		assert.Equal(t, "GolazosNFTCollection", nftCollectionDataView.PublicPath)
		assert.Equal(t, fmt.Sprintf("&A.%s.Golazos.Collection",
			contracts.GolazosAddress.Hex()), nftCollectionDataView.PublicCollection)
		assert.Equal(t, fmt.Sprintf("&A.%s.Golazos.Collection",
			contracts.GolazosAddress.Hex()),
			nftCollectionDataView.PublicLinkedType)

		traitsView := getMomentNFTTraitsMetadataView(
			t,
			b,
			contracts,
			userAddress,
			nftID,
		)
		assert.Contains(t, traitsView, TraitView{
			Name: metadatakey1, Value: playNames[editions[editionID].PlayID-1],
		})
		assert.Contains(t, traitsView, TraitView{
			Name: metadatakey2, Value: playerJerseyName,
		})
		assert.Contains(t, traitsView, TraitView{
			Name: metadatakey3, Value: playType,
		})
		assert.Contains(t, traitsView, TraitView{
			Name: "editionTier", Value: editions[editionID].Tier,
		})
		assert.Contains(t, traitsView, TraitView{
			Name: "seriesName", Value: series[editions[editionID].SeriesID].Name,
		})
		assert.Contains(t, traitsView, TraitView{
			Name: "setName", Value: sets[editions[editionID].SetID].Name,
		})
		assert.Contains(t, traitsView, TraitView{
			Name: "serialNumber", Value: fmt.Sprintf("%v", nftProperties.SerialNumber),
		})
	} else {
		assert.Equal(t, previousSupply, newSupply)
	}
}
