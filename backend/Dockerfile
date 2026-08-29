FROM golang:1.27.0 AS build
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o out.exe

FROM alpine:3.24.1
WORKDIR /app
COPY --from=build /app/out.exe /app/out.exe
EXPOSE 8080
CMD [ "/app/out.exe" ]