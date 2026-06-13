package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const tfDevDir = "infra/environments/dev"

var rootCmd = &cobra.Command{
	Use:   "voltctl",
	Short: "VoltStream developer environment CLI",
	Long:  "voltctl provisions and manages VoltStream IoT cloud environments for local development and hardware testing.",
}

// --- env commands ---

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage cloud environments",
}

var envUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Provision the dev environment via Terraform apply",
	RunE:  runEnvUp,
}

var envDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Destroy the dev environment via Terraform destroy",
	RunE:  runEnvDown,
}

// --- sim commands ---

var simCmd = &cobra.Command{
	Use:   "sim",
	Short: "Control the battery edge simulator",
}

var simStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the Go battery simulator",
	RunE:  runSimStart,
}

// --- flags ---

var (
	dryRun   bool
	machines int
)

func init() {
	envUpCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing")
	envDownCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing")
	simStartCmd.Flags().IntVar(&machines, "machines", 500, "number of machines to simulate")
	simStartCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing")

	envCmd.AddCommand(envUpCmd, envDownCmd)
	simCmd.AddCommand(simStartCmd)
	rootCmd.AddCommand(envCmd, simCmd)
}

func runEnvUp(cmd *cobra.Command, args []string) error {
	return runTerraform("apply", "-auto-approve")
}

func runEnvDown(cmd *cobra.Command, args []string) error {
	return runTerraform("destroy", "-auto-approve")
}

func runTerraform(subArgs ...string) error {
	args := append([]string{"-chdir=" + tfDevDir}, subArgs...)
	if dryRun {
		fmt.Printf("[dry-run] terraform %v\n", args)
		return nil
	}
	return run("terraform", args...)
}

func runSimStart(cmd *cobra.Command, args []string) error {
	simArgs := []string{
		"run", "./simulator/cmd/main.go",
		fmt.Sprintf("--machines=%d", machines),
		"--interval=2s",
	}
	if dryRun {
		fmt.Printf("[dry-run] go %v\n", simArgs)
		return nil
	}
	return run("go", simArgs...)
}

func run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
