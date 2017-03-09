#!/bin/bash

VERSION="0.4d"

if [ "$(basename $(pwd))" != "deploy" ]; then
	echo "must be in deploy dir" >&2
	exit 1
fi

deploy_win32_jdidco() {
  GOOS=windows GOARCH=386 go build -v -o lbcbot.exe jdid.co/lbcbot || exit 1
  test -x lbcbot.exe || exit 1

  tmp=$(mktemp -d)
  src=$tmp/lbcbot/
  mkdir -p $src

  target="lbcbot_win32-v${VERSION}.zip"
  cp -ar ../README.md ../config_sample.json ../html lbcbot.exe $src/
  echo "zipping $target"
  cd $(dirname $src)
  zip -r $target lbcbot/
  cd -
  mv $(dirname $src)/$target .
  rm -rf $tmp

  echo "sending $target to remote.."
  scp $target wservices:share/
}

deploy_win32_jdidco
