# UI-only Docker build for development
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install bun
RUN npm install -g bun@1.2.7

# Copy package files
COPY ui/package.json ui/bun.lock* ./

# Install dependencies
RUN bun install

# Copy source code
COPY ui/ .

# Build the application
RUN bun run build

# Install a simple static server
RUN npm install -g serve

# Expose port
EXPOSE 3000

# Command to serve the built application
CMD ["serve", "-s", "dist", "-p", "3000"]