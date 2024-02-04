# IPFS plugin for Dillo

Connect to [IPFS][] sites in [Dillo][].

Created by [Charles E. Lehner](https://celehner.com/) along with other
[Dillo plugins](https://celehner.com/projects.html#dillo-plugins).

## Installation

To install the plugin use:

```
$ make install
```

## Usage

Navigate to `ipfs://` and `ipns://` URLs in Dillo like any other URLs.

Examples:
- <ipfs://QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/architecture/index.html>
- <ipfs://QmYNQJoKGNHTpPxCBPh9KkDpaExgd2duMa3aF6ytMpHdao/>
- <ipns://dist.ipfs.io/>

You will need a local `ipfs daemon` running on port 8080.

[IPFS]: https://ipfs.io/
[Dillo]: https://dillo-browser.github.io/

## License

AGPLv3+
