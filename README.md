# goauth

## How to run

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
