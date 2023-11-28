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