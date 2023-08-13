package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
    "log"
    "bufio"
    "strings"
)

// go run main.go run <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
        cleanup()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func add_path_to_env(nix_store_path string, container_root_path string) {
    local_nix_store_path := filepath.Join(container_root_path, nix_store_path)
    fileInfo, err := os.Stat(nix_store_path)
    if err != nil {
        fmt.Printf("argh")
    }
    if fileInfo.IsDir() {
        _, err:= os.Stat(local_nix_store_path)
        if os.IsNotExist(err) {
            must(syscall.Mkdir(local_nix_store_path, 700))
        }
	    must(syscall.Mount(nix_store_path, local_nix_store_path, "", syscall.MS_BIND, ""))
    } else {
        _, err:= os.Stat(local_nix_store_path)
        if os.IsNotExist(err) {
	        // must(syscall.Link(nix_store_path, local_nix_store_path))
            cp_cmd := exec.Command("cp", nix_store_path, local_nix_store_path)
            cp_cmd.Run()
        }
    }
}

func remove_path_from_env(nix_store_path string) {
    // fmt.Printf("Trying to remove: %s\n", nix_store_path)
    fileInfo, err := os.Stat(nix_store_path)
    if err != nil {
        fmt.Printf("argh")
    }
    if fileInfo.IsDir() {
        // fmt.Printf("it is a dir\n")
        must(syscall.Unmount(nix_store_path, 0))
        must(syscall.Rmdir(nix_store_path))
    } else {
        // fmt.Printf("it is a file\n")
        rm_cmd := exec.Command("rm", nix_store_path)
        rm_cmd.Run()
    }
}

func set_env(deps_filename string, container_root_path string) {
    f, err := os.Open(deps_filename)
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        nix_store_path := scanner.Text()
        // fmt.Printf("Adding: %s\n", nix_store_path)
        add_path_to_env(nix_store_path, container_root_path)
        // fmt.Printf("Added: %s\n", nix_store_path)
    }
    f.Close()
}

func clean_env(deps_filename string) {
    f, err := os.Open(deps_filename)
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        nix_store_path := scanner.Text()
        // fmt.Printf("Removing: %s\n", nix_store_path)
        remove_path_from_env(nix_store_path)
        // fmt.Printf("Removed: %s\n", nix_store_path)
    }
    f.Close()
}

func setup_nix_env(env_filename string) {
    f, err := os.Open(env_filename)
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(f)
    skip := 5
    for scanner.Scan() {
        command := scanner.Text()
        if skip == 0 {
            s := command[11:]
            splitted := strings.Split(s, "=")
            var_name := splitted[0]
            var_value := strings.Join(splitted[1:], "=")
            must(syscall.Setenv(var_name, var_value))
        } else {
            skip = skip - 1
        }
    }
    f.Close()
}

func cleanup() {
    ici, err := os.Getwd()
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
    }
    fmt.Printf(" [*] Starting the cleanup\n")
    container_root_path := filepath.Join(ici, "nixos_root_empty")
    fmt.Printf(" [*] Rm /\n")
    must(os.RemoveAll(container_root_path))
}

func child() {
	fmt.Printf("Running %v as %d \n", os.Args[2:], os.Getpid())

	cg()


    ici, err := os.Getwd()
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
    }
    container_root_path := filepath.Join(ici, "nixos_root_empty")
    must(syscall.Mkdir(container_root_path, 700))
    must(syscall.Mkdir(filepath.Join(container_root_path, "etc/"), 700))
    must(syscall.Mkdir(filepath.Join(container_root_path, "nix/"), 700))
    must(syscall.Mkdir(filepath.Join(container_root_path, "nix/store"), 700))
    must(syscall.Mkdir(filepath.Join(container_root_path, "root"), 700))

    passwd_content := fmt.Sprintf("root:x:0:0:sysadmin:/root:%s/bin/bash", os.Args[3])
     if err := os.WriteFile(filepath.Join(container_root_path, "etc/passwd"), []byte(passwd_content), 0644); err != nil {
        log.Fatal(err)
    }

	must(syscall.Mount(ici, filepath.Join(container_root_path, "root"), "", syscall.MS_BIND, ""))

    set_env(os.Args[2], container_root_path)

	must(syscall.Sethostname([]byte("nix-shell-container")))
	fmt.Printf(" [*] Set container hostname\n")
	must(syscall.Chroot(container_root_path))
	fmt.Printf(" [*] Perform chroot\n")
	must(os.Chdir("/root"))
	fmt.Printf(" [*] Change dir\n")

	// must(syscall.Mount("/proc", "/proc", "proc", 0, ""))

    setup_nix_env(os.Args[4])

    if len(os.Args) > 5 {
        cmd3 := exec.Command("bash", "-c", os.Args[5])
        cmd3.Stdin = os.Stdin
        cmd3.Stdout = os.Stdout
        cmd3.Stderr = os.Stderr
        cmd3.Run()
    } else {
        cmd3 := exec.Command("bash")//, os.Args[2:]...)
        cmd3.Stdin = os.Stdin
        cmd3.Stdout = os.Stdout
        cmd3.Stderr = os.Stderr
        cmd3.Run()
    }

	// must(syscall.Unmount("/proc", 0))
	// fmt.Printf(" [*] Unmount /proc\n")


    clean_env(os.Args[2])
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "liz"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "liz/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
