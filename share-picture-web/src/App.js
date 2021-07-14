/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
'use strict';


import logo from './logo.svg';
import './App.css';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');


async function query() {
  try {
      // load the network configuration
      const ccpPath = path.resolve(__dirname, '..', '..', 'sp-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
      const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

      // Create a new file system based wallet for managing identities.
      const walletPath = path.join(process.cwd(), 'wallet');
      const wallet = await Wallets.newFileSystemWallet(walletPath);
      console.log(`Wallet path: ${walletPath}`);

      // Check to see if we've already enrolled the user.
      const identity = await wallet.get('appUser');
      if (!identity) {
          console.log('An identity for the user "appUser" does not exist in the wallet');
          console.log('Run the registerUser.js application before retrying');
          return;
      }

      // Create a new gateway for connecting to our peer node.
      const gateway = new Gateway();
      await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

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

query();


function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
      <div>
        Share Picture 　操作一覧
      </div>
      <div>
        <button type="button" onClick={this.query}>query</button>
      </div>
    </div>
  );
}

export default App;
