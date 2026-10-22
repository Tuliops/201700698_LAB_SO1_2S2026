#!/bin/bash
echo "=== 1. VERIFICANDO OPERATOR ==="
kubectl get pods -n operator-system

echo "=== 2. RECREANDO INSTANCIAS (Valkey y Grafana) ==="
kubectl delete valkeyinstance my-valkey -n default --ignore-not-found
kubectl delete grafanainstance my-grafana -n default --ignore-not-found
sleep 2
kubectl apply -f k8s/valkey-instance.yaml
kubectl apply -f k8s/grafana-instance.yaml

echo "=== 3. ESPERANDO A QUE VALKEY Y GRAFANA ESTÉN LISTOS ==="
kubectl wait --for=condition=Ready pods -l app=my-valkey-valkey -n default --timeout=120s
kubectl wait --for=condition=Ready pods -l app=my-grafana-grafana -n default --timeout=120s

echo "=== OPERATOR Y VISUALIZACIÓN LISTOS ==="