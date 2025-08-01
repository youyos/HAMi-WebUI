FROM --platform=$BUILDPLATFORM m.daocloud.io/docker.io/library/golang:1.23.1 AS builder

WORKDIR /src
ARG TARGETARCH

# 设置 Go 模块代理为清华源，避免国外网络超时
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 替换 APT 源为清华，并清理默认源配置
RUN set -eux; \
    release=$(grep -oP '(?<=VERSION_CODENAME=).*' /etc/os-release); \
    echo "Using Debian release: $release"; \
    echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian $release main contrib non-free\n\
deb http://mirrors.tuna.tsinghua.edu.cn/debian $release-updates main contrib non-free\n\
deb http://mirrors.tuna.tsinghua.edu.cn/debian-security $release-security main contrib non-free" \
    > /etc/apt/sources.list && \
    rm -rf /etc/apt/sources.list.d/* && \
    apt-get update && \
    apt-get install -y --no-install-recommends protobuf-compiler && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY . .

RUN make build TARGET_ARCH=${TARGETARCH}

# 运行阶段镜像
FROM --platform=$TARGETPLATFORM m.daocloud.io/docker.io/library/debian:stable-slim

# 如果运行时不需要 APT，这一段可以删除（否则保留）
RUN set -eux; \
    release=$(grep -oP '(?<=VERSION_CODENAME=).*' /etc/os-release); \
    echo "Using Debian release: $release"; \
    echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian $release main contrib non-free\n\
deb http://mirrors.tuna.tsinghua.edu.cn/debian $release-updates main contrib non-free\n\
deb http://mirrors.tuna.tsinghua.edu.cn/debian-security $release-security main contrib non-free" \
    > /etc/apt/sources.list && \
    rm -rf /etc/apt/sources.list.d/* && \
    apt-get update && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/build/ /apps/
COPY --from=builder /src /src

