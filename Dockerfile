FROM m.daocloud.io/docker.io/library/node:21.6.2 AS builder

WORKDIR /src

RUN npm install -g pnpm

COPY . .

RUN pnpm config set registry https://registry.npmmirror.com 
RUN make build-all

FROM m.daocloud.io/docker.io/library/node:21.6.2-slim

COPY --from=builder /src/dist/ /apps/dist/
COPY --from=builder /src/node_modules/ /apps/node_modules/
COPY --from=builder /src/public/ /apps/public/
