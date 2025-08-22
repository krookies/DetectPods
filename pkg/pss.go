package pkg

import (
	"github.com/spf13/pflag"
	"strings"
)

type Finding struct {
	Check        string   //表示进行安全检查的标识或名称
	Namespace    string   //表示容器所在的命名空间
	Pod          string   //表示容器所属的 Pod 名称
	Container    string   `json:",omitempty"` //表示容器的名称
	Capabilities []string `json:",omitempty"` //表示容器的 Linux 容器权限（capabilities）列表
	Hostport     int      `json:",omitempty"` //表示容器使用的主机端口
	Volume       string   `json:",omitempty"` //表示容器挂载的卷
	Path         string   `json:",omitempty"` //表示容器中的路径
	Sysctl       string   `json:",omitempty"` //表示容器的 sysctl 设置
	Image        string   `json:",omitempty"` //表示容器所使用的镜像
}

func Hostpid(options *pflag.FlagSet) []Finding {
	var hostPidCont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		pod.GetObjectMeta()
		if pod.Spec.HostPID {
			p := Finding{Check: "hostpid", Namespace: pod.Namespace, Pod: pod.Name, Container: ""}
			hostPidCont = append(hostPidCont, p)
		}
	}
	return hostPidCont
}

func Hostnet(options *pflag.FlagSet) []Finding {
	var hostNetCont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {

		if pod.Spec.HostNetwork {
			p := Finding{Check: "hostnet", Namespace: pod.Namespace, Pod: pod.Name}
			hostNetCont = append(hostNetCont, p)
		}
	}
	return hostNetCont
}

func Hostipc(options *pflag.FlagSet) []Finding {
	var hostIpcCont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		if pod.Spec.HostIPC {
			p := Finding{Check: "hostipc", Namespace: pod.Namespace, Pod: pod.Name, Container: ""}
			hostIpcCont = append(hostIpcCont, p)
		}
	}
	return hostIpcCont
}

func HostPorts(options *pflag.FlagSet) []Finding {
	var hostPorts []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// 容器是否指定了端口
			cports := container.Ports != nil
			if cports {
				for _, port := range container.Ports {
					// 该端口是主机端口?
					if port.HostPort != 0 {
						p := Finding{Check: "Host Ports", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name, Hostport: int(port.HostPort)}
						hostPorts = append(hostPorts, p)
					}
				}
			}
		}
		for _, init_container := range pod.Spec.InitContainers {
			cports := init_container.Ports != nil
			if cports {
				for _, port := range init_container.Ports {
					if port.HostPort != 0 {
						p := Finding{Check: "Host Ports", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name, Hostport: int(port.HostPort)}
						hostPorts = append(hostPorts, p)
					}
				}
			}
		}
		for _, eph_container := range pod.Spec.EphemeralContainers {
			cports := eph_container.Ports != nil
			if cports {
				for _, port := range eph_container.Ports {
					if port.HostPort != 0 {
						p := Finding{Check: "Host Ports", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name, Hostport: int(port.HostPort)}
						hostPorts = append(hostPorts, p)
					}
				}
			}
		}
	}
	return hostPorts
}

func HostPath(options *pflag.FlagSet) []Finding {
	var hostPath []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		host_path := pod.Spec.Volumes != nil
		if host_path {
			for _, vol := range pod.Spec.Volumes {
				if vol.HostPath != nil {
					p := Finding{Check: "Host Path", Namespace: pod.Namespace, Pod: pod.Name, Volume: vol.Name, Path: vol.HostPath.Path}
					hostPath = append(hostPath, p)
				}
			}
		}
	}
	return hostPath
}

func HostProcess(options *pflag.FlagSet) []Finding {
	var hostprocesscont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		hostProcessPod := pod.Spec.SecurityContext.WindowsOptions != nil && *pod.Spec.SecurityContext.WindowsOptions.HostProcess
		if hostProcessPod {
			p := Finding{Check: "HostProcess", Namespace: pod.Namespace, Pod: pod.Name}
			hostprocesscont = append(hostprocesscont, p)
		}
		for _, container := range pod.Spec.Containers {
			hostProcessCont := container.SecurityContext != nil && container.SecurityContext.WindowsOptions != nil && *container.SecurityContext.WindowsOptions.HostProcess
			if hostProcessCont {
				p := Finding{Check: "HostProcess", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name}
				hostprocesscont = append(hostprocesscont, p)
			}
		}
		for _, init_container := range pod.Spec.InitContainers {
			hostProcessCont := init_container.SecurityContext != nil && init_container.SecurityContext.WindowsOptions != nil && *init_container.SecurityContext.WindowsOptions.HostProcess
			if hostProcessCont {
				p := Finding{Check: "HostProcess", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name}
				hostprocesscont = append(hostprocesscont, p)
			}
		}
		for _, eph_container := range pod.Spec.EphemeralContainers {
			hostProcessCont := eph_container.SecurityContext != nil && eph_container.SecurityContext.WindowsOptions != nil && *eph_container.SecurityContext.WindowsOptions.HostProcess
			if hostProcessCont {
				p := Finding{Check: "HostProcess", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name}
				hostprocesscont = append(hostprocesscont, p)
			}
		}
	}
	return hostprocesscont
}

func Privileged(options *pflag.FlagSet) []Finding {
	var privCont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			privileged_container := container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged
			if privileged_container {
				p := Finding{Check: "privileged", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name}
				privCont = append(privCont, p)
			}
		}
		for _, init_container := range pod.Spec.InitContainers {
			privileged_container := init_container.SecurityContext != nil && init_container.SecurityContext.Privileged != nil && *init_container.SecurityContext.Privileged
			if privileged_container {
				p := Finding{Check: "privileged", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name}
				privCont = append(privCont, p)
			}
		}
		for _, eph_container := range pod.Spec.EphemeralContainers {
			privileged_container := eph_container.SecurityContext != nil && eph_container.SecurityContext.Privileged != nil && *eph_container.SecurityContext.Privileged
			if privileged_container {
				p := Finding{Check: "privileged", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name}
				privCont = append(privCont, p)
			}
		}
	}
	return privCont
}

func AllowPrivEsc(options *pflag.FlagSet) []Finding {
	var allowPrivEscCont []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// 如果没有安全上下文，或者有安全上下文并且没有提到允许权限提升，则默认情况为true
			allowPrivilegeEscalation := (container.SecurityContext == nil) || (container.SecurityContext != nil && container.SecurityContext.AllowPrivilegeEscalation == nil)
			if allowPrivilegeEscalation {
				p := Finding{Check: "allowprivesc", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name}
				allowPrivEscCont = append(allowPrivEscCont, p)
			}
		}
		for _, init_container := range pod.Spec.InitContainers {
			allowPrivilegeEscalation := (init_container.SecurityContext == nil) || (init_container.SecurityContext != nil && init_container.SecurityContext.AllowPrivilegeEscalation == nil)
			if allowPrivilegeEscalation {
				p := Finding{Check: "allowprivesc", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name}
				allowPrivEscCont = append(allowPrivEscCont, p)
			}
		}
		for _, eph_container := range pod.Spec.EphemeralContainers {
			allowPrivilegeEscalation := (eph_container.SecurityContext == nil) || (eph_container.SecurityContext != nil && eph_container.SecurityContext.AllowPrivilegeEscalation == nil)
			if allowPrivilegeEscalation {
				p := Finding{Check: "allowprivesc", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name}
				allowPrivEscCont = append(allowPrivEscCont, p)
			}
		}
	}
	return allowPrivEscCont
}

func AddedCapabilities(options *pflag.FlagSet) []Finding {
	var capAdded []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			cap_added := container.SecurityContext != nil && container.SecurityContext.Capabilities != nil && container.SecurityContext.Capabilities.Add != nil
			if cap_added {
				var added_caps []string
				for _, cap := range container.SecurityContext.Capabilities.Add {
					added_caps = append(added_caps, string(cap))
				}
				p := Finding{Check: "Added Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name, Capabilities: added_caps}
				capAdded = append(capAdded, p)
			}
		}

		for _, init_container := range pod.Spec.InitContainers {
			cap_added := init_container.SecurityContext != nil && init_container.SecurityContext.Capabilities != nil && init_container.SecurityContext.Capabilities.Add != nil
			if cap_added {
				var added_caps []string
				for _, cap := range init_container.SecurityContext.Capabilities.Add {
					added_caps = append(added_caps, string(cap))
				}
				p := Finding{Check: "Added Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name, Capabilities: added_caps}
				capAdded = append(capAdded, p)
			}
		}

		for _, eph_container := range pod.Spec.EphemeralContainers {
			cap_added := eph_container.SecurityContext != nil && eph_container.SecurityContext.Capabilities != nil && eph_container.SecurityContext.Capabilities.Add != nil
			if cap_added {
				var added_caps []string
				for _, cap := range eph_container.SecurityContext.Capabilities.Add {
					added_caps = append(added_caps, string(cap))
				}
				p := Finding{Check: "Added Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name, Capabilities: added_caps}
				capAdded = append(capAdded, p)
			}
		}
	}
	return capAdded
}

func DroppedCapabilities(options *pflag.FlagSet) []Finding {
	var capDropped []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			cap_dropped := container.SecurityContext != nil && container.SecurityContext.Capabilities != nil && container.SecurityContext.Capabilities.Drop != nil
			if cap_dropped {
				var dropped_caps []string
				for _, cap := range container.SecurityContext.Capabilities.Drop {
					dropped_caps = append(dropped_caps, string(cap))
				}
				p := Finding{Check: "Dropped Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name, Capabilities: dropped_caps}
				capDropped = append(capDropped, p)
			}
		}

		for _, init_container := range pod.Spec.InitContainers {
			cap_dropped := init_container.SecurityContext != nil && init_container.SecurityContext.Capabilities != nil && init_container.SecurityContext.Capabilities.Drop != nil
			if cap_dropped {
				var dropped_caps []string
				for _, cap := range init_container.SecurityContext.Capabilities.Drop {
					dropped_caps = append(dropped_caps, string(cap))
				}
				p := Finding{Check: "Dropped Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name, Capabilities: dropped_caps}
				capDropped = append(capDropped, p)
			}
		}

		for _, eph_container := range pod.Spec.EphemeralContainers {
			cap_dropped := eph_container.SecurityContext != nil && eph_container.SecurityContext.Capabilities != nil && eph_container.SecurityContext.Capabilities.Drop != nil
			if cap_dropped {
				var dropped_caps []string
				for _, cap := range eph_container.SecurityContext.Capabilities.Drop {
					dropped_caps = append(dropped_caps, string(cap))
				}
				p := Finding{Check: "Dropped Capabilities", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name, Capabilities: dropped_caps}
				capDropped = append(capDropped, p)
			}
		}
	}
	return capDropped
}

func Seccomp(options *pflag.FlagSet) []Finding {
	var seccomp []Finding
	pods := ConnectWithPods(options)
	// 如果pod是无限制的,容器也是无限制的
	// 理论上,如果pod中的所有容器都是无限制的,我们可以在pod级别对其进行标记
	for _, pod := range pods.Items {
		unconfined_pod := (pod.Spec.SecurityContext == nil) || (pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SeccompProfile == nil) || (pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SeccompProfile != nil && pod.Spec.SecurityContext.SeccompProfile.Type == "Unconfined")
		if unconfined_pod {
			//log.Printf("Pod name %s was unconfined at pod level", pod.Name)
			for _, container := range pod.Spec.Containers {
				unconfined_container := (container.SecurityContext == nil) || (container.SecurityContext != nil && container.SecurityContext.SeccompProfile == nil) || (container.SecurityContext != nil && container.SecurityContext.SeccompProfile != nil && container.SecurityContext.SeccompProfile.Type == "Unconfined")
				if unconfined_container {
					p := Finding{Check: "Seccomp Disabled", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name}
					seccomp = append(seccomp, p)
				}
			}
			for _, init_container := range pod.Spec.InitContainers {
				unconfined_init_container := (init_container.SecurityContext == nil) || (init_container.SecurityContext != nil && init_container.SecurityContext.SeccompProfile == nil) || (init_container.SecurityContext != nil && init_container.SecurityContext.SeccompProfile != nil && init_container.SecurityContext.SeccompProfile.Type == "Unconfined")
				if unconfined_init_container {
					p := Finding{Check: "Seccomp Disabled", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name}
					seccomp = append(seccomp, p)
				}
			}
			for _, eph_container := range pod.Spec.EphemeralContainers {
				unconfined_eph_container := (eph_container.SecurityContext == nil) || (eph_container.SecurityContext != nil && eph_container.SecurityContext.SeccompProfile == nil) || (eph_container.SecurityContext != nil && eph_container.SecurityContext.SeccompProfile != nil && eph_container.SecurityContext.SeccompProfile.Type == "Unconfined")
				if unconfined_eph_container {
					p := Finding{Check: "Seccomp Disabled", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name}
					seccomp = append(seccomp, p)
				}
			}
		}
	}
	return seccomp
}

func Apparmor(options *pflag.FlagSet) []Finding {
	var apparmor []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		// 默认值应该是apparmor已设置,所以我们只关心它是否明确设置为unconfined
		if pod.Annotations != nil {
			for key, val := range pod.Annotations {
				if val == "unconfined" && strings.Split(key, "/")[0] == "container.apparmor.security.beta.kubernetes.io" {
					p := Finding{Check: "Apparmor Disabled", Namespace: pod.Namespace, Pod: pod.Name}
					apparmor = append(apparmor, p)
				}
			}
		}
	}
	return apparmor
}

func Procmount(options *pflag.FlagSet) []Finding {
	var unmaskedProc []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			unmask := container.SecurityContext != nil && container.SecurityContext.ProcMount != nil && *container.SecurityContext.ProcMount == "Unmasked"
			if unmask {
				p := Finding{Check: "Unmasked procmount", Namespace: pod.Namespace, Pod: pod.Name, Container: container.Name}
				unmaskedProc = append(unmaskedProc, p)
			}
		}
		for _, init_container := range pod.Spec.InitContainers {
			unmask := init_container.SecurityContext != nil && init_container.SecurityContext.ProcMount != nil && *init_container.SecurityContext.ProcMount == "Unmasked"
			if unmask {
				p := Finding{Check: "Unmasked procmount", Namespace: pod.Namespace, Pod: pod.Name, Container: init_container.Name}
				unmaskedProc = append(unmaskedProc, p)
			}
		}
		for _, eph_container := range pod.Spec.EphemeralContainers {
			unmask := eph_container.SecurityContext != nil && eph_container.SecurityContext.ProcMount != nil && *eph_container.SecurityContext.ProcMount == "Unmasked"
			if unmask {
				p := Finding{Check: "Unmasked procmount", Namespace: pod.Namespace, Pod: pod.Name, Container: eph_container.Name}
				unmaskedProc = append(unmaskedProc, p)
			}
		}
	}
	return unmaskedProc
}

func Sysctl(options *pflag.FlagSet) []Finding {
	var sysctls []Finding
	pods := ConnectWithPods(options)
	for _, pod := range pods.Items {
		sysctl := pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.Sysctls != nil
		if sysctl {
			for _, sys := range pod.Spec.SecurityContext.Sysctls {
				safe := []string{"kernel.shm_rmid_forced", "net.ipv4.ip_local_port_range", "net.ipv4.ip_unprivileged_port_start", "net.ipv4.tcp_syncookies", "net.ipv4.ping_group_range"}
				safe_sys := false
				for _, s := range safe {
					if sys.Name == s {
						safe_sys = true
					}
				}
				if !safe_sys {
					p := Finding{Check: "Unsafe Sysctl", Namespace: pod.Namespace, Sysctl: sys.Name}
					sysctls = append(sysctls, p)
				}
			}
		}
	}
	return sysctls
}
