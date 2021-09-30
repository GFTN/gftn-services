FROM openjdk:8

RUN mkdir -p /opt/integrationtests
WORKDIR /opt/integrationtests
COPY . /opt/integrationtests
RUN ./gradlew compileJava \
        && ./gradlew compileGroovy \
        && ./gradlew compileTestJava \
        && ./gradlew compileTestGroovy



