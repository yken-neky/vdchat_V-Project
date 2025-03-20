# Etapa de construcción
FROM golang:1.23.6-alpine AS builder
WORKDIR /src

# Copia los archivos de módulos y descarga las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Construye la aplicación; se usa CGO_DISABLED para generar un binario estático
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd

# Etapa final, con una imagen más ligera
FROM alpine
WORKDIR /src

# Copia los recursos y el binario desde la etapa de construcción
COPY --from=builder /bin/app /bin/app
COPY --from=builder /src/views /src/views
COPY --from=builder /src/assets /src/assets

ENTRYPOINT ["/bin/app"]
