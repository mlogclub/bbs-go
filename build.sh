#!/bin/bash

### 构建要求
### 1. go >= 1.13
### 2. node >= 8

baseDir=`echo $PWD`
serverDir=`echo $baseDir/server`
siteDir=`echo $baseDir/site`
distDir=`echo $baseDir/dist`

# go version
# go env
# echo $baseDir
# echo $serverDir
# echo $siteDir


funcBuildServer() {
    echo 'server module building...'
    export GOPROXY=https://goproxy.cn
    cd $serverDir
    go mod download
    go build
    echo 'server module building...finished'
}


funcBuildSite() {
    echo 'site module building...'
    cd $siteDir
    npm install --registry=https://registry.npm.taobao.org
    npm run build
    echo 'site module building...finished'
}

funcTouchDir() {
    if [ ! -d "$1" ]; then
        mkdir $1
    fi
}

funcCleanBuild() {
    rm -rf $distDir/*

    funcTouchDir $distDir/server
    funcTouchDir $distDir/site

    cp $serverDir/bbs-go $distDir/server/
    
    cp -r $siteDir/.nuxt $distDir/site/
    cp -r $siteDir/static $distDir/site/
    cp -r $siteDir/nuxt.config.js $distDir/site/
    cp -r $siteDir/package.json $distDir/site/
}

funcBuildServer
funcBuildSite
funcCleanBuild

echo 'done!!!'
