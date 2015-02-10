env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC_FOR_TARGET="i686-w64-mingw32-gcc" CGO_LDFLAGS="-L/home/vrecan/zeromq/src/.libs" go build -v -x  -o rift.exe
