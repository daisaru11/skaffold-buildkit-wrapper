# syntax = docker/dockerfile:experimental
FROM alpine

RUN apk add --no-cache git

ARG TEST
RUN echo $TEST

RUN --mount=type=secret,id=gitconfig,dst=/root/.gitconfig \
    git clone https://github.com/daisaru11/sample_private_repo.git
RUN echo sample_private_repo/sample.txt