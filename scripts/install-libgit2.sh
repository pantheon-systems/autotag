set -e

version="0.23.2"
dirname="libgit2-$version"
basedir=$(pwd)


sudo apt-get install -y libssh2-1-dev openssl pkg-config
if [  -d ~/libgit2 ] ; then
  echo "libgit2 already setup"
  exit 0
fi


sudo apt-get update
sudo apt-get install -y curl build-essential python cmake

curl -L "https://github.com/libgit2/libgit2/archive/v${version}.tar.gz" -o ~/${dirname}.tar.gz
tar -xzvf ~/${dirname}.tar.gz -C ~/

cd ~/${dirname}
mkdir build && cd build
cmake .. -DCMAKE_INSTALL_PREFIX=~/libgit2
cmake --build . --target install
