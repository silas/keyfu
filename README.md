# KeyFu [![Build Status](https://drone.io/github.com/silas/keyfu/status.png)](https://drone.io/github.com/silas/keyfu/latest)

[KeyFu][keyfu] is a simple search site that accepts keywords and executes some type of
action based on that keyword.

## Usage

``` console
$ brew tap silas/silas
$ brew install keyfu
$ mkdir -p ~/.keyfu
$ echo 'link("https://github.com/", "https://github.com/search?q=%s");' > ~/.keyfu/github.js
$ ./keyfu &
$ curl -I 'localhost:8000/run?q=github'
HTTP/1.1 302 Found
Location: http://www.example.org/
$ curl -I 'localhost:8000/run?q=github+keyfu'
HTTP/1.1 302 Found
Location: https://github.com/search?q=keyfu
```

### License

This work is licensed under the MIT License (see the LICENSE file).

[keyfu]: http://www.keyfu.com/
