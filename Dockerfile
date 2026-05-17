FROM ubuntu:26.04

ENV DEBIAN_FRONTEND=noninteractive 
RUN apt update && apt install -y \
    python3-pip \
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
    libosmesa6 \
    libgl1 \
    libglx-mesa0 \
    libglu1-mesa-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install --break-system-packages \
    trimesh \
    Pillow \
    numpy \
    pyopengl \
    pyopengl-accelerate && \
    pip3 install --break-system-packages pyrender==0.1.18 && \
    pip3 install --break-system-packages --upgrade networkx

WORKDIR /opt 
RUN git clone https://github.com/mchehab/zbar.git 
COPY ean.c /opt/zbar/zbar/decoder/ean.c
COPY img_scanner.c /opt/zbar/zbar/img_scanner.c
WORKDIR /opt/zbar 
RUN autoreconf -vfi 
RUN ./configure --with-gtk=no --with-python=no && \
    make -j$(nproc) && \
    make install && \
    ldconfig 

WORKDIR /opt
RUN git clone https://github.com/zint/zint.git

WORKDIR /opt 
RUN git clone https://github.com/zxing-cpp/zxing-cpp.git
WORKDIR /opt/zxing-cpp 
RUN cmake -S . -B build \
 -DCMAKE_BUILD_TYPE=Release && \
 cmake --build build -j$(nproc) 

WORKDIR /work
CMD ["/bin/bash"]