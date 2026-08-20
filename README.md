# Universidad: Universidad de San Carlos de Guatemala (USAC)

# Facultad: Facultad de Ingeniería

# Curso: Sistemas Operativos 
# Proyecto 1: Desarrollo, Conexión y Gestión de Contenedores en Entornos Virtualizados

El presente proyecto tiene como objetivo principal el diseño e implementación de un
entorno virtualizado que integre el uso de máquinas virtuales (VMs) y contenedores,
empleando tecnologías modernas como Docker, Containerd, Podman, Go y Zot. Esta
arquitectura permite simular entornos reales de desarrollo utilizados en la industria,
enfocados en la contenerización, el almacenamiento de imágenes de contenedores,
conexión entre API´s y la gestión eficiente de recursos.

---

##  Recursos y Herramientas a Utilizar

Para el desarrollo exitoso del proyecto, se integra un ecosistema de herramientas de contenedores y virtualización orientadas a entornos de producción y pruebas locales:

*   **Docker:** Plataforma de contenedores utilizada para la creación de imágenes base y pruebas iniciales de despliegue de microservicios.
   <img width="491" height="407" alt="image" src="https://github.com/user-attachments/assets/fd4ad5c3-1814-448a-b902-e96b33df089f" />

*   **Zot:** Registro de contenedores OCI (Open Container Initiative) de alto rendimiento, seguro y ligero, utilizado para la distribución y almacenamiento local/remoto de imágenes.

  <img width="447" height="447" alt="image" src="https://github.com/user-attachments/assets/dded5730-0415-4b01-9cde-b5f150bae770" />

*   **qemu-kvm:** Hipervisor basado en kernel de Linux para la virtualización completa de hardware en los entornos virtuales asignados.
   <img width="132" height="109" alt="image" src="https://github.com/user-attachments/assets/ff6afc0b-f2aa-4246-8d63-40323a5c2a6b" />

*   **ContainerD:** Motor de contenedores de nivel industrial que proporciona el runtime subyacente para la ejecución y ciclo de vida de los contenedores.
   <img width="759" height="188" alt="image" src="https://github.com/user-attachments/assets/c132347e-7f6c-43c5-90ef-bf8d3f259ad7" />

*   **Go (Golang):** Lenguaje de programación principal empleado para el desarrollo de servicios backend, clientes API o herramientas CLI requeridas en la lógica del proyecto.
   <img width="735" height="272" alt="image" src="https://github.com/user-attachments/assets/7db214c3-ce63-4500-9bea-c8c40f87bbbb" />

*   **Podman:** Herramienta sin demonio (*daemonless*) para la gestión y ejecución de contenedores y pods compatibles con OCI, utilizada como alternativa robusta y segura a Docker en las máquinas virtuales.
<img width="884" height="293" alt="image" src="https://github.com/user-attachments/assets/1d7415de-a8f2-4abc-853c-ae215a828565" />

>  **Referencia técnica de apoyo:** Para profundizar en las diferencias y arquitectura entre runtimes en entornos Cloud Native, puedes consultar la guía técnica en [Container Runtimes para Proyectos Cloud Native](https://curzona.com/blog/tech-blog-1/container-runtimes-para-proyectos-cloud-native-4).

---

## 📋 Requisitos e Instalación

Asegúrate de contar con las siguientes herramientas instaladas en tu entorno de desarrollo o máquina virtual (VM):

1. **Go** (versión 1.18 o superior)
2. **Podman** y/o **Docker**
3. **QEMU/KVM** configurado en el sistema operativo base (Linux)

