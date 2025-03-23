#!/bin/bash

rm -rf public/images/album


cp album.zip public/images/

pushd public/images/

unzip album.zip 

popd