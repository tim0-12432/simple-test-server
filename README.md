# Simple servers for testing

This is a collection of simple servers for testing. They are written in Python and are intended to be run on a local machine or as a docker container.

## Server types

- [HTTP server](./http_server.py)
- [FTP server](./ftp_server.py)

## Running the servers

### Locally

To run the servers locally, you need to have Python 3 and all [requirements](./requirements.txt) installed. Then you can run the servers with the following commands:

```bash
python http_server.py
python ftp_server.py
```

or you can run it by configuring the environment variables via a `.env`-file and calling the `run.py` script directly:

```bash
python run.py
```

### Docker

To run the servers in a docker container, you need to have docker installed. Then you can run the servers with the following commands:

```bash
docker pull ghcr.io/tim0-12432/simple-test-server:latest
```

Example: Run the HTTP server on port 8080:

```bash
docker run --name http-test -d -p 8080:80 --env TYPE=http ghcr.io/tim0-12432/simple-test-server:latest
```

## Environment variables

The following environment variables can be used to configure the servers:

|   Variable | Description | Default value |
| ---------: | ----------- | :------------ |
|     `TYPE` | The type of the server. Can be `http` or `ftp`. | `http` |
|  `ADDRESS` | The host:port on which the server should listen. | `0.0.0.0:<application>` |

## Licensing

This project is licensed under the [MIT License](https://en.wikipedia.org/wiki/MIT_License) - see the [LICENSE](./LICENSE.md) file for details.
