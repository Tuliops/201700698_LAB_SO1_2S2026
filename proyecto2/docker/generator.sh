#!/bin/bash

# Script generador de carga mediante contenedores Docker
# Crea contenedores de prueba que consumen memoria deliberadamente

CONTAINER_NAME="stress_test_$(date +%s)"

echo "[GENERATOR] Creando contenedor de carga: $CONTAINER_NAME..."

# Levanta un contenedor en segundo plano que asigna memoria
docker run -d --name "$CONTAINER_NAME" alpine sh -c "dd if=/dev/zero of=/dev/null bs=20M count=1 & sleep 60"

echo "[GENERATOR] Contenedor $CONTAINER_NAME desplegado con exito."