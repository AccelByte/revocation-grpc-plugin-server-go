# Platform Service's Revocation Plugin gRPC Demo App (Go)

## Prerequisites

- Docker
- make

## Usage

### Setup

The following environment variables are used by this CLI demo app.

Put environment variables in .env file:

```shell
AB_BASE_URL='https://test.accelbyte.io'
AB_CLIENT_ID='<AccelByte IAM Client ID>'
AB_CLIENT_SECRET='<AccelByte IAM Client Secret>'

AB_NAMESPACE='namespace'
AB_USERNAME='<AccelByte account username>'
AB_PASSWORD='<AccelByte account password>'

GRPC_SERVER_URL='<gRPC server url accessible by internet>'
```

### Run 

Run the demo test cli using makefile with provided .env file containing required variables

```shell
make run ENV_FILE_PATH=.env
```