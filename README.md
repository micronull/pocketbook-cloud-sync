# PocketBook Cloud Sync

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