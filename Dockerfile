FROM registry.gitlab.com/ninjavan/devops/golang:go-alpine as builder

ARG SERVICE_NAME
ENV PROJECT_DIR=/go/src/eb2b
WORKDIR /go/src/eb2b/backend


# RUN go get -u github.com/go-bindata/go-bindata/... \
RUN apk add go-bindata \
    && mkdir -p /build && mkdir -p /build/bin
COPY . .

## Copy assets & templates
COPY com/web/ecom/assets /build/com/web/ecom/assets
COPY com/web/ecom/templates /build/com/web/ecom/templates
COPY com/report/templates /build/com/report/templates

RUN sh ./scripts/generate-release.sh \
    && CGO_ENABLED=1 GOOS=linux go build \
      -o /build/bin/${SERVICE_NAME} \
      -tags release \
      ./cmd/${SERVICE_NAME}

FROM registry.gitlab.com/ninjavan/devops/golang:alpine

COPY --from=builder /build/com  /go/src/eb2b/backend/com
COPY --from=builder /build/bin/* /usr/bin/

EXPOSE 8080 8180
