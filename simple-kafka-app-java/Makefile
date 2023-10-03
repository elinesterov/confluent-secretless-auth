.PHONY: build clean test docker

build:
	mvn package -DskipTests

test:
	mvn test

clean:
	mvn clean

docker:
	docker build -t my-application .