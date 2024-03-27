# Requirements
FROM python:3.11-slim AS builder
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc
COPY /requirements.txt /app/requirements.txt
WORKDIR /app
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
ARG DEBIAN_FRONTEND=noninteractive
ARG DEBCONF_NOWARNINGS="yes"
RUN python -m pip install --upgrade pip && \
    pip install -r requirements.txt

# Run App
FROM python:3.11-slim AS runner
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY /ftp-files /app/ftp-files
COPY /http-files /app/http-files
COPY /smb-share /app/smb-share
COPY /abstract_server.py /app/abstract_server.py
COPY /ftp_server.py /app/ftp_server.py
COPY /http_server.py /app/http_server.py
COPY /smtp_server.py /app/smtp_server.py
COPY /smb_server.py /app/smb_server.py
COPY /ssh_server.py /app/ssh_server.py
COPY /run.py /app/run.py
WORKDIR /app
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
