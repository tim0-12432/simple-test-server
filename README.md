# Simple servers for testing

This is a collection of simple servers for testing. They are written in Python and are intended to be run on a local machine or as a docker container.

## Quickstart with Docker

```bash
docker run --name http-test -t -d -p 8080:80 ghcr.io/tim0-12432/simple-test-server:latest
```

## Table of Contents

1. [Server types](#server-types)
2. [Running the servers](#running-the-servers)
   1. [Locally](#locally)
   2. [Docker](#docker)
3. [Environment variables](#environment-variables)
4. [Additional information](#additional-information-for-testing-if-the-server-works)
   1. [HTTP Server](#http-server)
   2. [FTP Server](#ftp-server)
   3. [SSH Server](#ssh-server)
   4. [SMTP Server](#smtp-server)
5. [Licensing](#licensing)

## Server types

- [HTTP server](./http_server.py)
- [FTP server](./ftp_server.py)
- [SSH server](./ssh_server.py)
- [SMTP server](./smtp_server.py)

## Running the servers

### Locally

To run the servers locally, you need to have Python 3 and all [requirements](./requirements.txt) installed. Then you can run the servers with the following commands:

```bash
python http_server.py
python ftp_server.py
python ssh_server.py
python smtp_server.py
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

**Example**: Run the HTTP server on port 8080:

```bash
docker run --name http-test -t -d -p 8080:80 --env TYPE=http ghcr.io/tim0-12432/simple-test-server:latest
```

**Example**: Run the FTP server on port 21 with custom files:

```bash
docker run --name ftp-test -t -d -p 21:21 --env TYPE=ftp -v ftp_files:/app/ftp-files ghcr.io/tim0-12432/simple-test-server:latest
```

- User: `anonymous`, no password -> only read access
- User: `admin`, `test` -> read and write access

## Environment variables

The following environment variables can be used to configure the servers:

|   Variable | Description | Default value |
| ---------: | ----------- | :------------ |
|     `TYPE` | The type of the server. Can be `http`, `ssh`, `smtp` or `ftp`. | `http` |
|  `ADDRESS` | The host:port on which the server should listen. | `0.0.0.0:<application>` |

## Additional information for testing if the server works

### HTTP server

By default, the HTTP server serves the files in the `http-files` directory. You can add your own files to this directory and they will be served by the server.
By default, the server is reachable at the port `80`. But you can change this by setting the `ADDRESS` environment variable.

Open the webpage by opening the address URL in your preferred browser!

### FTP server

By default, the FTP server serves the files in the `ftp-files` directory. You can add your own files to this directory and they will be served by the server.
By default, the server is reachable at the port `21`. But you can change this by setting the `ADDRESS` environment variable.

You can connect to the FTP server with your preferred FTP client.

```bash
ftp
open localhost
```

### SSH server

By default, the SSH server is reachable at the port `22`. But you can change this by setting the `ADDRESS` environment variable.

You can connect to the SSH server with your preferred SSH client.

```bash
ssh localhost
```

### SMTP server

By default, the SMTP server is reachable at the port `587`. But you can change this by setting the `ADDRESS` environment variable.

You can send emails to the SMTP server with your preferred SMTP client.

```powershell
Send-MailMessage -From 'User01 <user01@test.com>' -To 'User02 <user02@test.com>' -Subject 'Test mail' -Body 'This is a test message.' -SmtpServer '127.0.0.1' -Port 587
```

## Licensing

This project is licensed under the [MIT License](https://en.wikipedia.org/wiki/MIT_License) - see the [LICENSE](./LICENSE.md) file for details.
