# Iniciar o golang imagem
FROM golang:alpine as builder

# Enable go modules
ENV GO111MODULE=on

# Install git. 
RUN apk update && apk add --no-cache git

# Set current working directory
WORKDIR /app

# Copiar go mod e sum files
COPY go.mod ./
COPY go.sum ./

# Download de todas as dependencias.
RUN go mod download

# Copia todos os arquivos 
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

# Finally our multi-stage to build a small image
# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file
COPY --from=builder /app/bin/main .

# Execução
CMD ["./main"]