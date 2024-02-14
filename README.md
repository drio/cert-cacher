# cert-cacher

<p align="center">
  <img align="center" src="magic.webp" width="400px" alt="The cert cacher magician"/>
</p>

If you create machines on the fly to run your integration tests and your services run over TLS, you
may run into [rate limits](https://github.com/tailscale/tailscale/issues/10395#issuecomment-1934383393) 
when requesting certificates. 

Cert-cacher is a simple [tsnet](https://tailscale.com/kb/1244/tsnet) cert caching service that leverages the built-in identity on [tailnets](https://tailscale.com/glossary/tailnet).

### Usage

```
$ git clone github.com/drio/cert-cacher
$ cd cert-cacher
$ go build .
# This will store the certs in memory (use -disk to store in disk)
$ ./cert-cacher
```

Now you can make requests from any node in the tailnet. Remember to update your ([ACL](https://tailscale.com/kb/1018/acls))
to give proper permissions to the machines you want to have access to the cert-cacher.

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

All these require you to issue the certs. 
You will probably use [`tailscale cert`](https://tailscale.com/kb/1153/enabling-https) for that.
To help with that, cert-cacher comes with a shell script:

```
# Execute the script but only print the cmds you'd run (-p)
$ curl -s http://cert-cacher:9191/sh |  sh -s -- -p -d m3.XXX.ts.net
LOG> -p enabled, printing cmds only
LOG> /days status=404 days_to_expire=404 page not found
LOG> Cert not available in cacher. Requesting one and sending it to the cacher
tailscale cert m3.XXX.ts.net
curl --data-binary @./m3.XXX.ts.net.cert http://cert-cacher:9191
curl --data-binary @./m3.XXX.ts.net.key http://cert-cacher:9191
```

If the cert was cached:

```
$ curl -s http://cert-cacher:9191/sh |  sh -s -- -p -d m3.XXX.ts.net
LOG> -p enabled, printing cmds only
LOG> /days status=200 days_to_expire=359
LOG> Cert cached and valid. Getting it from the cacher
curl http://cert-cacher:9191/cert >m3.XXX.ts.net.cert
curl http://cert-cacher:9191/key >m3.XXX.ts.net.key
```
