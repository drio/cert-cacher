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

```
#$ go install github.com/drio/cert-catcher@latest
$ git clone github.com/drio/cert-catcher
$ cd cert-catcher
$ go build .
# This will store the certs in memory (use -disk to store in disk)
$ ./cert-catcher
```

Now you can make requests from any node in the tailnet. Remember to update your [ACL](https://tailscale.com/kb/1018/acls))
to give proper permissions to the machines you want to hace access to the cert-catcher.

```
# Save your certs (the service will look at the file header to determine what you are sending):
$ curl  --data-binary @./machine.tailnet-name.ts.net.cert http://cert-cacher:9191/
$ curl  --data-binary @./machine.tailnet-name.ts.net.key http://cert-cacher:9191/

# Get the private key/cert for the cert associated with the machine that makes the request
$ curl http://cert-cacher:9191/key
$ curl http://cert-cacher:9191/cert

# Check how many days before the cert expires:
$ curl http://cert-cacher:9191/cert
```

All these requires you to issue the certs. 
You will probably use [`tailscale cert`](https://tailscale.com/kb/1153/enabling-https) for that.
The service comes with a shell script that can help you with that:

```
$ curl -s http://cert-cacher:9191/sh | sh -s -- -d m3.tailnet.net
```
