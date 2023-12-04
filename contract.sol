// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

/**
 * @title Storage
 * @dev Store & retrieve value in a variable
 * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
 */
contract Storage {

    string private storedString;

    /**
     * @dev Store value in variable
     * @param str value to store
     */
    function store(string memory str) public {
        storedString = str;
    }

    /**
     * @dev Return value 
     * @return value of 'storedString'
     */
    function retrieve() public view returns (string memory){
        return storedString;
    }
}

