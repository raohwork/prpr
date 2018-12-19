Simple application fetching static HTML via Firefox Marionette.

It is designed most for internal using. YOU HAVE TO PREVENT SECURITY ISSUE YOURSELF.

# Configuration

Prpr accepts following environmental variables:

- `SECRET`: shared secret. Default to `""` which disables authenticating.
- `FIREFOX`: path to Firefox binary. Default to `"firefox"`.
- `BIND`: binding address passed to `ListenAndServe`. Default to `":9801"`.

# Usage

```sh
wget \
  --post-data='uri=http://google.com&wait=queryString&secret=mysecret' \
  http://127.0.0.1:9801/grab
```

- `uri`: URL you want to fetch from.
- `wait`: CSS-selector to wait for.
- `secret`: your shared secret.

It returns HTML code if success, 400 if shared secret mismatch and 500 when failed.

# WARNING

* To keep code simple, prpr just waits 10 seconds before actually connect to firefox. IT DOES NOT SUPPORT RECONNECTING.
* For best security, DO NOT EXPOSE IT TO PUBLIC NETWORK.

# License

WTFPL
