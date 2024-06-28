import path from "path";
import { 
  emulator, 
  init, 
  getAccountAddress,
  deployContract,
  deployContractByName, 
  sendTransaction, 
  shallPass,
  shallRevert, 
  executeScript 
} from "@onflow/flow-js-testing";
import fs from "fs";

// Increase timeout if tests failing due to timeout
// jest.setTimeout(10000);

// Set basepath of the project
const basePath = path.resolve(__dirname, "./../");

// Define the test suite for NFTProviderAggregator contract
describe("NFTProviderAggregator Contract Tests", () => {

  // Variables for holding the account addresses
  let serviceAccount;
  let manager;
  let supplier;
  let supplierTwo;
  let thirdParty;
  let nftTypeIdentifier;

  // Setup each test
  beforeEach(async () => {
	// Set logging flag to true will pipe emulator output to console
    const logging = false;

    // Initialize flow-js-testing
    await init(basePath);

    // Start the emulator
    await emulator.start({ logging });

    // Get account addresses
    serviceAccount = await getAccountAddress("ServiceAccount");
    manager = await getAccountAddress("Manager");
    supplier = await getAccountAddress("Supplier");
    supplierTwo = await getAccountAddress("SupplierTwo");
    thirdParty = await getAccountAddress("ThirdParty");
    nftTypeIdentifier = "A.".concat(serviceAccount.replace("0x", "")).concat(".ExampleNFT.Collection")

    // Deploy project contracts with service account
    await shallPass(
      deployContractByName({ to: serviceAccount, name: "NonFungibleToken"})
    );
    await shallPass(
      deployContractByName({ to: serviceAccount, name: "ExampleNFT"})
    );
    await shallPass(
      deployContractByName({ to: serviceAccount, name: "NFTProviderAggregator"})
    );
  
    // Setup exampleNFT Collection and mint example NFTs in manager and supplier's collections
    await sendTransaction({
      name: "exampleNFT/setup_exampleNFT",
      args: [],
      signers: [manager]
    });
    await sendTransaction({
      name: "exampleNFT/setup_exampleNFT",
      args: [],
      signers: [supplier]
    });
    await sendTransaction({
      name: "exampleNFT/setup_exampleNFT",
      args: [],
      signers: [supplierTwo]
    });
    await sendTransaction({
      name: "exampleNFT/setup_exampleNFT",
      args: [],
      signers: [thirdParty]
    });
    await sendTransaction({
      name: "exampleNFT/mint_exampleNFTBatched",
      args: [manager, "1"],
      signers: [serviceAccount]
    });
    await sendTransaction({
      name: "exampleNFT/mint_exampleNFTBatched",
      args: [supplier, "1"],
      signers: [serviceAccount]
    });

  });

  // Stop the emulator after each test (so it could be restarted)
  afterEach(async () => {
    return emulator.stop();
  });
  
  test("#1: should be able to bootstrap an Aggregator resource", async () => {
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
  });

  test("#2: should be able to add a NFT provider capability as manager", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add a NFT provider capability as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
  });
  
  test("#3: should NOT be able to add a NFT provider capability that is already existing as manager", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add a NFT provider capability as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Attempt adding the same NFT provider capability again
    const [_, error] = await shallRevert(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    const expected = expect.stringContaining("panic: NFT provider capability already exists!");
    expect(error).toEqual(expected);
  });

  test("#4: should be able to bootstrap Supplier resources", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [manager, supplier], ["SupplierAccess", "SupplierAccess2"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess2"],
        signers: [supplier]
      })
    );
    // Step 4: Publish additional supplier capabilities signing as manager
    await shallPass(
      sendTransaction({
        name: "publish_additional_supplier_factory_capabilities",
        args: [[supplierTwo], ["SupplierAccess3"]],
        signers: [manager]
      })
    );
    // Step 5: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess3"],
        signers: [supplierTwo]
      })
    );
  });

  test("#5: should be able to add a NFT provider capability as supplier", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );

    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
  });

  test("#6: should NOT be able to add a NFT provider capability that is already existing as supplier", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 4: Attempt adding the same NFT provider capability again
    const [_, error] = await shallRevert(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    const expected = expect.stringContaining("panic: NFT provider capability already exists!");
    expect(error).toEqual(expected);
  });

  test("#7: should NOT be able to add a NFT provider capability that targets a collection with invalid NFT type", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], false],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Deploy an alternative NFT contract (AltExampleNFT) to supplier
    const altExampleNFT_contract_code = fs.readFileSync(path.resolve(basePath, "./test/cadence/contracts/AltExampleNFT.cdc"), {encoding:'utf8', flag:'r'});
    await shallPass(
      deployContract({
        to: supplier,
        name: "AltExampleNFT",
        code: altExampleNFT_contract_code
      })
    );
    // Step 4: Mint an NFT to supplier so that the collection is not empty
    await shallPass(
      sendTransaction({
        name: "./../test/cadence/transactions/mint_altExampleNFTBatched",
        args: [supplier, "1"],
        signers: [supplier]
      })
    );
    // Step 5: Attempt adding a NFT provider capability signing as supplier
    const [_r, error] = await shallRevert(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["altExampleNFTProvider", "altExampleNFTCollection"],
        signers: [supplier]
      })
    );
    const expected = expect.stringContaining("assertion failed: NFT provider capability targets a collection with invalid NFT type!");
    expect(error).toEqual(expected);
  });

  test("#8: should be able to withdraw NFTs from Aggregator's aggregated provider held in both the supplier and the manager's own collections", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
    });
    // Step 6: Check that array of NFT ids returned by Aggregator is the concatenation of result1[0] and result2[0]]
    const [result3, _e3] = await executeScript({
      name: "get_ids",
      args: [supplier]
    });
    expect(result3.length).toBe(2);
    expect(result3).toContain(result1[0]);
    expect(result3).toContain(result2[0]);
    // Step 7: Transfer NFT with ID = result1[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
        name: "transfer_from_aggregated_nft_provider_as_manager",
        args: [serviceAccount, result1[0]],
        signers: [manager]
      })
    );
    // Step 8: Transfer NFT with ID = result2[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result2[0]],
          signers: [manager]
      })
      );
    // Step 9: Check if service account has NFT with ID = result1[0] and NFT with ID = result2[0]
    const [result4, _e4] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
    });
    expect(result4.length).toBe(2);
    expect(result4).toContain(result1[0]);
    expect(result4).toContain(result2[0]);
  });

  test("#9: should be able to withdraw NFTs from Aggregator's aggregated provider even if manager's own collection is empty but supplier's is not", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    // Step 6: Transfer NFT with ID = result1[0] from manager's own provider to service account and verify the manager's collection is empty
    await shallPass(
      sendTransaction({
        name: "exampleNFT/transfer_exampleNFT",
        args: [serviceAccount, result1[0]],
        signers: [manager]
      })
    );
    const [result3, _e3] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    expect(result3).toEqual([]);
    // Step 7: Transfer NFT with ID = result2[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result2[0]],
          signers: [manager]
      })
      );
    // Step 8: Check if service account has NFT with ID = result1[0] and NFT with ID = result2[0]
    const [result4, _e4] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
    });
    expect(result4.length).toBe(2);
    expect(result4).toContain(result1[0]);
    expect(result4).toContain(result2[0]);
  });

  test("#10: should be able to withdraw NFTs from Aggregator's aggregated provider even if manager's capability gets unlinked but supplier's does not", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    // Step 6: Unlink manager's provider capability
    await shallPass(
      sendTransaction({
        name: "./../test/cadence/transactions/unlink_nft_provider",
        args: [],
        signers: [manager]
      })
    );
    // Step 7: Attempt transferring NFT with ID = result1[0] from manager's aggregated provider to service account
    const [_r3, error3] = await shallRevert(
      sendTransaction({
        name: "transfer_from_aggregated_nft_provider_as_manager",
        args: [serviceAccount, result1[0]],
        signers: [manager]
      })
    );
    const expected = expect.stringContaining("panic: missing NFT");
    expect(error3).toEqual(expected);
    // Step 8: Check that array of NFT ids returned by Aggregator is result2[0]]
    const [result4, _e4] = await executeScript({
      name: "get_ids",
      args: [supplier]
    });
    expect(result4).toEqual([result2[0]]);
    // Step 9: Transfer NFT with ID = result2[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result2[0]],
          signers: [manager]
      })
      );
    // Step 10: Check if service account has NFT with ID = result2[0]
    const [result5, _e5] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
    });
    expect(result5).toEqual([result2[0]]);
  });

  test("#11: should be able to withdraw NFTs from Aggregator's aggregated provider even if supplier's own collection is empty but manager's is not", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    // Step 6: Transfer NFT with ID = result1[0] from supplier's own provider to service account and verify the supplier's collection is empty
    await shallPass(
      sendTransaction({
        name: "exampleNFT/transfer_exampleNFT",
        args: [serviceAccount, result2[0]],
        signers: [supplier]
      })
    );
    const [result3, _e3] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
    });
    expect(result3).toEqual([]);
    // Step 7: Transfer NFT with ID = result2[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result1[0]],
          signers: [manager]
      })
      );
    // Step 8: Check if service account has NFT with ID = result1[0] and NFT with ID = result2[0]
    const [result4, _e4] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
    });
    expect(result4.length).toBe(2);
    expect(result4).toContain(result1[0]);
    expect(result4).toContain(result2[0]);
  });

  test("#12: should be able to withdraw NFTs from Aggregator's aggregated provider even if supplier's capability gets unlinked but manager's does not", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    // Step 6: Unlink supplier's NFT provider capability
    await shallPass(
      sendTransaction({
        name: "./../test/cadence/transactions/unlink_nft_provider",
        args: [],
        signers: [supplier]
      })
    );
    // Step 7: Attempt transferring NFT with ID = result2[0] from manager's aggregated provider to service account
    const [_r3, error3] = await shallRevert(
      sendTransaction({
        name: "transfer_from_aggregated_nft_provider_as_manager",
        args: [serviceAccount, result2[0]],
        signers: [manager]
      })
    );
    const expected = expect.stringContaining("panic: missing NFT");
    expect(error3).toEqual(expected);
    // Step 8: Check that array of NFT ids returned by Aggregator is result1[0]]
    const [result4, _e4] = await executeScript({
      name: "get_ids",
      args: [supplier]
    });
    expect(result4).toEqual([result1[0]]);
    // Step 8: Transfer NFT with ID = result1[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result1[0]],
          signers: [manager]
      })
      );
    // Step 9: Check if service account has NFT with ID = result2[0]
    const [result5, _e5] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
    });
    expect(result5).toEqual([result1[0]]);
  });

  test("#13: should NOT be able to withdraw NFTs from Aggregator's aggregated provider if both the supplier and the manager's own collections are empty", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Add NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_manager",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [manager]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 4: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 5: Get NFT Ids in manager and supplier's collections
    const [result1, _e1] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    // Step 6: Transfer NFTs with ID = result1[0] and ID = result2[0] from manager and supplier's own providers
    // to service account and verify both the manager and supplier's collections are empty
    await shallPass(
      sendTransaction({
        name: "exampleNFT/transfer_exampleNFT",
        args: [serviceAccount, result1[0]],
        signers: [manager]
      })
    );
    const [result3, _e3] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [manager]
    });
    expect(result3).toEqual([]);
    await shallPass(
      sendTransaction({
        name: "exampleNFT/transfer_exampleNFT",
        args: [serviceAccount, result2[0]],
        signers: [supplier]
      })
    );
    const [result4, _e4] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
    });
    expect(result4).toEqual([]);
    // Step 7: Attempt withdrawing NFT with ID = result1[0] from manager's aggregated provider to service account
    const [_r1, error1] = await shallRevert(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result1[0]],
          signers: [manager]
      })
    );
    const expected1 = expect.stringContaining("error: panic: missing NFT");
    expect(error1).toEqual(expected1);
    // Step 8: Attempt withdrawing NFT with ID = result2[0] from manager's aggregated provider to service account
    const [_r2, error2] = await shallRevert(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result2[0]],
          signers: [manager]
      })
    );
    const expected2 = expect.stringContaining("error: panic: missing NFT");
    expect(error2).toEqual(expected2);
  });
  
  test("#14: should be able to remove a NFT provider capability as supplier previously added by themselves", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 4: Remove NFT provider capability signing as supplier
    const [result, _e] = await executeScript({
      name: "get_supplier_added_collection_uuids",
      args: [supplier]
    });
    await shallPass(
      sendTransaction({
        name: "remove_nft_provider_capability_as_supplier",
        args: [result[0]],
        signers: [supplier]
      })
    );
  });

  test("#15: should be able to remove a NFT provider capability added by supplier as manager", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 4: Remove NFT provider capability signing as manager
    const [result, _e] = await executeScript({
      name: "get_collection_uuids",
      args: [supplier]
    });
    await shallPass(
      sendTransaction({
        name: "remove_nft_provider_capability_as_manager",
        args: [result[0]],
        signers: [manager]
      })
    );
  });

  test("#16: should NOT be able to remove a NFT provider capability added by a separate supplier as supplier", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier, supplierTwo], ["SupplierAccess", "AggregatorAccessTwo"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Bootstrap a Supplier resource signing as supplierTwo depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "AggregatorAccessTwo"],
        signers: [supplierTwo]
      })
    );
    // Step 4: Mint a NFT to supplierTwo's collection so that it is not empty and add NFT provider capability signing as supplierTwo
    await sendTransaction({
      name: "exampleNFT/mint_exampleNFTBatched",
      args: [supplierTwo, "1"],
      signers: [serviceAccount]
    });
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplierTwo]
      })
    );
    // Step 5: Check cannot remove NFT provider capability added by supplierTwo signing as supplier
    const [result1, _e] = await executeScript({
      name: "get_supplier_added_collection_uuids",
      args: [supplierTwo]
    });
    const [_r, error] = await shallRevert(
      sendTransaction({
        name: "remove_nft_provider_capability_as_supplier",
        args: [result1[0]],
        signers: [supplier]
      })
    );
    const expected = expect.stringContaining("error: pre-condition failed: Collection UUID does not exist in added collection UUIDs!");
    expect(error).toEqual(expected);
  });

  test("#17: should NOT be able to withdraw NFTs in supplier's collection using Aggregator's aggregated provider if supplier removed their NFT provider capability", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
      sendTransaction({
        name: "bootstrap_supplier",
        args: [manager, "SupplierAccess"],
        signers: [supplier]
      })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProvider", "exampleNFTCollection"],
        signers: [supplier]
      })
    );
    // Step 4: Remove NFT provider capability signing as supplier
    const [result1, _e1] = await executeScript({
      name: "get_supplier_added_collection_uuids",
      args: [supplier]
    });
    await shallPass(
      sendTransaction({
        name: "remove_nft_provider_capability_as_supplier",
        args: [result1[0]],
        signers: [supplier]
      })
    );
    // Step 5: Get a NFT Id in supplier's collection
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [supplier]
      });
    expect(result2.length).toBe(1);
    // Step 6: Attempt withdrawing NFT with ID = result[0] from manager's aggregated provider to service account
    const [_r3, error3] = await shallRevert(
      sendTransaction({
          name: "transfer_from_aggregated_nft_provider_as_manager",
          args: [serviceAccount, result2[0]],
          signers: [manager]
      })
    );
    const expected = expect.stringContaining("error: panic: missing NFT");
    expect(error3).toEqual(expected);
  });

  test("#18: should be able to publish the manager's aggregated NFT provider capability", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
      sendTransaction({
        name: "bootstrap_aggregator",
        args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
        signers: [manager]
      })
    );
    // Step 2: Publish a private capability to third party account inbox
    await shallPass(
      sendTransaction({
        name: "publish_aggregated_nft_provider_capability",
        args: [thirdParty, "AggregatedNFTProviderCapability"],
        signers: [manager]
      })
    );
  });

  test("#19: should be able to claim the manager's aggregated NFT provider capability and withdraw from it", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
    sendTransaction({
      name: "bootstrap_aggregator",
      args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
      signers: [manager]
    })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
    sendTransaction({
      name: "bootstrap_supplier",
      args: [manager, "SupplierAccess"],
      signers: [supplier]
    })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
    sendTransaction({
      name: "add_NFT_provider_capability_as_supplier",
      args: ["exampleNFTProvider", "exampleNFTCollection"],
      signers: [supplier]
    })
    );
    const [result1, _e1] = await executeScript({
    name: "get_ids",
    args: [supplier]
    });
    // Step 4: Publish a private capability to third party account inbox
    await shallPass(
      sendTransaction({
        name: "publish_aggregated_nft_provider_capability",
        args: [thirdParty, "AggregatedNFTProviderCapability"],
        signers: [manager]
      })
    );
    // Step 5: Claim capability from manager
    await shallPass(
      sendTransaction({
        name: "claim_aggregated_nft_provider_capability",
        args: [manager, "AggregatedNFTProviderCapability"],
        signers: [thirdParty]
      })
    );
    // Step 6: Check serviceAccount's collection is empty
    const [result2, _e2] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
      });
    expect(result2.length).toBe(0);
    // Step 7: Transfer NFT with ID = result[0] from manager's aggregated provider to service account
    await shallPass(
      sendTransaction({
        name: "transfer_from_aggregated_nft_provider_as_thirdparty",
        args: [serviceAccount, result1[0]],
        signers: [thirdParty]
      })
    );
    // Step 8: Check serviceAccount's collection has 1 NFT
    const [result3, _e3] = await executeScript({
      name: "exampleNFT/balance_exampleNFT",
      args: [serviceAccount]
      });
    expect(result3.length).toBe(1);
  });

  test("#20: should be able to remove supplied NFT provider capabilities when a Supplier resource is destroyed", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
    sendTransaction({
      name: "bootstrap_aggregator",
      args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
      signers: [manager]
    })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
    sendTransaction({
      name: "bootstrap_supplier",
      args: [manager, "SupplierAccess"],
      signers: [supplier]
    })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
    sendTransaction({
      name: "add_NFT_provider_capability_as_supplier",
      args: ["exampleNFTProvider", "exampleNFTCollection"],
      signers: [supplier]
    })
    );
    // Step 4: Add a second NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "./../test/cadence/transactions/setup_exampleNFT",
        args: ["exampleNFTCollectionTwo", "exampleNFTCollectionTwo"],
        signers: [supplier]
      })
    );
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProviderTwo", "exampleNFTCollectionTwo"],
        signers: [supplier]
      })
    );
    // Step 5: Add a third NFT provider capability signing as supplier
    await shallPass(
      sendTransaction({
        name: "./../test/cadence/transactions/setup_exampleNFT",
        args: ["exampleNFTCollectionThree", "exampleNFTCollectionThree"],
        signers: [supplier]
      })
    );
    await shallPass(
      sendTransaction({
        name: "add_NFT_provider_capability_as_supplier",
        args: ["exampleNFTProviderThree", "exampleNFTCollectionThree"],
        signers: [supplier]
      })
    );
    const [result1, _e1] = await executeScript({
      name: "get_collection_uuids",
      args: [supplier]
      });
    expect(result1.length).toBe(3);
    const [result2, _e2] = await executeScript({
      name: "get_supplier_added_collection_uuids",
      args: [supplier]
      });
    expect(result2.length).toBe(3);
    // Step 6: Remove a NFT provider capability signing as manager
    await shallPass(
      sendTransaction({
        name: "remove_NFT_provider_capability_as_manager",
        args: [result1[0]],
        signers: [manager]
      })
    );
    // Step 7: Check that the NFT capability was removed from the Aggregator but that the
    // length of supplierAddedCollectionUUIDs hasn't changed
    const [result3, _e3] = await executeScript({
      name: "get_collection_uuids",
      args: [supplier]
      });
    expect(result3.length).toBe(2);
    const [result4, _e4] = await executeScript({
      name: "get_supplier_added_collection_uuids",
      args: [supplier]
      });
    expect(result4.length).toBe(3);
    // Step 8: Destroy the supplier resource even if a NFT provider capability was removed by the manager
    await shallPass(
      sendTransaction({
        name: "destroy_supplier",
        args: [],
        signers: [supplier]
      })
    );
  });

  test("#21: should be able to nullify the aggregated NFT provider and child Supplier resources when the Aggregator resource is destroyed", async () => {
    // Step 1: Bootstrap an Aggregator resource signing as manager depositing to manager
    await shallPass(
    sendTransaction({
      name: "bootstrap_aggregator",
      args: [nftTypeIdentifier, [supplier], ["SupplierAccess"], true],
      signers: [manager]
    })
    );
    // Step 2: Bootstrap a Supplier resource signing as supplier depositing to supplier
    await shallPass(
    sendTransaction({
      name: "bootstrap_supplier",
      args: [manager, "SupplierAccess"],
      signers: [supplier]
    })
    );
    // Step 3: Add NFT provider capability signing as supplier
    await shallPass(
    sendTransaction({
      name: "add_NFT_provider_capability_as_supplier",
      args: ["exampleNFTProvider", "exampleNFTCollection"],
      signers: [supplier]
    })
    );
    const [result, _e] = await executeScript({
    name: "get_ids",
    args: [supplier]
    });
    // Step 4: Publish a private capability to third party account inbox
    await shallPass(
      sendTransaction({
        name: "publish_aggregated_nft_provider_capability",
        args: [thirdParty, "AggregatedNFTProviderCapability"],
        signers: [manager]
      })
    );
    // Step 5: Claim capability from manager
    await shallPass(
      sendTransaction({
        name: "claim_aggregated_nft_provider_capability",
        args: [manager, "AggregatedNFTProviderCapability"],
        signers: [thirdParty]
      })
    );
    // Step 6: Destroy the aggregator resource
    await shallPass(
      sendTransaction({
        name: "destroy_aggregator",
        args: [],
        signers: [manager]
      })
      );
    // Step 7: Attempt transfer NFT with ID = result[0] from manager's aggregated provider to service account
    const [_, error] =await shallRevert(
      sendTransaction({
        name: "transfer_from_aggregated_nft_provider_as_thirdparty",
        args: [serviceAccount, result[0]],
        signers: [thirdParty]
      })
    );
    const expected = expect.stringContaining("panic: Could not get capability and borrow reference");
    expect(error).toEqual(expected);
  });
});
 