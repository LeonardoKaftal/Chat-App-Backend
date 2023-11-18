FROM maven:4.0.0-openjdk-21 AS build
COPY . .
RUN mvm clean package -DskipTests

FROM openjdk:21-jdk-slim
COPY --from=build /target/chat-0.0.1-SNAPSHOT.jar chat.jar
EXPOSE 8080
ENTRYPOINT ["JAVA", "-jar", "chat.jar"]