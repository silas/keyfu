# KeyFu

KeyFu is a simple search site that accepts keywords and executes some type of
action based on that keyword. The most common action is to map a keyword and
query string to a URL.

You can see examples of these mappings in the `keyfu.conf` configuration file.

One such example is the `gh` keyword, which maps the bare `gh` keyword to
`github.com` and the query keyword to `github.com/search?q=%s`, where `%s` is
replaced with the query parameter. This allows you to enter `gh` into your
search box and get redirected to the GitHub front page, or `gh keyfu` to search
GitHub for the KeyFu project. You can try this and all keywords in the
`keyfu.conf` configuration file on the [KeyFu][keyfu] demo site.

## Usage

 1. [Install Go][go-install]

    Ensure the `GOPATH` environment variable is set and `$GOPATH/bin` is in your `PATH`.

    ``` console
    $ echo 'export GOPATH="$HOME/go"' >> ~/.bashrc
    $ echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.bashrc
    $ source ~/.bashrc
    ```

 1. Get `github.com/silas/keyfu` and setup dependencies

    ``` console
    $ go get github.com/silas/keyfu
    $ cd $GOPATH/src/github.com/silas/keyfu
    $ make setup
    ```

 1. Build and install `keyfu`

    ``` console
    $ make install
    ```

 1. Run `keyfu`

    ``` console
    $ keyfu -c keyfu.conf
    ```

 1. Open site [localhost:8000](http://localhost:8000/)

    You can specify an alternative listen interface/port by setting the `listen` option in `keyfu.conf`.

    ``` toml
    listen = "127.0.0.1:8888"
    ```

    Or via the `HOST` and `PORT` environment variables.

 1. There are various ways start `keyfu` on boot, see the `contrib` directory for examples.

## Config

#### General

``` toml
listen = "127.0.0.1:8888"
```

#### Link

``` toml
[keyword.gh]
type = "link"
url = "https://github.com/"
query_url = "https://github.com/search?q=%s"
```

#### Alias

``` toml
[keyword.github]
type = "alias"
name = "gh"
```

### License

This work is licensed under the MIT License (see the LICENSE file).

[keyfu]: http://www.keyfu.com/
[go-install]: http://golang.org/doc/install
