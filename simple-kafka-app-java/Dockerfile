# Builder stage
FROM maven:3.8.1-openjdk-11-slim AS builder
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY pom.xml .
# Download all required dependencies into one layer
RUN mvn dependency:go-offline -B
COPY src ./src
COPY checkstyle ./checkstyle
COPY checkstyle.xml .
RUN mvn package -DskipTests

# Application stage
FROM openjdk:22-slim
WORKDIR /app
COPY --from=builder /app/target/client-7.4.0.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/app/app.jar"]