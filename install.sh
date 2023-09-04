#!/bin/bash
binary_name=gitd
goos=$(uname)
version=0.0.4
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	goos=linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
	goos=macos
else
	echo "Error: The current os is not supported at this time" 1>&2
	exit 1
fi

# check the gitd if exist
ret=0
command -v gitd >/dev/null 2>&1 || { local ret=1; }

# fail on non-zero return value
if [ "$ret" -ne 0 ]; then
	echo "delete the old gitd."
	rm -rf $(which gitd)
fi

file_name=gitd-v${version}-${goos}

url=https://github.com/codexfield/gitd/releases/download/v${version}/${file_name}
echo "Download url:${url}"

curl "$url" -OL --retry 2 2>&1

mv ${file_name} ${binary_name}
chmod u+rwx $binary_name
sudo mv $binary_name /usr/local/bin/

echo "gitd install success."
