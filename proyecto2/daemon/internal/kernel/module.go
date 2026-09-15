package kernel

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// CheckAndLoadModule verifica e inserta el modulo de kernel si no esta cargado
func CheckAndLoadModule() error {
	out, err := exec.Command("lsmod").Output()
	if err != nil {
		return fmt.Errorf("error ejecutando lsmod: %w", err)
	}

	if strings.Contains(string(out), "continfo_pr2_so1_201700698") {
		log.Println("[KERNEL] Modulo continfo_pr2_so1_201700698 ya esta cargado.")
		return nil
	}

	log.Println("[KERNEL] Cargando modulo continfo_pr2_so1_201700698...")
	cmd := exec.Command("sudo", "insmod", "../kernel_module/continfo_pr2_so1_201700698.ko")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error insertando modulo de kernel: %w", err)
	}

	log.Println("[KERNEL] Modulo cargado exitosamente.")
	return nil
}
