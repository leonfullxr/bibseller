# SvelteKit (adapter-node) web server. Build context is frontend/.
# Build stage runs the full install + build, then prunes to prod deps; the
# runtime stage carries only build/ + production node_modules.
# syntax=docker/dockerfile:1
FROM node:26-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# Origin-specific build: vite's paths.assets bakes this in so prod asset URLs
# are absolute and Vite's __vitePreload deps are not resolved relative to the
# entry module (which doubled the path -> no CSS/JS, #66). Must be passed for
# prod images; empty falls back to dev behavior.
ARG PUBLIC_ORIGIN
ENV PUBLIC_ORIGIN=$PUBLIC_ORIGIN
RUN npm run build && npm prune --omit=dev

FROM node:26-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
COPY --from=build /app/build ./build
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./package.json
EXPOSE 3000
USER node
CMD ["node", "build/index.js"]
