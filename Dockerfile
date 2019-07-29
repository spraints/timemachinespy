FROM debian:stable-slim AS base
RUN apt-get update && apt-get install -y \
  cmake \
  curl \
  fuse libfuse-dev \
  libboost-system-dev libboost-filesystem-dev \
  build-essential pkg-config \
  hfsprogs parted

FROM base AS sparsebundlefs
RUN mkdir sparsebundlefs && \
  curl -L https://github.com/torarnv/sparsebundlefs/tarball/3cfb00979ad80a81d3901e9b77b3c1671a81e9d5 \
  | tar xvz -C /sparsebundlefs --strip-components 1
RUN cd sparsebundlefs && make

FROM base AS tmfs
RUN mkdir tmfs && \
  curl -L https://github.com/abique/tmfs/tarball/a5e201fb7693df253325896e61e7c26e6cc48114 \
  | tar xvz -C /tmfs --strip-components 1
RUN cd tmfs && \
  mkdir build && \
  cd build && \
  cmake -DCMAKE_INSTALL_PREFIX=/usr .. && \
  make
  
FROM base
COPY --from=sparsebundlefs /sparsebundlefs /sparsebundlefs
COPY --from=tmfs /tmfs /tmfs
RUN cp /sparsebundlefs/sparsebundlefs /usr/bin/
RUN cd tmfs/build && make install

RUN mkdir -p /mnt/sparse
RUN mkdir -p /mnt/hfs
RUN mkdir -p /mnt/tmfs
