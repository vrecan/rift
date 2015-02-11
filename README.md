# rift
[![Build Status](https://travis-ci.org/vrecan/rift.svg?branch=master)](https://travis-ci.org/vrecan/rift)<br>
Zeromq process that takes a zeromq pull socket and relays the data to any number of end points.


# Windows install help
The windows install script was derived from the following: <br>
[Windows Install guide](http://www.limitlessfx.com/cross-compile-golang-app-for-windows-from-linux.html) <br>
On centos 7 execute this script
```
sh windows_install.sh
```
<p>This will install mingw & zeromq from source and then compile rift with both. Then copy the dll's required to the current working directory.</p>
