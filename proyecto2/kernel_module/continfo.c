/*
Cabeceras del Kernel (Headers)
+*/
#include <linux/module.h>      // Requerido por todos los modulos de kernel
#include <linux/kernel.h>      // Funciones del kernel como pr_info()
#include <linux/init.h>        // Macros __init y __exit para gestion de memoria
#include <linux/proc_fs.h>     // Funciones para crear y eliminar entradas en /proc
#include <linux/seq_file.h>    // API seq_file para manejo seguro de lectura de datos
#include <linux/sched/signal.h>// Necesario para iterar procesos con for_each_process()
#include <linux/sched.h>       // Definicion de la estructura task_struct
#include <linux/mm.h>          // Calculo de memoria RAM, RSS y VSZ


// Info del modulo
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tulio P Sistemas Operativos 1");
MODULE_DESCRIPTION(" Sistemas Operativos 1 - Registro de procesos y memoria");
MODULE_VERSION("1.0");

// Nombre del Modulo de Procesos de Contenedores
#define PROC_FILENAME "continfo_pr2_so1_201700698"