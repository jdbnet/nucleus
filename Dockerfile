# ── Stage 1: Build Vue SPA ──
FROM node:22-alpine AS frontend
WORKDIR /build
COPY frontend/package*.json ./
RUN npm ci --no-audit --no-fund 2>/dev/null || npm install --no-audit --no-fund
COPY frontend/ ./
RUN npm run build

# ── Stage 2: Build Go server ──
FROM golang:1.26-alpine AS server
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /cmd/nucleus/dist ./cmd/nucleus/dist
ARG VERSION=dev
RUN CGO_ENABLED=1 go build \
  -ldflags="-s -w -X nucleus/internal/version.Version=${VERSION}" \
  -o nucleus ./cmd/nucleus

# ── Stage 3: Runtime ──
FROM alpine:3.24

ARG NUCLEI_VERSION=3.11.1
ARG TARGETARCH

RUN apk add --no-cache ca-certificates tzdata wget unzip \
  && wget -qO /tmp/nuclei.zip \
    "https://github.com/projectdiscovery/nuclei/releases/download/v${NUCLEI_VERSION}/nuclei_${NUCLEI_VERSION}_linux_${TARGETARCH}.zip" \
  && unzip -q /tmp/nuclei.zip -d /tmp/nuclei \
  && mv /tmp/nuclei/nuclei /usr/local/bin/nuclei \
  && chmod +x /usr/local/bin/nuclei \
  && rm -rf /tmp/nuclei.zip /tmp/nuclei \
  && nuclei -update-templates -silent

WORKDIR /app
COPY --from=server /build/nucleus ./nucleus

EXPOSE 8080
ENTRYPOINT ["./nucleus"]
