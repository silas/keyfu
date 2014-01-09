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

[keyfu]: http://www.keyfu.com/
