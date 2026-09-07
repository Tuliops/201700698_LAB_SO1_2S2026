savedcmd_continfo.ko := ld -r -m elf_x86_64 -z noexecstack --no-warn-rwx-segments --build-id=sha1  -T /home/tj/wsl-kernel/scripts/module.lds -o continfo.ko continfo.o continfo.mod.o .module-common.o
