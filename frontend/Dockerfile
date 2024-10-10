FROM node:16
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
RUN npm install -g http-server
COPY . .
EXPOSE 3000
CMD ["http-server", "./", "-p", "3000"]
