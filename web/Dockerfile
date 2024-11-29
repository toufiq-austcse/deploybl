# Build Stage
FROM node:18-alpine AS BUILD_IMAGE
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build


# Production Stage
FROM node:18-alpine AS PRODUCTION_STAGE
WORKDIR /app
COPY --from=BUILD_IMAGE /app/package*.json ./
COPY --from=BUILD_IMAGE /app/.next/static ./.next/static
COPY --from=BUILD_IMAGE /app/.next/standalone .
COPY --from=BUILD_IMAGE /app/public ./public
ENV NODE_ENV=production
ENTRYPOINT ["node", "server.js"]