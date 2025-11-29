# Frontend Dockerfile for TestGen Vue.js application

# Build stage
FROM node:24.11.1-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments for Vite environment variables
ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ARG VITE_MAX_FILE_SIZE=52428800
ARG VITE_SUPPORTED_FORMATS=pdf,docx,pptx,txt

# Set environment variables for build
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_MAX_FILE_SIZE=${VITE_MAX_FILE_SIZE}
ENV VITE_SUPPORTED_FORMATS=${VITE_SUPPORTED_FORMATS}

# Build application (skip type checking for Docker build)
RUN npx vite build

# Production stage
FROM node:24.11.1-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ALL dependencies (vite preview needs devDependencies)
RUN npm ci

# Copy built files from build stage
COPY --from=build /app/dist ./dist

# Copy vite config for preview
COPY vite.config.ts ./

# Expose port
EXPOSE 5173

# Start application in preview mode
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "5173"]
