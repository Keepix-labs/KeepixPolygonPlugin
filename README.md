# Keepix Polygon Plugin

## Contribute

### Prerequisites

- Golang 1.21+
- Node 18+

### Install

`npm install`

### Build

`npm run build`

The built executables are found in /dist folder.

### Run frontend

`npm run start`  

Go on http://localhost:3000
Also an api mock of keepix-server are available on http://localhost:2000/plugins/keepix-polygon-plugin/status like routes where you can see on keepix-server.  

### Plugin Front-end

The plugin need a front-end static code in the final build directory index.html file  
Here we are using a React framework  
The Front-end application will be loaded by the Keepix with an iframe at the following endpoint url:  
  
`http|https://hostname/plugins/keepix-polygon-plugin/view`  
  
### Plugin Dev Api from Front-end dev
  
For developping locally your plugin on the config-overrides.js you can see a  
express.js running on 0.0.0.0:2000 is a simulation server copying routes of the real keepix server.  

`GET /plugins/nameOfThePlugin/:key`  
`POST /plugins/nameOfThePlugin/:key`  
`GET /plugins/nameOfThePlugin/watch/tasks/:taskId`  

### Integration data:
"platforms": ["linux-arm64", "linux-x64", "win-x64", "oxs-x64", "oxs-arm64"],
"ports": [
{
    "internal": 30303,
    "external": 30303,
    "protocol": "TCP",
    "description": "keepix-upnp-erigon-tcp"
},
{
    "internal": 30303,
    "external": 30303,
    "protocol": "UDP",
    "description": "keepix-upnp-erigon-udp"
},
{
    "internal": 30304,
    "external": 30304,
    "protocol": "TCP",
    "description": "keepix-upnp-erigon-tcp2"
},
{
    "internal": 30304,
    "external": 30304,
    "protocol": "UDP",
    "description": "keepix-upnp-erigon-udp2"
},
{
    "internal": 42069,
    "external": 42069,
    "protocol": "TCP",
    "description": "keepix-erigon-snapshot-tcp"
},
{
    "internal": 42069,
    "external": 42069,
    "protocol": "UDP",
    "description": "keepix-erigon-snapshot-udp"
},
{
    "internal": 26656,
    "external": 26656,
    "protocol": "TCP",
    "description": "keepix-heimdall-p2p-tcp"
},
{
    "internal": 26656,
    "external": 26656,
    "protocol": "UDP",
    "description": "keepix-heimdall-p2p-udp"
}
],
"installForm": {
    "inputs": [
        {
            "key": "mnemonic",
            "type": "wallet",
            "walletType": "ethereum",
            "walletKey": "mnemonic",
            "filter": "wallet.asMnemonic === true",
            "label": "Select a Polygon Wallet.",
            "labelIfEmptyWalletList": "Please go on wallets page and create or import (with mnemomic) a Polygon Wallet."
        },
        {
            "key": "testnet",
            "type": "checkbox-2",
            "true": "true",
            "false": "false",
            "defaultValue": false,
            "label": "Mainnet or Testnet (mumbai)"
        },
        {
            "key": "ethereumRPC",
            "type": "string",
            "label": "Ethereum mainnet rpc (optional)"
        },
        {
            "key": "autostart",
            "type": "checkbox-2",
            "true": "true",
            "false": "false",
            "defaultValue": true,
            "label": "Auto Start"
        }
    ],
    "description": "Please fill the form before installing the plugin."
}