FROM archlinux:latest

RUN pacman -Syy
RUN pacman -S --noconfirm mariadb go

WORKDIR /go/src/srv
COPY . .
RUN go build

EXPOSE 80

CMD ["./srv"]
