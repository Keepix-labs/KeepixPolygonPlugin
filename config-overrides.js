const path = require('path');
const http = require('http');
const express = require('express');
const fs = require('fs');
const { exec } = require('child_process');
const cors = require('cors');

const packageJsonContent = fs.readFileSync('package.json', 'utf8');
const packageJson = JSON.parse(packageJsonContent);
const pluginId = packageJson.name;

console.log(pluginId)

const execPluginFunction = async (argObject) => {
    return await new Promise((resolve) => {
        // double stringify for escapes double quotes
        exec(`./build/dist/${pluginId} ${JSON.stringify(JSON.stringify(argObject))}`, (error, stdout, stderr) => {
            const result = JSON.parse(stdout);
            console.log("OUT", result)

            resolve({
                result: JSON.parse(result.jsonResult),
                stdOut: result.stdOut
            });
        });
    });
};

const executeServer = () => {
    if (!process.argv.join(' ').includes('start.js')) {
        return ;
    }
    console.log('Start Dev Server.');

    const app = express();

    app.use(cors());
    app.use(express.json());
    app.set('port', process.env.PORT || 2000);
    app.set('host', '0.0.0.0');
    var runningTasks = {};

    const executePlugin = async (dto, isAsync = false) => {
        console.log("executing")
        if (isAsync) {
            const taskId = `${pluginId}-${dto.key}`;

            // skip duplicate running
            if (runningTasks[taskId] != undefined && runningTasks[taskId].status != 'FINISHED') {
                return {
                    taskId: taskId,
                    aborted: true,
                    reason: "Already running"
                };
            }

            // save the task status
            runningTasks[taskId] = {
                status: 'RUNNING'
            };

            // run asynchronously
            execPluginFunction(dto).then((result) => {
                runningTasks[taskId] = {
                    status: 'FINISHED',
                    description: result
                };
            });

            return { taskId: taskId };
        } else {
            return await execPluginFunction(dto);
        }
    }

    app.get(`/plugins/${pluginId}/:key`, async (req, res) => {
        const isAsync = req.query.isAsync === 'true';
        res.send(await executePlugin({
            key: req.params.key
        }, isAsync));
    });
    app.post(`/plugins/${pluginId}/:key`, async (req, res) => {
        const isAsync = req.query.isAsync === 'true';

        console.log(req.body);

        res.send(await executePlugin({
            key: req.params.key,
            ... req.body
        }, isAsync));
    });
    app.get(`/plugins/${pluginId}/watch/task/:taskId`, (req, res) => {
        res.send(runningTasks[req.params.taskId]);
    });

    // plugin view static files
    app.use((req, res, next) => {
        if (req?.url?.startsWith('/plugins/')) {
        const pluginIdAndSubPath = (req.url.split('/plugins/')[1]).split('/');
        // const pluginId = pluginIdAndSubPath[0];
        const subPathWithView = pluginIdAndSubPath.slice(1).join('/');

        if (subPathWithView.startsWith('view')) {
            const subPath = subPathWithView.substr(4);

            const staticPluginPath = path.join(__dirname, 'public');
            
            if (!fs.existsSync(path.join(staticPluginPath, subPath.split('?')[0]))) {
            req.url = '/index.html';
            }
            req.url = subPath;

            express.static(staticPluginPath)(req, res, next);
            return ;
        }
        }
        next();
    });

    http.createServer(app).listen(app.get('port'), app.get('host'), function(){
        console.log('Express server listening on port ' + app.get('port'));
    });
};

module.exports = {
    paths: function (paths, env) {
        executeServer();
        paths.appIndexJs = path.resolve(__dirname, 'src-front/index.tsx');
        paths.appSrc = path.resolve(__dirname, 'src-front');
        return paths;
    }
}