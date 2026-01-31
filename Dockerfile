# Gunakan image dasar Golang versi 1.24.5
FROM golang:1.24.5

# Tambahkan user non-root untuk keamanan
RUN useradd -m -u 1001 appuser

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh kode
COPY . .

RUN mkdir -p /app/images /app/logs /app/audio && \
    chmod -R 777 /app/images /app/logs /app/audio && \
    chown -R appuser:appuser /app/images /app/logs /app/audio

# Build aplikasi
RUN go build -o main .

# Beralih ke user non-root
USER appuser

# Set environment variables for Hugging Face Spaces (Secrets will override these)
ENV DB_HOST=""
ENV DB_PORT=""
ENV DB_USER=""
ENV DB_PASSWORD=""
ENV DB_NAME=""

# Set Server Host vars
ENV HOST_ADDRESS=0.0.0.0
ENV HOST_PORT=7860

# Expose port untuk Hugging Face Spaces
EXPOSE 7860

# Jalankan aplikasi
CMD ["./main"]