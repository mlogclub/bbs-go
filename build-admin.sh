#!/bin/sh

rm -rf ./build/mlog-admin.zip

cd ./web/admin/
npm run build

zip -r ../../build/mlog-admin.zip ./dist