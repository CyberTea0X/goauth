# goauth

## How to run

Make sure you installed mysql and mysql server is running.
Configure user, database and password for this service.

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

Place the binary where you need it and then copy the `example.toml` file there.
Rename `example.toml` to `config.toml`, configure this file, make sure mysql is running
and you can start the server by running executable.
