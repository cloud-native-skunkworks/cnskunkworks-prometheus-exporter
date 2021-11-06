docker-push:
	docker push tibbar/cnskunkworks-prometheus-exporter
docker-build:
	docker build -t tibbar/cnskunkworks-prometheus-exporter .
deploy:
	helm upgrade --install cnskunkworks-prometheus-exporter ./chart --namespace=default