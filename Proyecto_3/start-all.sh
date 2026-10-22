#!/bin/bash

echo "=== 1. Iniciando Clúster y Registro ==="
docker start kind-control-plane kind-registry
sleep 5

echo "=== 2. Configurando Port-Forwards ==="
# Matar cualquier túnel viejo
pkill -f "port-forward" || true

# Gateway (puerto 8080)
nohup kubectl -n envoy-gateway-system port-forward service/envoy-default-military-gateway-33657d09 8080:80 > /dev/null 2>&1 &
echo "✅ Gateway expuesto en http://localhost:8080"

# Grafana (puerto 3000)
nohup kubectl port-forward service/my-grafana-grafana-service 3000:3000 > /dev/null 2>&1 &
echo "✅ Grafana expuesto en http://localhost:3000"

# RabbitMQ (puerto 15672)
nohup kubectl port-forward service/rabbitmq-service 15672:15672 > /dev/null 2>&1 &
echo "✅ RabbitMQ expuesto en http://localhost:15672"

echo "=== 3. Verificando Pods ==="
sleep 5
kubectl get pods

echo "=== 4. Todo listo ==="
echo "Abre Locust con: locust -f locust/locustfile.py --host=http://localhost:8080"