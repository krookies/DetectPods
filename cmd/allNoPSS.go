/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"getNoPSS/pkg"
	"github.com/spf13/cobra"
)

// allNoPSSCmd represents the allNoPSS command
var allNoPSSCmd = &cobra.Command{
	Use:   "allNoPSS",
	Short: "获取所有不安全 Pod",
	Long:  `检索所有不符合安全标准的 Pod`,
	Run: func(cmd *cobra.Command, args []string) {
		options := cmd.Flags()
		hostPidCont := pkg.Hostpid(options) //检索使用主机 PID的 Pod
		pkg.ReportPSS(hostPidCont, "Host PID")

		hostNetCont := pkg.Hostnet(options) //检索使用主机网络的 Pod
		pkg.ReportPSS(hostNetCont, "Host Network")

		hostIpcCont := pkg.Hostipc(options) //检索使用主机 IPC 的 Pod
		pkg.ReportPSS(hostIpcCont, "Host IPC")

		hostPorts := pkg.HostPorts(options) //检索使用主机端口的容器
		pkg.ReportPSS(hostPorts, "Host Ports")

		hostPath := pkg.HostPath(options) //检索挂载主机路径卷的 Pod
		pkg.ReportPSS(hostPath, "Host Path")

		hostProcessCont := pkg.HostProcess(options) //检索使用 hostprocess 权限运行的容器
		pkg.ReportPSS(hostProcessCont, "Host Process")

		privCont := pkg.Privileged(options) //检索特权容器
		pkg.ReportPSS(privCont, "Privileged Container")

		allowPrivEscCont := pkg.AllowPrivEsc(options) //检索允许权限提升的容器
		pkg.ReportPSS(allowPrivEscCont, "Allow Privilege Escalation")

		capAdded := pkg.AddedCapabilities(options) //检索具有比默认配置更高权限的容器
		pkg.ReportPSS(capAdded, "Added Capabilities")

		capDropped := pkg.DroppedCapabilities(options) //检索通过降低权限来提高安全性的容器
		pkg.ReportPSS(capDropped, "Dropped Capabilities")

		seccomp := pkg.Seccomp(options) //检索未启用 Seccomp 的容器
		pkg.ReportPSS(seccomp, "Seccomp Disabled")

		apparmor := pkg.Apparmor(options) // 检索未启用 AppArmor 的 Pod
		pkg.ReportPSS(apparmor, "Apparmor Disabled")

		unmaskedProc := pkg.Procmount(options) //检索使用 "Unmasked" proc mount 的容器
		pkg.ReportPSS(unmaskedProc, "Unmasked Procmount")

		sysctls := pkg.Sysctl(options) //检索配置了不在安全列表中的 sysctl 参数的 Pod
		pkg.ReportPSS(sysctls, "Unsafe Sysctl")
	},
}

func init() {
	rootCmd.AddCommand(allNoPSSCmd)
}
