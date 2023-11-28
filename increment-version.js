const fs = require('fs');

const packageJsonPath = './package.json';
const packageJson = require(packageJsonPath);

function incrementVersion(version) {
    const parts = version.split('.');
    parts[2] = parseInt(parts[2], 10) + 1; // Increment patch version
    return parts.join('.');
}

packageJson.version = incrementVersion(packageJson.version);

fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2));
