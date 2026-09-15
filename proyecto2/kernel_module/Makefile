obj-m += continfo.o

KDIR := /home/tj/wsl-kernel
PWD := $(shell pwd)
PROC_PATH := /proc/continfo_pr2_so1_201700698

all:
	@echo "[+] Compilando modulo de kernel..."
	$(MAKE) -C $(KDIR) M=$(PWD) KBUILD_MODPOST_WARN=1 modules

clean:
	@echo "[+] Limpiando archivos binarios y temporales..."
	$(MAKE) -C $(KDIR) M=$(PWD) clean
	@rm -f *.o *.ko *.mod *.mod.c *.symvers *.order .*.cmd

load: all
	@echo "[+] Removiendo version previa si esta cargada..."
	-sudo rmmod continfo 2>/dev/null || true
	@echo "[+] Cargando modulo de kernel en el sistema..."
	sudo insmod continfo.ko
	@echo "[+] Modulo continfo.ko cargado exitosamente."

unload:
	@echo "[+] Desmontando modulo de kernel..."
	-sudo rmmod continfo 2>/dev/null || true
	@echo "[+] Modulo removido correctamente."

# Regla de prueba rapida para verificar la interfaz /proc
test:
	@echo "\n=== VERIFICACION DE /proc/continfo_pr2_so1_201700698 ==="
	@if [ -f $(PROC_PATH) ]; then \
		echo "[SUCCESS] El archivo proc existe. Mostrando estructura JSON:"; \
		cat $(PROC_PATH) | head -n 30; \
	else \
		echo "[ERROR] El archivo $(PROC_PATH) NO existe. Asegurate de ejecutar 'make load'."; \
	fi