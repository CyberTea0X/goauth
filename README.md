# goauth

## How to run

Make sure you installed:

    - golang, at least 1.21 version.
    - mysql (mysql server running)

Then go to the backend directory and create .env file like this:

```bash
DB_HOST=127.0.0.1                       
DB_DRIVER=mysql                          
DB_USER=someuser
DB_PASSWORD=somepw
DB_NAME=goauth
DB_PORT=3306 
ACCESS_TOKEN_SECRET=secret1
REFRESH_TOKEN_SECRET=secret2
ACCESS_TOKEN_MINUTE_LIFESPAN=15
REFRESH_TOKEN_HOUR_LIFESPAN=8760
```

Then type

```bash
go run .
```
