package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval int   `yaml:"interval"`
	Targets []string `yaml:"targets"`
}

func loadConfig(filename string) (*Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter caminho do executável: %w", err)
	}
	
	execDir := filepath.Dir(execPath)

	yamlPath := filepath.Join(execDir, filename)

	file, err := os.Open(yamlPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ping(address string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS != "windows" {
		cmd = exec.Command("ping", "-c", "1", address)
	} else {
		cmd = exec.Command("ping", "-n", "1", address)
	}

	err := cmd.Run()
	return err == nil
}

func isValidAddress(address string) bool {
	_, err := net.LookupHost(address)
	return err == nil
}

func main () {
	config, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo de configurações:", err)
		return
	}

	if len(config.Targets) == 0 {
		fmt.Println("Nenhum endereço configurado no arquivo de configurações.")
		return
	}

	for {
		fmt.Println("Destino                           | Status")
		fmt.Println("------------------------------------------")

		for _, target := range config.Targets {
			if !isValidAddress(target) {
				color.New(color.FgRed).Printf("%-33s | Erro (endereço inválido)\n", target)
				continue
			}

			if ping(target) {
				color.New(color.FgGreen).Printf("%-33s | OK\n", target)
			} else {
				color.New(color.FgRed).Printf("%-33s | Erro\n", target)
			}
		}

		time.Sleep(time.Duration(config.Interval) * time.Second)
		clearScreen()
	}

}

func clearScreen() {
	if runtime.GOOS != "windows" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
