

# Objeto final que genera el compilador (continfo.o -> continfo.ko)
obj-m += continfo.o

# Ubicacion de las cabeceras y reglas del kernel actual en tu distribucion Linux
KDIR := /lib/modules/$(shell uname -r)/build
PWD := $(shell pwd)

all:
	# Llama al Makefile oficial del Kernel pasando nuestra ruta como modulo externo (M=)
	$(MAKE) -C $(KDIR) M=$(PWD) modules

clean:
	# Limpia archivos temporales (.o, .mod, .ko)
	$(MAKE) -C $(KDIR) M=$(PWD) clean

