# PocketBook Cloud Sync

[![Docker Pulls](https://img.shields.io/docker/pulls/micronull/pbcsync)](https://hub.docker.com/r/micronull/pbcsync)

Backup your library from [PocketBook Cloud](https://cloud.pocketbook.digital).

## Usage

### Docker

```shell
mkdri books
docker run \
--rm \
-v books:/books \
-e PBC_CLIENT_ID=qNAx1RDb \
-e PBC_CLIENT_SECRET=K3YYSjCgDJNoWKdGVOyO1mrROp3MMZqqRNXNXTmh \
-e PBC_USERNAME=your@email.some \
-e PBC_PASSWORD=your_super_mega_password_123 \
-e DEBUG=true \
micronull/pbcsync:latest
```

### Build and run

```shell
go build -o pbcsync ./cmd/pbcsync/

./pbcsync \
-client-id qNAx1RDb \
-client-secret K3YYSjCgDJNoWKdGVOyO1mrROp3MMZqqRNXNXTmh \
-username your@email.some \
-password your_super_mega_password_abc \
-debug \
-dir /some/dir
```

## Help sync

```txt
Usage of sync:
  -client-id string
        Client ID of PocketBook Cloud API.
        Read the readme to find out how to get it.
  -client-secret string
        Client Secret of PocketBook Cloud API.
        Read the readme to find out how to get it.
  -daemon
        Enable daemon mode. Use the daemon-timeout flag for setting sync interval.
  -daemon-timeout duration
        Timeout for sync operation. 
        Used only daemon mode. (default 24h0m0s)
  -debug
        Enable debug output.
  -dir string
        Directory to sync files. (default "books")
  -env
        Enable environment variables mode.
        Ignores all command-line flags and loads values from environment variables:
        PBC_CLIENT_ID as -client-id
        PBC_CLIENT_SECRET as -client-secret
        PBC_USERNAME as -username
        PBC_PASSWORD as -password
        DEBUG as -debug
        DIR as -dir
        DAEMON as -daemon
        DAEMON_TIMEOUT as -daemon-timeout
  -password string
        Password from your PocketBook Cloud account.
  -username string
        Username of PocketBook Cloud. Usually it's your email.
```
