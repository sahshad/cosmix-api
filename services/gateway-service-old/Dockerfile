# Build stage
FROM node:20-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy tsconfig.json
COPY tsconfig.json ./

# Copy rest of the code
COPY src ./src

# Copy .env
COPY .env .env

# Build the application
RUN npm run build

# Expose port
EXPOSE 5000

# Run the application
CMD ["node", "dist/server.js"]
