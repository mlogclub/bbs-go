#!/bin/bash
cd ./
path=`pwd`
echo "build path:"$path
files=$(ls $path)
for filename in $files
do
	if [[ ${filename:0:6} == "cscope"  ]];then
		rm cscope*
		break
	fi
done

if [ -f tags  ];then
	rm tags
fi
gotags -R ./ > tags
find ./ -name "*.go" > ./cscope.files
cscope -Rbq
echo "build tags succeed!"
