## -- GO BUILD ------- 

  FROM golang:1.21 as builder

  WORKDIR /app
  
  COPY go.mod go.sum ./
  RUN go mod download
  
  COPY ./src ./src
  
  WORKDIR /app/src
  
  RUN go test -count=1 ./... && go build -o /app/main .
  
  ## -- YARN BUILD ------- 
  FROM node:18 as frontend-builder
  
  WORKDIR /app
  
  COPY package.json yarn.lock* package-lock.json* ./
  COPY tailwind.config.js postcss.config.js ./
  
  RUN if [ -f yarn.lock ]; then yarn install; else npm install; fi
  
  COPY ./src ./src
  
  RUN npm run build
  
  ## -- ASSEMBLE --------
  FROM golang:1.21
  
  WORKDIR /app
  
  COPY --from=builder /app/main .
  COPY --from=frontend-builder /app/src/templates ./src/templates
  
  EXPOSE 8080
  
  CMD ["./main"]
  