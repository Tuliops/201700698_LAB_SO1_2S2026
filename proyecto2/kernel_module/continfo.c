/*
Cabeceras del Kernel (Headers)
+*/
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/sched/mm.h>
#include <linux/mm.h>
#include <linux/sched/cputime.h>

#define PROC_NAME "continfo_pr2_so1_201700698"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("201700698");
MODULE_DESCRIPTION("Modulo de Kernel SO1 -  CPU, Memoria y E/S en formato JSON");

/*
 * Funcion principal que genera la salida en /proc
 */
static int continfo_show(struct seq_file *m, void *v) {
    struct task_struct *task;
    bool first = true;

    seq_puts(m, "[\n");

    /* Bloqueo de lectura RCU para iteracion segura en la lista de procesos */
    rcu_read_lock();
    for_each_process(task) {
        unsigned long vsz_pages = 0;
        unsigned long rss_pages = 0;
        unsigned long vsz_kb = 0;
        unsigned long rss_kb = 0;

        u64 utime = 0, stime = 0;
        u64 total_cpu_jiffies = 0;

        u64 read_bytes = 0;
        u64 write_bytes = 0;

        /* Validacion de puntero de memoria para omitir hilos de kernel sin mm */
        if (task->mm) {
            vsz_pages = task->mm->total_vm;
            rss_pages = get_mm_rss(task->mm);
        }

        /* Conversion de paginas de memoria a Kilobytes */
        vsz_kb = vsz_pages * (PAGE_SIZE / 1024);
        rss_kb = rss_pages * (PAGE_SIZE / 1024);

        /* Obtencion del tiempo de CPU en tiempo de usuario y sistema */
        task_cputime_adjusted(task, &utime, &stime);
        total_cpu_jiffies = utime + stime;

        /* Obtencion de contadores de entrada y salida de disco */
#ifdef CONFIG_TASK_IO_ACCOUNTING
        read_bytes = task->ioac.read_bytes;
        write_bytes = task->ioac.write_bytes;
#endif

        if (!first) {
            seq_puts(m, ",\n");
        }
        first = false;

        seq_puts(m, "  {\n");
        seq_printf(m, "    \"pid\": %d,\n", task->pid);
        seq_printf(m, "    \"name\": \"%s\",\n", task->comm);
        seq_printf(m, "    \"cpu_jiffies\": %llu,\n", total_cpu_jiffies);

        /* Formateo dinamico de Memoria Virtual (VSZ) en KB o MB */
        if (vsz_kb < 1024) {
            seq_printf(m, "    \"vsz\": \"%lu KB\",\n", vsz_kb);
        } else {
            seq_printf(m, "    \"vsz\": \"%lu.%02lu MB\",\n", vsz_kb / 1024, ((vsz_kb % 1024) * 100) / 1024);
        }

        /* Formateo dinamico de Memoria Residente (RSS) en KB o MB */
        if (rss_kb < 1024) {
            seq_printf(m, "    \"rss\": \"%lu KB\",\n", rss_kb);
        } else {
            seq_printf(m, "    \"rss\": \"%lu.%02lu MB\",\n", rss_kb / 1024, ((rss_kb % 1024) * 100) / 1024);
        }

        /* Salida de metricas E/S en bytes */
        seq_printf(m, "    \"io_read_bytes\": %llu,\n", read_bytes);
        seq_printf(m, "    \"io_write_bytes\": %llu\n", write_bytes);

        seq_puts(m, "  }");
    }
    rcu_read_unlock();

    seq_puts(m, "\n]\n");
    return 0;
}

static int continfo_open(struct inode *inode, struct file *file) {
    return single_open(file, continfo_show, NULL);
}

static const struct proc_ops continfo_ops = {
    .proc_open    = continfo_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};

static int __init continfo_init(void) {
    if (!proc_create(PROC_NAME, 0444, NULL, &continfo_ops)) {
        pr_err("Error al crear la entrada en /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    pr_info("Modulo continfo cargado correctamente en /proc/%s\n", PROC_NAME);
    return 0;
}

static void __exit continfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    pr_info("Modulo continfo desinstalado correctamente\n");
}

module_init(continfo_init);
module_exit(continfo_exit);