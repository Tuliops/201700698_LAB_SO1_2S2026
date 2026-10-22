# GUÍA DE INSTALACIÓN

GUÍA DE INSTALACIÓN Y DESPLIEGUE: M.U.M.N.K8s
===

Esta guía detalla paso a paso cómo preparar el entorno, desplegar todos los componentes del sistema M.U.M.N.K8s y realizar las pruebas funcionales en un entorno local (**Linux Mint**) utilizando **Kind** como clúster de Kubernetes y **Zot** como registro local de imágenes.

---

## 1. Requisitos Previos

Asegúrate de tener instaladas las siguientes herramientas en tu sistema operativo:

* **Docker** (Instalado y en ejecución).
* **Kubernetes CLI** (`kubectl`).
* **Helm** (Gestor de paquetes para Kubernetes).
* **Kind** (Herramienta para crear clústeres en Docker).
* **Operator SDK**.
* **Go** (Versión 1.22 o superior).
* **Locust** (`pip install locust`).

---

## 2. Infraestructura Base (Kind y Zot)

### 2.1. Crear el Clúster de Kind con mapeo de puertos
El Gateway API necesita escuchar en el puerto 80 de tu máquina anfitriona.

1. Crea un archivo llamado `kind-config.yaml` en la raíz de tu proyecto con el siguiente contenido:
   ```yaml
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   containerdConfigPatches: []
   nodes:
     - role: control-plane
       extraPortMappings:
         - containerPort: 80
           hostPort: 80
           protocol: TCP
   ```

2. Crea el clúster ejecutando en la terminal:
   ```bash
   kind create cluster --config kind-config.yaml --name kind
   ```

### 2.2. Levantar el Registro Local (Zot)
Usaremos un registro `registry:2` estándar en el puerto 5001 para almacenar nuestras imágenes.

1. Descarga y levanta el contenedor del registro:
   ```bash
   docker run -d --restart=always -p "127.0.0.1:5001:5000" --name kind-registry registry:2
   ```

2. Conectar Kind con Zot (configurar como registro inseguro local):
   * Ejecutar el script de configuración de `/etc/containerd/certs.d/` y conectar la red:
   ```bash
   docker network connect "kind" "kind-registry"
   ```

---

## 3. Instalación de Kubernetes Gateway API (Envoy)

1. Instalar Envoy Gateway usando Helm:
   ```bash
   helm install eg oci://docker.io/envoyproxy/gateway-helm --version v1.1.0 -n envoy-gateway-system --create-namespace
   ```

2. Esperar a que el controlador esté listo:
   ```bash
   kubectl wait --timeout=5m -n envoy-gateway-system deployment/envoy-gateway --for=condition=Available
   ```

3. Crear el `GatewayClass` necesario:
   ```bash
   kubectl apply -f - <<EOF
   apiVersion: gateway.networking.k8s.io/v1
   kind: GatewayClass
   metadata:
     name: eg
   spec:
     controllerName: gateway.envoyproxy.io/gatewayclass-controller
   EOF
   ```

---

## 4. Compilación y Despliegue de Microservicios

### 4.1. Compilación de Imágenes Docker
Desde la raíz del proyecto, compila cada servicio y súbelo a Zot (`localhost:5001`).

* **API en Rust:**
  ```bash
  docker build -t localhost:5001/api-rust:v2 . && docker push localhost:5001/api-rust:v2
  ```
* **Go Deployment 1 (API y gRPC Client):** 
  * Compilar y subir `go-d1-api:v1` y `go-d1-grpc:v1`.
* **Go Deployment 2 (gRPC Server y RabbitMQ Sidecar):** 
  * Compilar y subir `go-d2-grpc:v3` y `go-d2-sidecar:v2`.
* **Consumer (Go):** 
  * Compilar y subir `go-consumer:v12`.

### 4.2. Despliegue de Manifiestos en Kubernetes
Aplica todos los manifiestos YAML para crear los Deployments, Services y HPA:
```bash
kubectl apply -f k8s/rust-api.yaml
kubectl apply -f k8s/go-deployment-1.yaml
kubectl apply -f k8s/go-deployment-2.yaml
kubectl apply -f k8s/consumer.yaml
kubectl apply -f k8s/rabbitmq.yaml
kubectl apply -f k8s/gateway.yaml
```
Verifica que los pods estén corriendo con: 
```bash
kubectl get pods
```

### 4.3. Configurar el HPA (Metrics Server)
Kind no trae Metrics Server por defecto.

1. Instalar Metrics Server:
   ```bash
   kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
   ```

2. Parchear para que funcione en local (evitar errores TLS):
   ```bash
   kubectl patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
   ```

3. Aplicar el HPA de Rust:
   ```bash
   kubectl apply -f k8s/rust-hpa.yaml
   ```

---

## 5. Despliegue del Operator (Valkey y Grafana)

1. Compilar y subir la imagen del Operator:
   ```bash
   cd operator && docker build -t localhost:5001/operator:v9 . && docker push localhost:5001/operator:v9
   ```

2. Desplegar el Operator en el clúster:
   ```bash
   make deploy
   ```

3. Aplicar los Custom Resources (CRs) para instanciar Valkey y Grafana:
   ```bash
   kubectl apply -f k8s/valkey-instance.yaml
   kubectl apply -f k8s/grafana-instance.yaml
   ```

4. Verificar que el Operator creó los pods:
   ```bash
   kubectl get pods -l app=my-valkey-valkey
   ```

---

## 6. Instalación de Seguridad (Falco)

Añadir el repositorio e instalar Falco con la configuración para Kind:
```bash
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update
helm install falco falcosecurity/falco --namespace falco --create-namespace --set driver.kind=modern_ebpf --set falcoctl.artifact.install.enabled=false --set falcoctl.artifact.follow.enabled=false -f k8s/falco-custom-rules.yaml
```

---

## 7. Pruebas de Funcionalidad (Verificación)

### 7.1. Abrir los Túneles (Port-Forward)
Abre terminales separadas para cada uno:

* **Gateway (Entrada):**
  ```bash
  kubectl -n envoy-gateway-system port-forward service/envoy-default-military-gateway-33657d09 8080:80
  ```
* **Grafana (Dashboard):**
  ```bash
  kubectl port-forward service/my-grafana-grafana-service 3000:3000
  ```
* **RabbitMQ (Admin):**
  ```bash
  kubectl port-forward service/rabbitmq-service 15672:15672
  ```

### 7.2. Generar Tráfico con Locust
1. Ejecuta:
   ```bash
   locust -f locust/locustfile.py --host=http://localhost:8080
   ```
2. Abre `http://localhost:8089` en tu navegador, configura **50 usuarios** y haz clic en **Start swarming**.

### 7.3. Verificar Grafana y Falco

1. Abre `http://localhost:3000` (`admin/admin`) y observa el dashboard.
2. Simula amenaza:
   ```bash
   kubectl exec -it \$POD_RUST -- cat /etc/shadow > /dev/null
   ```
3. Revisa alerta:
   ```bash
   kubectl logs -n falco -l app.kubernetes.io/name=falco --tail=10
   ```
