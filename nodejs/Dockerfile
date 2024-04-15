FROM grobuzin/main-agent:latest as main-agent

# This is gonna be executed from the root of the project (grobuzin-agent), with only the user-code and corresponding engine/agent (nodejs) folder in the context


FROM alpine:3.19 as builder

RUN apk add nodejs npm

WORKDIR /build

COPY ./nodejs .

RUN npm ci

RUN npm run build

# Build user code 
FROM alpine:3.19 as user-code

RUN apk add nodejs npm

WORKDIR /build

COPY ./user-code .

RUN ls

RUN npm install

# Optional: prepare the user code (build, etc.)
RUN npm run prepare || true

FROM alpine:3.19

WORKDIR /app

RUN apk add nodejs npm

# Copy the config file

COPY ./nodejs/agent-config.json /etc/grobuzin/agent-config.json

COPY --from=main-agent /build/main-agent .
COPY --from=builder /build/dist/* .
COPY --from=user-code /build/. ./function
