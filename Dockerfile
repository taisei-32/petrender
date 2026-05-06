FROM ubuntu:26.04

ENV DEBIAN_FRONTEND=noninteractive 
RUN apt update && apt install -y \
    autopoint \
    gettext \
    autoconf \
    automake \
    libtool \
    pkg-config \
    libjpeg-dev \
    libmagickwand-dev \
    imagemagick \
    python3 \
    python3-dev \
    zint \
    git \
    make \
    gcc \
    g++ \
    cmake \
    build-essential \
    golang-go \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /opt 
RUN git clone https://github.com/mchehab/zbar.git 
WORKDIR /opt/zbar 
RUN autoreconf -vfi 
RUN ./configure --with-gtk=no --with-python=no && \
    make -j$(nproc) && \
    make install && \
    ldconfig 
WORKDIR /opt 
RUN git clone https://github.com/zxing-cpp/zxing-cpp.git
WORKDIR /opt/zxing-cpp 
RUN cmake -S . -B build \
 -DCMAKE_BUILD_TYPE=Release && \
 cmake --build build -j$(nproc) 
CMD ["/bin/bash"]