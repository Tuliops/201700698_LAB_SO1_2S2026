# Manual Técnico: Sistema M.U.M.N.K8s

## 1. Arquitectura General

El sistema **M.U.M.N.K8s** es una arquitectura de microservicios distribuida y orientada a eventos, desplegada sobre un clúster local de Kubernetes (usando **Kind**). Su propósito es simular la recepción, procesamiento, almacenamiento y visualización de reportes militares en tiempo real.

### Componentes de la Arquitectura
1. **Generador de Tráfico (Locust):** Envía reportes militares simulados en formato JSON mediante peticiones HTTP.
2. **Gateway API (Envoy):** Actúa como punto de entrada único, exponiendo la ruta `/grpc-201700698`.
3. **API REST (Rust):** Primer punto de contacto, valida el JSON y lo reenvía internamente.
4. **Servicios en Go (Patrón Sidecar):**
    *   **Deployment 1 (Cliente):** Recibe de Rust y actúa como cliente gRPC.
    *   **Deployment 2 (Servidor/Writer):** Recibe por gRPC y publica en el broker mediante un sidecar.
5. **Broker (RabbitMQ):** Desacopla la ingesta del procesamiento mediante la cola `war_reports`.
6. **Consumidor (Go):** Extrae mensajes de RabbitMQ y los procesa.
7. **Almacenamiento (Valkey):** Base de datos en memoria, gestionada por un Operator.
8. **Visualización (Grafana):** Lee de Valkey y muestra el dashboard.
9. **Seguridad (Falco):** Monitorea syscalls y detecta accesos a archivos sensibles.
10. **Registro (Zot):** Almacena todas las imágenes Docker construidas localmente.

---

## 2. Flujo Completo de Datos

1. **Locust** genera una petición POST con el JSON del reporte militar hacia el Gateway.
2. **Envoy Gateway** recibe la petición en la ruta `/grpc-201700698` y la redirige al servicio de la **API Rust**.
3. **API Rust** recibe el JSON, lo procesa y realiza una petición HTTP al **Go Deployment 1**.
4. **Go Deployment 1 (Contenedor gRPC Client)** toma los datos y realiza una llamada gRPC al **Go Deployment 2**.
5. **Go Deployment 2 (Contenedor gRPC Server)** recibe la llamada y la pasa por `localhost` al **Sidecar (RabbitMQ Writer)**.
6. **Sidecar** publica el mensaje en la cola `war_reports` de **RabbitMQ**.
7. **Go Consumer** está suscrito a la cola, toma el mensaje, extrae los datos y los guarda en **Valkey**.
8. **Grafana** consulta periódicamente a Valkey para actualizar los paneles del dashboard.

---

## 3. Configuración de Gateway API

En sustitución de un Ingress Controller tradicional, se utilizó **Envoy Gateway** para exponer los servicios.

*   **Instalación:** Se instaló el CRD de Gateway API y el controlador de Envoy Gateway mediante Helm.
*   **Ruta Configurada:** Se creó un `HTTPRoute` que expone la ruta exacta `/grpc-201700698` y la redirige al servicio `rust-api-service` en el puerto 8080.

---

## 4. Comunicación REST y gRPC

*   **REST (Externo):** Locust se comunica con Rust mediante HTTP/JSON a través del Gateway.
*   **gRPC (Interno):** La comunicación entre los servicios de Go se realiza mediante Protocol Buffers (gRPC), definidos en el archivo `report.proto`.
    *   **Contrato:** El enum `Countries` asegura que solo se envíen códigos de país válidos (USA, RUS, CHN, ESP, GTM).

---

## 5. Uso de RabbitMQ

RabbitMQ actúa como el broker de mensajería principal.
*   **Cola:** `war_reports`.
*   **Flujo:** El Writer (Sidecar de Go D2) publica en la cola. El Consumer (Go) consume los mensajes con `auto-ack=false` para garantizar que se guardan en Valkey antes de borrarlos de la cola.

---

## 6. Diseño y Funcionamiento del Operator

Se desarrolló un **Operator basado en Helm** utilizando Operator SDK. Su función es gestionar el ciclo de vida de Valkey y Grafana.

*   **CRDs Definidos:**
    *   `ValkeyInstance`: Define el despliegue de Valkey (Deployment, Service, PVC).
    *   `GrafanaInstance`: Define el despliegue de Grafana (Deployment, Service, Datasource).
*   **Lógica de Reconciliación:** Cuando se aplica un Custom Resource (CR), el Operator detecta el cambio, renderiza el Helm chart y crea/actualiza los recursos en el clúster.

---

## 7. Configuración de HPA (Horizontal Pod Autoscaler)

Se configuró un HPA para el Deployment de la API Rust.
*   **Métrica:** Uso de CPU.
*   **Umbral:** 30%.
*   **Rango:** Mínimo 1 réplica, máximo 3 réplicas.
*   **Prueba:** Al ejecutar Locust con alta carga, el HPA escaló automáticamente el número de réplicas.

---

## 8. Publicación y Consumo de Imágenes desde Zot

*   **Zot Local:** Se desplegó un registro Zot en un contenedor Docker fuera del clúster de Kind.
*   **Integración:** Se configuró el nodo de Kind para confiar en el registro local (`/etc/containerd/certs.d/`).
*   **Publicación:** Todas las imágenes (Rust, Go, Operator) se etiquetaron como `localhost:5001/<nombre>:vX` y se subieron a Zot.

---

## 9. Seguridad en Tiempo Real con Falco

Se instaló Falco con el driver `modern_ebpf`. Se creó una regla personalizada para detectar la lectura de archivos sensibles (`/etc/shadow`). Al ejecutar `cat /etc/shadow` dentro del pod de Rust, Falco disparó una alerta de seguridad inmediata, demostrando la capacidad de auditoría del sistema.

---

## 10. Pruebas Realizadas y Conclusiones

### 10.1 Pruebas de Carga (Locust)
Se ejecutaron pruebas con Locust enviando reportes militares simulados. El sistema procesó miles de peticiones sin fallos (0% de errores).

### 10.2 Análisis de Réplicas (1 vs 2 réplicas)
Se realizaron pruebas comparativas en los componentes de Go. Con 2 réplicas, el sistema mostró una mejor tolerancia a fallos y capacidad de respuesta bajo carga concurrente.

### 10.3 Persistencia
Se verificó que los datos sobreviven al reinicio de los pods de Valkey gracias al PVC gestionado por el Operator.

### Conclusiones
El proyecto demostró la implementación exitosa de una arquitectura de microservicios escalable y resiliente en Kubernetes. El uso del Operator simplificó la gestión de servicios complejos como Valkey y Grafana, mientras que la integración de Gateway API y gRPC permitió una comunicación eficiente y moderna entre componentes. La observabilidad y seguridad se garantizaron con Grafana y Falco respectivamente.