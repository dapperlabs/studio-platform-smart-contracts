package atlas

import (
	_ "embed"
)

// Transactions is a list of all the transactions we export with imports mapped
var (
	UserBuyPacksPrimarySale []byte
	//go:embed transactions/user/buy_packs_primary_sale.cdc
)
