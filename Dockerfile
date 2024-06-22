FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o go-schedule-manager .

RUN echo 'deb https://deb.debian.org/debian stable non-free contrib' >> /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y gcc \
    build-essential \
    ffmpeg \
    python3-pip \
    autoconf \
    automake \
    libtool \
    git \
    libasound2-dev \
    espeak \
    espeak-ng \
    yt-dlp \
    make \
    mbrola \
    mbrola-br4

ENTRYPOINT ["./go-schedule-manager"]

EXPOSE 8000