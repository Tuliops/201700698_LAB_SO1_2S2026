package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel bpf bpf/sys_kill.bpf.c -- -I./bpf
