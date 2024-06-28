# NFT Provider Aggregator

`NFTProviderAggregator` is a general-purpose Cadence contract for aggregating multiple NFT providers into a single provider capability conforming to the [`NonFungibleToken`](https://github.com/onflow/flow-nft/blob/master/contracts/NonFungibleToken.cdc#L107) standard. It makes it possible to withdraw NFTs scattered across multiple collections without having to move NFTs into a single collection. See a [high level diagram of expected usage](https://sketchboard.me/BDDABYYoORLI#/).

> Note: Whenever possible, the collection that holds a given NFT to withdraw should be determined off-chain so that the activity on-chain is minimized to the withdrawal only instead of having first to spend gas retrieving which collection the NFT is stored in. The `NFTProviderAggregator` contract may be helpful in situations such as when a third-party contract is intended to withdraw NFTs using a single provider capability and that the NFTs are scattered across multiple collections.

## Usage

There are two types of accounts:
 - **Manager**: An account holding an `Aggregator` resource - any account can create `Aggregator` resources.
 - **Supplier**: An account holding a `Supplier` resource created using a capability from the parent `Aggregator` resource.

 The **manager** has access to the aggregated NFT provider and can add/remove any NFT provider capability. **Suppliers** can remove only NFT provider capabilities they added themselves.

Setup steps:
 1. `bootstrap_aggregator.cdc`: Create an `Aggregator` resource, save it in the **manager** account's storage, and publish a *SupplierFactory* capability for each designated **supplier**.
 2. `bootstrap_supplier.cdc`: Claim the *SupplierFactory* capability, create a `Supplier` resource, and save it in the **supplier** account's storage (repeat for each **supplier**).
 3. `add_nft_provider_capability_as_supplier.cdc`: Add a NFT provider capability (repeat as needed for each **supplier** and each collection) - the transaction may be merged with that of step 2. Only NFT provider capabilities targeting collections of valid NFT type can be added (i.e., the type defined when the `Aggregator` resource is created).
 
Once the setup steps are completed, NFTs scattered across the multiple collections added to the `Aggregator` resource can be withdrawn:
- By the **manager** with `transfer_from_aggregated_NFT_provider_as_manager.cdc`.
- By a designated **third-party account** with `transfer_from_aggregated_NFT_provider_as_thirdparty.cdc` after the **aggregated provider capability** being published by the **manager** with `publish_aggregated_nft_provider_capability.cdc` and claimed by the **third-party** with `claim_aggregated_nft_provider_capability.cdc`. 

NFT provider capabilities should be removed when they are not needed anymore with `remove_nft_provider_capability_as_supplier.cdc`. Destroying a `Supplier` resource removes all the NFT provider capabilities it previously added to the parent `Aggregator` resource. Destroying an `Aggregator` resource
removes all the resource's NFT provider capabilities and render child `Supplier` resources inoperable, they should be destroyed too.

## Architecture Considerations

> July 2024 Update: NFT Provider Aggregator was created prior to Cadence 1.0 and then updated for Cadence 1.0 with the goal of ensuring backwards compatibility. Even though Cadence 1.0 introduced relevant new features that  (e.g., the ability to add tags to keep track of capability usage), this current version of NFT Provider Aggregator doesn't have any architecture changes compared to the pre-Cadence 1.0 version.

The use of the `Supplier` resource allows:
- Explicitly keeping track of the NFT provider capabilities that are expected to be valid all in one place, the parent `Aggregator` resource, and emitting dedicated events when a capability is added or removed.
- Reversibly exposing NFT provider capabilities to the parent `Aggregator` resource without the capabilities being retrievable individually by the **manager**. This is made possible by the `Aggregator` resource’s `nftWithdrawCapabilities` dictionary being defined with `access(self)` access control and devoid of any getter function, though, if the `NFTProviderAggregator` contract is updatable, a contract update may potentially change that.

> Note: Another architecture to aggregate NFT providers exists where **suppliers** link NFT provider capabilities at custom private paths and directly publish them to the **manager** holding the `Aggregator` resource, thus bypassing the need for the `Supplier` resource. It should be noted that:
> - Off-chain records may be needed to keep track of each private path where a NFT provider capability was linked, even if it was then unlinked by the **supplier** to revoke **manager** access because it is important to be aware that the **manager** would regain access if it is linked at the same path again later.
> - Supplied NFT provider capabilities are not necessarily confined to the `Aggregator` resource because the **manager** can copy and save them somewhere else as part of the claim transaction.

## Transactions

The transactions included in this repository and listed below are all single-signer transactions.

#### Resource Bootstrapping:

- `bootstrap_aggregator.cdc`

- `bootstrap_supplier.cdc`

#### NFT Provider Capability Addition/Removal:

- `add_nft_provider_capability_as_supplier.cdc`

- `remove_nft_provider_capability_as_supplier.cdc`

> Admin Operations:
> - `add_nft_provider_capability_as_manager.cdc`
> - `remove_nft_provider_capability_as_manager.cdc`

#### NFT Withdrawal:

- `transfer_from_aggregated_nft_provider_as_manager.cdc`

- `transfer_from_aggregated_nft_provider_as_thirdparty.cdc`

#### Capability Distribution:

- `publish_additional_supplier_factory_capabilities.cdc`

- `publish_aggregated_nft_provider_capability.cdc`

- `claim_aggregated_nft_provider_capability.cdc`

- `unpublish_capability.cdc`

#### Resource Destruction:

- `destroy_aggregator.cdc`

- `destroy_supplier.cdc`

## Scripts

The scripts included in this repository and listed below are scripts that access the `Supplier` resource's getter functions using the `SupplierPublic` interface. Each script has one argument, a supplier address.

- `get_aggregator_uuid.cdc`

- `get_collection_uuids.cdc`

- `get_ids.cdc`

- `get_supplier_added_collection_uuids.cdc`

## Run Contract Tests

```
git clone git@github.com:dapperlabs/nft-provider-aggregator.git
cd ./nft-provider-aggregator/test
npm install
npm test
```
