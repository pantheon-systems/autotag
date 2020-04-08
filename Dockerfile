FROM debian:buster-slim

RUN apt-get update \
  && apt-get install --yes --no-install-recommends \
    git \
  && apt-get autoremove --purge \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY ./autotag /autotag

ENTRYPOINT [ "/autotag" ]
