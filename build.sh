#!/usr/bin/env bash

dir=$(cd "$(dirname "$0")";pwd)
projectName=mlog
buildDir=$dir/build
outDir=$buildDir/$projectName

if [[ ! -d build  ]];then
  mkdir build
fi


# 清理历史
rm -rf $buildDir/*
mkdir $outDir

# 构建
echo 'building...'
GOOS=linux GOARCH=386 go build

## 构建admin
#cd $dir/web/admin
#npm run build

# 构建结果转移到build目录
mv $dir/$projectName $buildDir/$projectName
mkdir $outDir/web
cp -r $dir/web/views $outDir/web/views
cp -r $dir/web/static $outDir/web/static
#cp -r $dir/web/admin/dist $outDir/web/admin

echo "nohup ./$projectName -config /data/$projectName.yaml > ./output.log 2>&1 &" >> $outDir/start.sh
chmod +x $outDir/$projectName
chmod +x $outDir/start.sh

# 打包
cd $buildDir
zip -r $projectName.zip ./$projectName
