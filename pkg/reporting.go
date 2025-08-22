package pkg

import (
	"fmt"
	"os"
	"strings"
)

func ReportPSS(f []Finding, check string) {
	var rep *os.File
	rep = os.Stdout
	fmt.Fprintf(rep, "Findings for the %s check\n", check)
	if f != nil {
		for _, i := range f {
			switch i.Check {
			case "hostpid", "hostnet", "hostipc", "privileged", "allowprivesc", "HostProcess", "Seccomp Disabled", "Unmasked Procmount", "Apparmor Disabled":
				if i.Container != "" {
					fmt.Fprintf(rep, "namespace %s : pod %s : container %s\n", i.Namespace, i.Pod, i.Container)
				} else {
					fmt.Fprintf(rep, "namespace %s : pod %s\n", i.Namespace, i.Pod)
				}
			case "Added Capabilities":
				fmt.Fprintf(rep, "namespace %s : pod %s : container %s added capabilities %s \n", i.Namespace, i.Pod, i.Container, strings.Join(i.Capabilities[:], ","))
			case "Dropped Capabilities":
				fmt.Fprintf(rep, "namespace %s : pod %s : container %s dropped capabilities %s \n", i.Namespace, i.Pod, i.Container, strings.Join(i.Capabilities[:], ","))
			case "Host Ports":
				fmt.Fprintf(rep, "namespace %s : pod %s : container %s : port %d\n", i.Namespace, i.Pod, i.Container, i.Hostport)
			case "Host Path":
				fmt.Fprintf(rep, "namespace %s : pod %s : volume %s : path %s\n", i.Namespace, i.Pod, i.Volume, i.Path)
			case "Unsafe Sysctl":
				fmt.Fprintf(rep, "namespace %s : pod %s : unsafe sysctl %s", i.Namespace, i.Pod, i.Sysctl)

			}
		}
	} else {
		fmt.Fprintln(rep, "No findings!")
	}
	fmt.Fprintln(rep, "")
}
