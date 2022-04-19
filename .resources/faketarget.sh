#!/bin/sh

# export AWS_ACCESS_KEY_ID=access_key
# export AWS_ACCESS_KEY_ID=secret_key
# PORT=4569
# docker run --rm --name s3 -p $PORT:4569 -d stanislavt/s3-emulator

# -v <host-dir>/users.conf:/etc/sftp/users.conf:ro
docker run --rm -p 22:22 -d atmoz/sftp foo:pass:::upload

# minioadmin:minioadmin
docker run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"