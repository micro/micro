# Micro Web Dashboard

Micro web is the dashboard for visualising and exploring services.


## Overview

Run the micro web app like so

```bash
micro web
```

This assumes you are in the root dir of github.com/micro/micro. In the event you want to run from a different directory specify the 
static dir path.

```
micro web --static_dir=/path/to/new/ui/dist/dir
```

The default location is ./web/webapp/dist

## Development

Micro web makes use of vuejs

### npm modules

```bash
npm install
```

### initialize eslint(optional)

```bash
./node_modules/.bin/eslint --init

? How would you like to use ESLint? To check syntax and find problems
? What type of modules does your project use? JavaScript modules (import/export)
? Which framework does your project use? Vue.js
? Where does your code run? (Press <space> to select, <a> to toggle all, <i> to invert selection)Browser
? What format do you want your config file to be in? JSON
Checking peerDependencies of eslint-config-eslint:recommended,plugin:vue/essential@latest
The config that you've selected requires the following dependencies:

eslint-plugin-vue@latest eslint-config-eslint:recommended,plugin:vue/essential@latest
? Would you like to install them now with npm? Yes

```

## run

```bash
npm run serve
```

## build

```bash
npm run build
```

# Thanks

## [elementUI](https://element.eleme.io/#/)
## [Baidu eCharts](https://github.com/apache/incubator-echarts)
## [vue-material-admin](https://github.com/tookit/vue-material-admin)
