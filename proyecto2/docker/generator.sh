#!/bin/bash

# ==============================================================================
# SCRIPT GENERADOR DE CARGA Y MANTENIMIENTO DE CONTENEDORES - PROYECTO 2 SO1
# ==============================================================================

INTERVALO_CHEQUEO=5

echo "[GENERATOR] Iniciando generador de carga autonomo..."
echo "[GENERATOR] Regla activa: 3 contenedores Bajo Consumo (<15MB) | 2 contenedores Alto Consumo (>15MB)"
echo "[GENERATOR] Presiona [CTRL+C] para detener."
echo "------------------------------------------------------------------------------"

while true; do
  # ----------------------------------------------------------------------------
  # 1. Asegurar 3 Contenedores de Bajo Consumo (<15 MB)
  # ----------------------------------------------------------------------------
  for i in {1..3}; do
    NAME="low_app_$i"
    if ! docker ps --format '{{.Names}}' | grep -q "^${NAME}$"; then
      echo "[GENERATOR $(date +%H:%M:%S)] Levando contenedor de bajo consumo: $NAME"
      # Si existia un contenedor detenido previo con ese nombre, lo elimina primero
      docker rm -f "$NAME" >/dev/null 2>&1
      # Levanta un proceso liviano (consumo < 2 MB RAM)
      docker run -d --name "$NAME" alpine sh -c "sleep 3600" >/dev/null
    fi
  done

  # ----------------------------------------------------------------------------
  # 2. Asegurar 2 Contenedores de Alto Consumo (>15 MB)
  # ----------------------------------------------------------------------------
  for i in {1..2}; do
    NAME="high_app_$i"
    if ! docker ps --format '{{.Names}}' | grep -q "^${NAME}$"; then
      echo "[GENERATOR $(date +%H:%M:%S)] Levando contenedor de alto consumo (>15MB): $NAME"
      docker rm -f "$NAME" >/dev/null 2>&1
      # Asigna 20 MB en /dev/shm (RAM) para garantizar superar el umbral de 15 MB
      docker run -d --name "$NAME" alpine sh -c "dd if=/dev/zero of=/dev/shm/fill bs=1M count=20 && sleep 3600" >/dev/null
    fi
  done

  # Esperar antes del siguiente ciclo de verificacion
  sleep $INTERVALO_CHEQUEO
done
