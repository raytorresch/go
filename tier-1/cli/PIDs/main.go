package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// sleepCmd returns a blocking command equivalent to `sleep <seconds>`,
// since Windows has no "sleep" binary in %PATH% by default.
func sleepCmd(seconds string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Start-Sleep -Seconds %s", seconds))
	}
	return exec.Command("sleep", seconds)
}

func verifyProcessState(pid int, etapa string) {
	// Windows has no zombie/defunct state visible in the process list:
	// once the process exits, tasklist simply stops listing it,
	// regardless of whether Wait() was called. This check only shows
	// the real zombie-vs-clean difference on macOS/Linux via `ps`.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
	} else {
		cmd = exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "stat,cmd")
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("[%s] PID %d NON-EXISTENT Complete cleaning.\n", etapa, pid)
		return
	}

	fmt.Printf("[%s] Process in table PID %d:\n%s", etapa, pid, string(out))
}

func main() {

	// -------------------------------------------------------------
	// CASE 1: cmd.Run()
	// -------------------------------------------------------------
	fmt.Println("\n--- 1.cmd.Run() ---")
	cmdRun := sleepCmd("2")

	fmt.Println("[Run] Ejecutando comando bloqueante (espera 1s)...")
	start := time.Now()

	// Run() run process and call Wait()
	if err := cmdRun.Run(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("[Run] Finished %v.\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("[Run] Finished process PID :", cmdRun.Process.Pid)

	// Verify PID existence in OS
	verifyProcessState(cmdRun.Process.Pid, "Post-Run")

	// -------------------------------------------------------------
	// CASE 2: cmd.Start() WITHOUT cmd.Wait() (Make ZOMBIE)
	// -------------------------------------------------------------
	fmt.Println("\n--- 2. cmd.Start() without Wait() ---")
	cmdStart := sleepCmd("4")

	if err := cmdStart.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	pidZombie := cmdStart.Process.Pid
	fmt.Printf("[Start] Procces in background PID: %d\n", pidZombie)

	fmt.Println("[Start] Waiting 2 seconds to 'sleep' finished in SO...")
	time.Sleep(2 * time.Second) // extra time till dead

	fmt.Println("\nPROCESS FINISHED, Wait() doesn't called:")
	verifyProcessState(pidZombie, "Without Wait()")

	// -------------------------------------------------------------
	// CASE 3: Cleaning  Zombies cmd.Wait()
	// -------------------------------------------------------------
	fmt.Println("\n--- 3. Calling cmd.Wait() to clean Zombies ---")
	if err := cmdStart.Wait(); err != nil {
		fmt.Println("Wait error:", err)
	}

	fmt.Println("Wait() runed. Verifing process table:")
	verifyProcessState(pidZombie, "Post-Wait()")
}
