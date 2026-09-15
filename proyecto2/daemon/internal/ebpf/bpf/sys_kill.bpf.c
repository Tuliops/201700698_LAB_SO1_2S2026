#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct kill_event {
    u32 pid_sender;
    u32 pid_target;
    s32 sig;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
} events SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_kill")
int trace_sys_kill(struct trace_event_raw_sys_enter *ctx) {
    long pid_target = (long)ctx->args[0];
    long sig = (long)ctx->args[1];

    // Filtrar unicamente SIGKILL (9) y SIGTERM (15)
    if (sig != 9 && sig != 15) {
        return 0;
    }

    struct kill_event *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    u64 pid_tgid = bpf_get_current_pid_tgid();
    e->pid_sender = (u32)(pid_tgid >> 32);
    e->pid_target = (u32)pid_target;
    e->sig = (s32)sig;

    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";