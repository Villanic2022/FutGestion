# Build stage
FROM node:20-alpine AS build

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source code and build
COPY . .
RUN npm run build

# Serve stage using Nginx
FROM nginx:alpine

# Copy the custom nginx setting if you have one, or just the default
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built files to nginx document root
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Run nginx
CMD ["nginx", "-g", "daemon off;"]
