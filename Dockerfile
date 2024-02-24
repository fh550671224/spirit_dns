FROM ubuntu:latest
LABEL authors="Richard"

ENTRYPOINT ["top", "-b"]