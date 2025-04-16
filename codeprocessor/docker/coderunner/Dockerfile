FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
    gcc \
    clang \
    python3 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* 

RUN useradd -u 1001 executor

WORKDIR /app

RUN chown -R executor:executor /app

USER executor

CMD ["/bin/bash", "-c", "while true; do sleep 1; done"]