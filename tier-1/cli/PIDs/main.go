package main

import (
	"fmt"
	"os/exec"
	"time"
)

func verifyProcessState(pid int, etapa string) {
	// 'ps -p <PID> -o stat,cmd' show process state (STAT)
	out, err := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "stat,cmd").Output()
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
	cmdRun := exec.Command("sleep", "1")

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
	cmdStart := exec.Command("sleep", "1")

	if err := cmdStart.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	pidZombie := cmdStart.Process.Pid
	fmt.Printf("[Start] Procces in background PID: %d\n", pidZombie)

	fmt.Println("[Start] Waiting 2 seconds to 'sleep' finished in SO...")
	time.Sleep(2 * time.Second) // extra time till dead

	fmt.Println("\n⚠️  PROCESS FINISHED, Wait() doesn't called:")
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
