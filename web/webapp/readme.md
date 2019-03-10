# Micro Web Dashboard

Dashboard is the web front page of micro admin service.

# Important

Be different from the old version which used backend template to render html, new UI separate backend and frontend codes.
It makes more easy to write frontend code and import modern UI framework like vue/vuetify, etc.

So, when using our new UI, you should copy the dist [dir](./dist) to some path you want, and tell the **micro** command where it is. 

```bash

$ micro web --static-dir=/path/to/new/ui/dist/dir

```

or use the default dir **/usr/local/var/www/micro-web**:

```bash
$ cp -r dist /usr/local/var/www
$ mv /usr/local/var/www/dist /usr/local/var/www/micro-web
$ micro web

```

## Install

### npm modules

```bash
$ npm install
```

### initialize eslint(optional)

```bash
$ ./node_modules/.bin/eslint --init

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
$ npm run serve
```

## build

```bash
$ npm run build

```