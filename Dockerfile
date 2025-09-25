# Stage 1: Build both Nuxt.js and Go application
FROM node:20-alpine AS builder

# Install pnpm and Go
RUN npm install -g pnpm && \
    apk add --no-cache go git

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy the entire project
COPY . .

# Run the build:prod command from package.json
# This will generate the Nuxt app, copy it to pocketbase dir, and build the Go binary
RUN pnpm run build:prod

# Stage 3: Create the final minimal image
FROM alpine:3.18

# Create a user to run the app
RUN addgroup -S massolit && adduser -S massolit -G massolit

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/massolit .

# Ensure correct permissions
# RUN chown -R massolit:massolit /app

# Use the non-root user
# USER massolit

# Expose the port the app runs on
EXPOSE 8090

# Command to run the application
CMD ["./massolit", "serve", "--http=0.0.0.0:8090"]
