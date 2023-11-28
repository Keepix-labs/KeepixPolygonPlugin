const { expect } = require('chai');
const { exec } = require('child_process');
const os = require('os');
const util = require('util');
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

describe('KeepixPolygonPlugin', function() {
    
    before(async function() {
        await enableExecutable();
    });

    it('should be able to install', async function() {
        const result = await execute({"key":"install"});
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to report installation status', async function() {
        const result = await execute({"key":"installed"});
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to report status', async function() {
        const result = await execute({"key":"status"});
        console.log(result)
        expect(result.jsonResult).to.equal(true);
    });

    it('should be able to uninstall', async function() {
        const result = await execute({"key":"uninstall"});
        expect(result.jsonResult).to.equal(true);
    });

    // Add more tests as needed
});
