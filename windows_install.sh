sudo yum install epel-release
sudo yum install mingw64*
rm -rf zeromq-4.0.5.zip
wget http://download.zeromq.org/zeromq-4.0.5.zip
unzip zeromq-4.0.5.zip ~/zeromq
cwd=`pwd`
cd ~/zeromq
mingw64-configure configure
mingw64-make
cd $cwd
cd /usr/local/go/src
sudo  env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC_FOR_TARGET="x86_64-w64-mingw32-gcc" ./make.bash
cd $cwd

env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC_FOR_TARGET="x86_64-w64-mingw32-gcc" CGO_LDFLAGS="-L/home/vrecan/zeromq/src/.libs -static-libgcc" go build -v -x  -o rift.exe
minigw_bin=/usr/x86_64-w64-mingw32/sys-root/mingw/bin/
cp $minigw_bin/libwinpthread-1.dll .
cp $minigw_bin/libgcc_s_seh-1.dll .
cp $minigw_bin/libstdc++-6.dll .
cp ~/zeromq/src/.libs/libzmq.dll .
