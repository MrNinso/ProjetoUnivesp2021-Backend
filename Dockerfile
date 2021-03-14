FROM golang:alpine AS build

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o server.exe

FROM scratch AS bin

COPY --from=build /src/server.exe /

CMD ["/server.exe"]