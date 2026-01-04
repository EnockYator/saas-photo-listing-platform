# Stage 1: Build React App
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy frontend code
COPY . .

# Install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Build the React appts
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:1.25-alpine


# Copy built React app to Nginx html folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom Nginx config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
