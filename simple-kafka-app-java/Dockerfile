# Builder stage
FROM maven:3.8.1-openjdk-11-slim AS builder
WORKDIR /app
COPY pom.xml .
# Download all required dependencies into one layer
RUN mvn dependency:go-offline -B
COPY src ./src
RUN mvn package -DskipTests

# Application stage
FROM openjdk:22-slim
WORKDIR /app
COPY --from=builder /app/target/your-artifact-1.0-SNAPSHOT.jar ./app.jar
EXPOSE 8080
CMD ["java", "-jar", "/app/app.jar"]