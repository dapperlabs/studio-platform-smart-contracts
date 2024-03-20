// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Capped.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract Kitty is ERC20, Ownable, ERC20Capped, ReentrancyGuard {
    IERC721 public cryptoKitties;
    mapping(uint256 => bool) private katapultedKitties;
    uint256 public constant KITTYLOAF = 1000 * 10**18; // 1000 tokens, considering 18 decimals
    event Katapult(address indexed from, string flowAddress, uint256[] kittyIds);
    event Meow(address indexed from, string flowAddress, uint256[] kittyIds);


    constructor(address _cryptoKittiesAddress)
        ERC20("Kitty", "MEOW")
        ERC20Capped(888_000_000_000 *  1e18) // 888 billion tokens with 18 decimals
    {
        require(_cryptoKittiesAddress != address(0), "Invalid CryptoKitties address");
        cryptoKitties = IERC721(_cryptoKittiesAddress);
    }

    function mint(address to, uint256 amount) public onlyOwner {
        // Ensure total supply + mint does not exceed cap
        require(totalSupply() + amount <= cap(), "Cap exceeded");
        _mint(to, amount);
    }

    function katapult(uint256[] calldata kittyIds, string calldata flowAddress) external nonReentrant {
        require(kittyIds.length > 0, "No kitty IDs provided");
        uint256 rewardAmount = KITTYLOAF * kittyIds.length;
        require(totalSupply() + rewardAmount <= cap(), "Cap exceeded");

        for(uint i = 0; i < kittyIds.length; i++) {
            require(cryptoKitties.ownerOf(kittyIds[i]) == msg.sender, "Caller is not the owner");
            require(!katapultedKitties[kittyIds[i]], "Kitty has already been katapulted");

            katapultedKitties[kittyIds[i]] = true;
        }

        _mint(msg.sender, rewardAmount);
        emit Katapult(msg.sender, flowAddress, kittyIds);
    }

    function meow(string calldata message) external nonReentrant {
        require(balanceOf(msg.sender) >= KITTYLOAF * 1e18, "Kittyloaf required");

        emit Meow(message);
    }
}