# dillo-ipfs

Connect to [IPFS][] sites in [Dillo][].

# Install

```sh
git clone ssb://%C35b+MlZ/y5TT1e7SG66eNKEIdX5DRl9PRUxbhvO89k=.sha256 dillo-ipfs
cd dillo-ipfs
go build ./ipfs.dpi.go
mkdir -p ~/.dillo/dpi/ipfs
ln -rs ipfs.dpi ~/.dillo/dpi/ipfs
test -f ~/.dillo/dpidrc || cp /etc/dillo/dpidrc ~/.dillo/dpidrc
echo 'proto.ipfs=ipfs/ipfs.dpi' >> ~/.dillo/dpidrc
echo 'proto.ipns=ipfs/ipfs.dpi' >> ~/.dillo/dpidrc
dpidc stop
```

# Usage

Navigate to `ipfs://` and `ipns://` URLs in dillo like any other URLs.

Examples:
- <ipfs://QmYNQJoKGNHTpPxCBPh9KkDpaExgd2duMa3aF6ytMpHdao/>
- <ipns://dist.ipfs.io/>

You will need a local `ipfs daemon` running on port 8080.

[IPFS]: https://ipfs.io/
[Dillo]: https://dillo.org/

## License

AGPLv3+
