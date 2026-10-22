#!/bin/bash

echo "=================================================="
echo "   VERIFICACIÓN PASO A PASO - RÚBRICA PROYECTO 3"
echo "=================================================="
echo ""

echo "--- 1.1 Arquitectura general, despliegue (Gateway API) ---"
kubectl get nodes -o wide
echo ">>> Ruta expuesta (debe ser /grpc-201700698):"
kubectl get httproute rust-route -n default -o jsonpath='{.spec.rules[*].matches[*].path}'
echo ""
kubectl get gateway,httproute -n default
echo ""

echo "--- 1.2 API REST en Rust ---"
kubectl get hpa,deployments -n default
echo ""

echo "--- 1.3 Servicios en Go (API/gRPC Client, gRPC Server/Writers) ---"
echo ">>> Logs de Rust (REST):"
kubectl logs -l app=rust-api -n default --tail=5
echo ">>> Logs de Go D1 (gRPC Client):"
kubectl logs -l app=go-d1 -c grpc-client -n default --tail=5
echo ""

echo "--- 1.4 Arquitectura general, despliegue en Zot ---"
echo ">>> Catálogo de imágenes en Zot:"
curl -s http://localhost:5001/v2/_catalog
echo ""
echo ""

echo "--- 1.5 Message broker y consumidor RabbitMQ ---"
kubectl get pods -n default
echo ">>> Mensajes en la cola de RabbitMQ:"
kubectl exec -it deployment/rabbitmq -n default -- rabbitmqctl list_queues name messages
echo ""

echo "--- 1.6 Almacenamiento (Valkey) ---"
kubectl get crd | grep valkey
kubectl get valkeyinstances -n default
echo ">>> Datos persistidos en Valkey (DBSIZE):"
kubectl exec -it deployment/my-valkey-valkey -n default -- valkey-cli DBSIZE
echo ""

echo "--- 1.8 Falco (Seguridad) ---"
echo ">>> Simulando amenaza: cat /etc/shadow en pod de Rust..."
POD_RUST=$(kubectl get pods -l app=rust-api -n default -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POD_RUST -n default -- cat /etc/shadow > /dev/null 2>&1
echo ">>> Buscando alerta en Falco..."
kubectl logs -n falco -l app.kubernetes.io/name=falco --tail=5 | grep -i "shadow" || echo "⚠️ Espera 5 segundos y revisa de nuevo."
echo ""

echo "=================================================="
echo "   FIN DE LA VERIFICACIÓN"
echo "=================================================="