#!/bin/bash
echo "=== 1. INICIANDO CLÚSTER Y REGISTRO ==="
docker start kind-control-plane kind-registry 2>/dev/null || true
sleep 5

echo "=== 2. ESPERANDO A QUE EL NODO ESTÉ LISTO ==="
kubectl wait --for=condition=Ready nodes/kind-control-plane --timeout=60s

echo "=== 3. APLICANDO MANIFIESTOS BASE ==="
kubectl apply -f k8s/rust-api.yaml
kubectl apply -f k8s/go-deployment-1.yaml
kubectl apply -f k8s/go-deployment-2.yaml
kubectl apply -f k8s/consumer.yaml
kubectl apply -f k8s/rabbitmq.yaml
kubectl apply -f k8s/gateway.yaml
kubectl apply -f k8s/rust-hpa.yaml

echo "=== 4. ESPERANDO A QUE LOS PODS ESTÉN LISTOS ==="
kubectl wait --for=condition=Ready pods --all -n default --timeout=180s || true

echo "=== INFRAESTRUCTURA BASE LISTA ==="