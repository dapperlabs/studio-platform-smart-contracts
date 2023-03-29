package test

import (
	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEPLDeployContracts(t *testing.T) {
	b := newEmulator()
	EPLDeployContracts(t, b)
}

func TestEPLSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupEPLAccount(t, b, userAddress, userSigner, contracts)

	t.Run("Account should be set up", func(t *testing.T) {
		isAccountSetUp := isAccountSetup(
			t,
			b,
			contracts,
			userAddress,
		)
		assert.Equal(t, true, isAccountSetUp)
	})
}

func TestSeries(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
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
		series := getSeriesData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, series.ID)
		assert.Equal(t, seriesName, series.Name)
		assert.Equal(t, true, series.Active)
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

func TestSets(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
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

func TestTags(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
	createTestTag(t, b, contracts)
}

func createTestTag(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new tag", func(t *testing.T) {
		testCreateTag(
			t,
			b,
			contracts,
			"Tag One",
			1,
			false,
		)
	})

	t.Run("Should be able to create a new tag with an incrementing ID from the first", func(t *testing.T) {
		testCreateTag(
			t,
			b,
			contracts,
			"Tag Two",
			2,
			false,
		)
	})
}

func testCreateTag(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	tagName string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createTag(
		t,
		b,
		contracts,
		tagName,
		false,
	)

	if !shouldRevert {
		tag := getTagData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, tag.ID)
		assert.Equal(t, tagName, tag.Name)
	}
}

func TestPlays(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
	createTestPlays(t, b, contracts)
}

func createTestPlays(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new play", func(t *testing.T) {
		tagIds := []uint64{1}
		testCreatePlay(
			t,
			b,
			contracts,
			map[string]string{"key1": "Erling", "key2": "Haaland"},
			tagIds,
			1,
			false,
		)
	})

	t.Run("Should be able to create a new play with an incrementing ID from the first", func(t *testing.T) {
		tagIds := []uint64{2}
		testCreatePlay(
			t,
			b,
			contracts,
			map[string]string{"key1": "Erling", "key2": "Haaland"},
			tagIds,
			2,
			false,
		)
	})
}

func testCreatePlay(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	metadata map[string]string,
	tagIds []uint64,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createTag(
		t,
		b,
		contracts,
		"test tag",
		false,
	)
	createPlay(
		t,
		b,
		contracts,
		metadata,
		tagIds,
		false,
	)

	if !shouldRevert {
		play := getPlayData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, play.ID)
		assert.Equal(t, tagIds, play.TagIds)
		assert.Equal(t, metadata, play.Metadata)
	}
}

func TestEditions(t *testing.T) {
	b := newEmulator()
	contracts := EPLDeployContracts(t, b)
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
		//if maxMintSize != nil {
		//	assert.Equal(t, &maxMintSize, &edition.MaxMintSize)
		//}
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
