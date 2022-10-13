# Requirements
FROM python:3.7.15-slim
COPY /requirements.txt /app/requirements.txt
WORKDIR /app
RUN pip install -r requirements.txt

# Run App
COPY /ftp-files /app/ftp-files
COPY /http-files /app/http-files
COPY /ftp_server.py /app/ftp_server.py
COPY /http_server.py /app/http_server.py
COPY /run.py /app/run.py
ENTRYPOINT [ "python" ]
CMD ["run.py"]

# Display metadata
LABEL org.opencontainers.image.title "Simple servers for testing"
LABEL org.opencontainers.image.description "This is a simple container providing servers for testing written in python."
LABEL org.opencontainers.image.url "https://github.com/tim0-12432/simple-test-server"
LABEL org.opencontainers.image.source "https://github.com/tim0-12432/simple-test-server"
LABEL org.opencontainers.image.documentation "https://github.com/tim0-12432/simple-test-server"
LABEL org.opencontainers.image.authors "tim0-12432"
LABEL org.opencontainers.image.licenses "MIT"
