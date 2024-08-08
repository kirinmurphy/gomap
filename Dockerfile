## -- GO INIT ------- 

FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

## -- YARN INIT ------- 
FROM node:18 as frontend-builder

WORKDIR /app

COPY package.json yarn.lock* package-lock.json* ./

RUN if [ -f yarn.lock ]; then yarn install; else npm install; fi

COPY . .

RUN npm run build 

FROM golang:1.21

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=frontend-builder /app/templates ./templates

EXPOSE 8080

CMD ["./main"]
