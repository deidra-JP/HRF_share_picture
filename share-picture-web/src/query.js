/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */



'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');


export async function Main() {
    try {
        // load the network configuration
        const connectionProfileJson = (await fs.promises.readFile('/home/deidra/share-picture/sp-network/organizations/peerOrganizations/org1.example.com/connection-org1.json')).toString();
        const connectionProfile = JSON.parse(connectionProfileJson);
       
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        await wallet.get('Org1');

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(connectionProfile, { wallet, identity: 'Org1'});

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('spmainchannel');

        // Get the contract from the network.
        const contract = network.getContract('sp-logic');

        // Evaluate the specified transaction.
        // queryPicture transaction - requires 1 argument, ex: ('queryPicture', 'PICTURE4')
        // queryAllPictures transaction - requires no arguments, ex: ('queryAllPictures')
        const result = await contract.evaluateTransaction('queryAllPictures');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

        // Disconnect from the gateway.
        await gateway.disconnect();
        
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}
