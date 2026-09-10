package kernel

import (
	"fmt"
	"os/exec"
	"strings"
)

// EnsureModuleLoaded ejecuta 'make load' forzando el contexto del directorio del modulo.
func EnsureModuleLoaded(modulePath string) error {
	fmt.Println("[+] Verificando/Cargando modulo de kernel...")

	cmd := exec.Command("make", "-C", modulePath, "load")
	cmd.Dir = modulePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error al cargar el modulo de kernel: %w\nSalida: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("[+] Modulo de kernel verificado/cargado con exito.")
	return nil
}

// UnloadModule desmonta el modulo del kernel y limpia los binarios compilados (make clean).
func UnloadModule(modulePath string) error {
	fmt.Println("[+] Desmontando modulo de kernel...")

	//  Ejecutar make unload
	cmdUnload := exec.Command("make", "-C", modulePath, "unload")
	cmdUnload.Dir = modulePath

	outputUnload, errUnload := cmdUnload.CombinedOutput()
	if errUnload != nil {
		fmt.Println("[!] 'make unload' fallo, intentando rmmod directo de respaldo...")
		fallbackCmd := exec.Command("sudo", "rmmod", "continfo")
		if fallbackOutput, fallbackErr := fallbackCmd.CombinedOutput(); fallbackErr != nil {
			fmt.Printf("[WARN] No se pudo desmontar el modulo con rmmod: %s\n", strings.TrimSpace(string(fallbackOutput)))
		} else {
			fmt.Println("[+] Modulo desmontado con comando rmmod directo.")
		}
	} else {
		fmt.Printf("[+] Modulo desmontado via Makefile:\n%s\n", strings.TrimSpace(string(outputUnload)))
	}

	//  Ejecutar make clean para eliminar los binarios .o, .ko, etc.
	fmt.Println("[+] Ejecutando 'make clean' para remover archivos temporales...")
	cmdClean := exec.Command("make", "-C", modulePath, "clean")
	cmdClean.Dir = modulePath

	outputClean, errClean := cmdClean.CombinedOutput()
	if errClean != nil {
		return fmt.Errorf("error al ejecutar 'make clean': %w\nSalida: %s", errClean, strings.TrimSpace(string(outputClean)))
	}

	fmt.Printf("[+] Limpieza de binarios completada con exito:\n%s\n", strings.TrimSpace(string(outputClean)))
	return nil
}