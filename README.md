# cert-catcher

If you create machines on the fly to run your integration tests and your services run over TLS, you
may run into [rate limits](https://github.com/tailscale/tailscale/issues/10395#issuecomment-1934383393) 
when requesting the certificates. 

<p align="center">
  <img align="center" src="magic.webp" width="600px" alt="The cert cacher magician"/>
</p>

Cert-cacher is a simple [tsnet](https://tailscale.com/kb/1244/tsnet) cert caching service that leverages the
built-in identity of [tailnet](https://tailscale.com/glossary/tailnet).

### Usage

```sh
$ go install github.com/drio/cert-catcher@latest
# Run the service
$ cert-catcher 
```

Now you can make requests to...

```sh
# Get the private key/cert for the cert associated with the machine that makes the request
$ curl http://cert-cacher:9191/key
$ curl http://cert-cacher:9191/cert

# Check how many days before the cert expires:
$ curl http://cert-cacher:9191/cert

#  
```
