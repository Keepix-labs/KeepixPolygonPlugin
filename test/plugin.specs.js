const { expect } = require('chai');
const { exec } = require('child_process');
const os = require('os');
const util = require('util');
const path = require('path');
const fs = require('fs');
const execAsync = util.promisify(exec);

function getExecutablePath() {
    const arch = os.arch()
    switch (os.platform()) {
        case 'win32':
            return `./build/dist/win-${arch}/keepix-polygon-plugin.exe`;
        case 'darwin':
            return `./build/dist/osx-${arch}/keepix-polygon-plugin`;
        case 'linux':
            return `./build/dist/linux-${arch}/keepix-polygon-plugin`;
        default:
            return `./build/dist/keepix-polygon-plugin`;
    }
}

async function enableExecutable() {
    if(os.platform() === 'win32') return;
    const executable = getExecutablePath();
    await execAsync("chmod +x " + executable);
}

async function execute(jsonInput) {
    const escapedJsonInput = JSON.stringify(jsonInput).replace(/"/g, '\\"');
    const executablePath = getExecutablePath();
    const command = `${executablePath} "${escapedJsonInput}"`;
    const {stdout} = await execAsync(command);
    return JSON.parse(stdout);
}

function checkLocalPackageVersion() {
    try {
        // Construct the path to the package.json file
        const packageJsonPath = path.join(__dirname, '../package.json');

        // Check if package.json exists
        if (!fs.existsSync(packageJsonPath)) {
            console.error('package.json not found in the current directory');
            return;
        }

        // Load the package.json file
        const packageJson = require(packageJsonPath);

        // Return the version number
        return packageJson.version;
    } catch (error) {
        console.error('An error occurred:', error);
    }
}

const delay = ms => new Promise(resolve => setTimeout(resolve, ms));

describe('KeepixPolygonPlugin', function() {
    
    before(async function() {
        await enableExecutable();
    });
    
    it.skip('should be able to report version', async function() {
        const executablePath = getExecutablePath();
        const {stdout} = await execAsync(`${executablePath} --version`);
        const packageVersion = checkLocalPackageVersion();
        expect(stdout).to.equal('"' + packageVersion + '"');
    });

    it('should be able to install', async function() {
        const result = await execute({"key":"install","ethereumRPC":"https://eth-goerli.g.alchemy.com/v2/94XF2HyO7HcROFZuuBJ7EBxn1c68LdQm","testnet":"true","autostart":"false","mnemonic":"test test test test test test test test test test test junk"});
        console.log(result)
        expect(result.jsonResult).to.equal("true");
    });

    it('should be able to report installation status', async function() {
        const result = await execute({"key":"installed"});
        expect(result.jsonResult).to.equal("true");
    });

    it('should be able to import wallet from mnemonic', async function() {
        let result = await execute({"key":"wallet-load","mnemonic":"test test test test test test test test test test test junk","privateKey":""});
        console.log(result)
        expect(result.jsonResult).to.equal("true");
    });

    it('should be able to import wallet from private key', async function() {
        const result = await execute({"key":"wallet-load","mnemonic":"","privateKey":"rAl0vsOaF+NrpKa00jj/lEustHjL7V78rnhNe/Ty/4A="});
        expect(result.jsonResult).to.equal("true");
    });

    it('should be able to fetch wallet', async function() {
        const result = await execute({"key":"wallet-fetch"});
        const parsedResult = JSON.parse(result.jsonResult);
        expect(parsedResult.Wallet).to.equal('0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266');
    });

    it('should be able to start the nodes', async function() {
        const result = await execute({"key":"start"});
        expect(result.jsonResult).to.equal("true");
        await delay(30000);
    });

    it('should be able to restart the nodes', async function() {
        const result = await execute({"key":"restart"});
        console.log(result)
        expect(result.jsonResult).to.equal("true");
        await delay(30000);
    });

    it('should be able to report status', async function() {
        const result = await execute({"key":"status"});
        const parsedResult = JSON.parse(result.jsonResult);
        expect(parsedResult.NodeState).to.equal(`NodeStarted`);
    });

    it.skip('should be able to report sync state', async function() { // this will not work during tests since the introduction of snapshot downloader
        const result = await execute({"key":"sync-state"});
        const jsonResult = JSON.parse(result.jsonResult);
        console.log(result)
        expect(jsonResult.IsSynced).to.equal(false);
    });

    it('should be able to report chain', async function() {
        await delay(10000);
        const result = await execute({"key":"chain"});
        console.log(result)
        expect(result.jsonResult == "mainnet" || result.jsonResult == "testnet").to.true
    });

    it('should be able to report logs', async function() {
        const result = await execute({"key":"logs","erigon":"true","heimdall":"true", "lines": "1"});
        const jsonResult = JSON.parse(result.jsonResult);
        console.log(jsonResult)
        expect(jsonResult.erigonLogs).to.not.be.equal("");
        expect(jsonResult.heimdallLogs).to.not.be.equal("");
    });

    it('should be able to resync', async function() {
        const result = await execute({"key":"resync","erigon":"true","heimdall":"true"});
        console.log(result)
        expect(result.jsonResult).to.equal("true");
        await delay(10000);
    });

    it('should be able to stop the nodes', async function() {
        const result = await execute({"key":"stop"});
        expect(result.jsonResult).to.equal("true");
        await delay(10000);
    });

    it('should be able to uninstall', async function() {
        const result = await execute({"key":"uninstall"});
        console.log(result)
        expect(result.jsonResult).to.equal("true");
    });

    // Add more tests as needed
});
