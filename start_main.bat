@echo off

REM Starts the docker container for application.
docker run -d ^
  --name simple-test-server ^
  -p 8080:8080 ^
  -v /var/run/docker.sock:/var/run/docker.sock ^
  --restart unless-stopped ^
  simple-test-server:latest
