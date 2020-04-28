FROM golang:1.14-alpine as build
WORKDIR /src/pageViewCounter
COPY . .
RUN go build

FROM alpine:latest
COPY --from=build /src/pageViewCounter .
CMD ["./pageViewCounter", "-port", "3000"]
EXPOSE 3000
