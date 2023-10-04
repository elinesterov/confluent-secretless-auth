
load('./tilt/ko/Tiltfile', 'ko_build')

# Build spiffe-proxy
ko_build('spiffe-proxy', './spiffe-proxy')

docker_build('simple-kafka-app', './simple-kafka-app-java')
k8s_yaml('./simple-kafka-app-java/deploy/deployment.yaml')
