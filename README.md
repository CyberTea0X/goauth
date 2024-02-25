# goauth

Simple authentification service written in golang that uses refresh and access tokens and communicates with external services.

[![Tests](https://github.com/CyberTea0X/goauth/actions/workflows/tests.yml/badge.svg)](https://github.com/CyberTea0X/goauth/actions/workflows/tests.yml)
[![Docker CI](https://github.com/CyberTea0X/goauth/actions/workflows/docker-image.yml/badge.svg)](https://github.com/CyberTea0X/goauth/actions/workflows/docker-image.yml)


## How to run

### Dockerhub

Copy config_test.toml from this repository and rename it to config.toml.
Edit config.toml and set your preferred values.

Then you can enter the following command:

```bash
docker run \
    -v ./config.toml:/goauth/config.toml \
    --network="host" \
    --rm cybertea0x/goauth
```

### Build from source

Make sure you installed and configured the latest golang version.

Clone this repository

```bash
git clone https://github.com/CyberTea0X/goauth
```

Then build binary

```bash
cd goauth
go build -o goauth
```

Place the binary where you need it and then copy the `config_test.toml` file there.
Rename `config_test.toml` to `config.toml`, configure this file, make sure mysql is running
and then you can start the server by running the executable.
