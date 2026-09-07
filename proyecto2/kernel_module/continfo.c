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

/*Funcion que recolecta y genera los datos (continfo_show)*/

/** continfo_show - Callback encargado de iterar sobre el kernel y renderizar metricas en /proc */

static int continfo_show(struct seq_file *m, void *v) {
    struct sysinfo sys_info;
    struct task_struct *task;
    unsigned long total_ram, free_ram, used_ram;

    /* Obtencion de  datos   globales del sistema */
    si_meminfo(&sys_info);

    /*datos vienen expresados en paginas o unidades de memoria (mem_unit) */
    total_ram = (sys_info.totalram * sys_info.mem_unit) / (1024 * 1024);
    free_ram = (sys_info.freeram * sys_info.mem_unit) / (1024 * 1024);
    used_ram = total_ram - free_ram;

    /* Inicio de estructura JSON */
    seq_printf(m, "{\n");
    seq_printf(m, "  \"total_ram_mb\": %lu,\n", total_ram);
    seq_printf(m, "  \"free_ram_mb\": %lu,\n", free_ram);
    seq_printf(m, "  \"used_ram_mb\": %lu,\n", used_ram);
    seq_printf(m, "  \"processes\": [\n");

    bool first = true;

    /*  Recorrido de la lista circular de procesos task_struct */
    for_each_process(task) {

        /*filtramos para enfocarnos en procesos reales y contenedores.*/
        if (task->mm) {
            /* Conversion de paginas de memoria a Kilobytes */

            //task->mm->total_vm: Representa el VSZ (Memoria Virtual).
            unsigned long vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
            //get_mm_rss(task->mm): Representa el RSS (Memoria Fisica Residente).
            unsigned long rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);


            if (!first) seq_printf(m, ",\n");
            seq_printf(m, "    {\n");
            seq_printf(m, "      \"pid\": %d,\n", task->pid);
            seq_printf(m, "      \"name\": \"%s\",\n", task->comm);
            seq_printf(m, "      \"vsz_kb\": %lu, KB \n", vsz);
cat /proc/continfo_pr2_so1_201700698 | jq .            seq_printf(m, "      \"rss_kb\": %lu KB \n", rss);
            seq_printf(m, "    }");
            first = false;
        }
    }
    seq_printf(m, "\n  ]\n}\n");
    return 0;
}



//Conectar /proc con la API de operaciones (proc_ops)


static int continfo_open(struct inode *inode, struct file *file) {

    //Inicializa la estructura seq_file
    return single_open(file, continfo_show, NULL);
}

/* Enrutamiento de operaciones sobre /proc */
static const struct proc_ops continfo_fops = {
    .proc_open    = continfo_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};


//Inicializacion y Salida del Modulo

/** Constructor ejecutado al usar insmod */
static int __init continfo_init(void) {
    // Permiso 0444 = Lectura para todos los usuarios
    proc_create(PROC_FILENAME, 0444, NULL, &continfo_fops);
    pr_info("Modulo /proc/%s inicializado con exito.\n", PROC_FILENAME);
    return 0;
}

/** continfo_cleanup - Destructor ejecutado al usar 'rmmod'  */
static void __exit continfo_cleanup(void) {
    remove_proc_entry(PROC_FILENAME, NULL);
    pr_info("Modulo /proc/%s removido del kernel.\n", PROC_FILENAME);
}

module_init(continfo_init);
module_exit(continfo_cleanup);

