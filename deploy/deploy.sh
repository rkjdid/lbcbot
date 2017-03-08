#!/bin/bash

VERSION="0.3"

deploy_win32_jdidco() {
  GOOS=windows GOARCH=386 go build -v -o lbcbot.exe jdid.co/lbcbot || exit 1
  test -x lbcbot.exe || exit 1

	tmp=$(mktemp -d)
  src=$tmp/lbcbot/
  mkdir -p $src


	target="lbcbot_win32-v${VERSION}.zip"
  cp -a README.md config_sample.json lbcbot.exe $src/
  echo "zipping $target"
  zip -r $target $src
  rm -rf $tmp

  echo "sending $target to remote.."
	scp $target wservices:share/
}

deploy_win32_jdidco
