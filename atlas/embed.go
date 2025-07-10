package atlas

import (
	_ "embed"
)

// Transactions is a list of all the transactions we export with imports mapped
var (
	//go:embed transactions/user/buy_packs_primary_sale.cdc
	UserBuyPacksPrimarySale []byte

	//go:embed transactions/admin/fulfill_pack_buyback_offer.cdc
	AdminFulfillPackBuybackOffer []byte
)
