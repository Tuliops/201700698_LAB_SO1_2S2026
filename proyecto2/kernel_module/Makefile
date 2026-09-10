MODULE_NAME = continfo
# Ruta absoluta directa hacia tu usuario (evita la resolucion a /root bajo sudo)
KERNEL_DIR = /home/tj/wsl-kernel
PROC_FILE = /proc/continfo_pr2_so1_201700698

# Obtener el directorio actual de forma dinamica dentro de Make
CURRENT_DIR := $(shell pwd)

obj-m += $(MODULE_NAME).o

all:
	make -C $(KERNEL_DIR) M=$(CURRENT_DIR) SHELL=/bin/bash KERNELRELEASE="6.18.33.2-microsoft-standard-WSL2" LOCALVERSION="" KBUILD_MODPOST_WARN=1 modules

clean:
	make -C $(KERNEL_DIR) M=$(CURRENT_DIR) clean

load: all
	@if lsmod | grep -q $(MODULE_NAME); then \
		echo "El modulo ya esta cargado. Removiendo previa version..."; \
		sudo rmmod $(MODULE_NAME); \
	fi
	sudo insmod $(MODULE_NAME).ko
	@echo "Modulo cargado exitosamente."

unload:
	@if lsmod | grep -q $(MODULE_NAME); then \
		sudo rmmod $(MODULE_NAME); \
		echo "Modulo removido correctamente."; \
	else \
		echo "El modulo no se encuentra cargado."; \
	fi

test:
	@if [ -f $(PROC_FILE) ]; then \
		cat $(PROC_FILE) | head -n 30; \
	else \
		echo "Error: El archivo $(PROC_FILE) no existe. Ejecute 'make load' primero."; \
	fi