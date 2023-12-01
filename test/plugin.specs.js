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

describe('KeepixPolygonPlugin', function() {
    
    before(async function() {
        await enableExecutable();
    });

    it('should be able to report version', async function() {
        const executablePath = getExecutablePath();
        const {stdout} = await execAsync(`${executablePath} --version`);
        const packageVersion = checkLocalPackageVersion();
        expect(stdout).to.equal('"' + packageVersion + '"');
    });

    it('should be able to install', async function() {
        const result = await execute({"key":"install","ethereumRPC":"https://eth-mainnet.g.alchemy.com/v2/dWXI2QkWnTMsr7XAhlNzcD44m1qqemMS"});
        console.log(result)
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to report installation status', async function() {
        const result = await execute({"key":"installed"});
        console.log(result)
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to report status', async function() {
        const result = await execute({"key":"status"});
        console.log(result)
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to uninstall', async function() {
        const result = await execute({"key":"uninstall"});
        console.log(result)
        expect(result.jsonResult).to.equal(true);
    });

    // Add more tests as needed
});
