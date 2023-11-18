FROM maven:3.9-amazoncorretto-21 AS build
COPY . .
RUN mvn clean package -DskipTests

FROM openjdk:21-jdk-slim
COPY --from=build /target/chat-0.0.1-SNAPSHOT.jar chat.jar
EXPOSE 8080
ENTRYPOINT ["JAVA", "-jar", "chat.jar"]
