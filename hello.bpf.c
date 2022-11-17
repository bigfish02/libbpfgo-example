// +build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>


struct {
        __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
        __uint(key_size, sizeof(u32));
        __uint(value_size, sizeof(u32));
}
events SEC(".maps");

// Example: tracing a message on a kprobe
SEC("tracepoint/syscalls/sys_enter_execve")
int hello(void *ctx)
{
    bpf_printk("I'm alive!");
    return 0;
}

// Example of passing data using a perf map
// Similar to bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count();}'
SEC("tracepoint/syscalls/sys_exit_execve")
int hello_bpftrace(void *ctx)
{
    char data[100];
    bpf_get_current_comm(&data, 100);
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &data, 100);
    return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";