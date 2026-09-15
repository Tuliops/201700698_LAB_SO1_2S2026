#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/mm.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tulio Pirir");
MODULE_DESCRIPTION("Sonda Kernel SO1 - USAC 201700698");

#define PROC_NAME "continfo_pr2_so1_201700698"

static int continfo_show(struct seq_file *m, void *v) {

    struct task_struct *task;
    struct sysinfo i;
    bool first = true;

    // Obtener informacion global de memoria RAM del sistema
    si_meminfo(&i);
    unsigned long total_ram_mb = (i.totalram * i.mem_unit) / (1024 * 1024);
    unsigned long free_ram_mb = (i.freeram * i.mem_unit) / (1024 * 1024);
    unsigned long used_ram_mb = total_ram_mb - free_ram_mb;

    seq_printf(m, "{\n");
    seq_printf(m, "  \"total_ram_mb\": %lu,\n", total_ram_mb);
    seq_printf(m, "  \"free_ram_mb\": %lu,\n", free_ram_mb);
    seq_printf(m, "  \"used_ram_mb\": %lu,\n", used_ram_mb);
    seq_printf(m, "  \"processes\": [\n");

    rcu_read_lock();
    for_each_process(task) {
        if (!first) {
            seq_printf(m, ",\n");
        }
        first = false;

        unsigned long vsz = 0;
        unsigned long rss = 0;

        if (task->mm) {
            vsz = task->mm->total_vm << (PAGE_SHIFT - 10); // KB
            rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10); // KB
        }

        seq_printf(m, "    {\n");
        seq_printf(m, "      \"pid\": %d,\n", task->pid);
        seq_printf(m, "      \"name\": \"%s\",\n", task->comm);
        seq_printf(m, "      \"cpu_jiffies\": %llu,\n", (u64)(task->utime + task->stime));        seq_printf(m, "      \"vsz_kb\": %lu,\n", vsz);
        seq_printf(m, "      \"rss_kb\": %lu\n", rss);
        seq_printf(m, "    }");
    }
    rcu_read_unlock();

    seq_printf(m, "\n  ]\n}\n");
    return 0;
}

static int continfo_open(struct inode *inode, struct file *file) {
    return single_open(file, continfo_show, NULL);
}

static const struct proc_ops continfo_fops = {
    .proc_open    = continfo_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};

static int __init continfo_init(void) {
    proc_create(PROC_NAME, 0444, NULL, &continfo_fops);
    printk(KERN_INFO "Modulo %s cargado.\n", PROC_NAME);
    return 0;
}

static void __exit continfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "Modulo %s removido.\n", PROC_NAME);
}

module_init(continfo_init);
module_exit(continfo_exit); 