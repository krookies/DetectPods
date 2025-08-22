package pkg

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func initKubeClient() (*kubernetes.Clientset, error) {
	//初始化 Kubernetes 客户端,使用 clientcmd 包加载 kubeconfig 文件并创建一个客户端配置
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Error().Err(err).Msg("initKubeClient: failed creating ClientConfig")
		return nil, err
	}
	// 使用创建的客户端配置创建 Kubernetes 客户端集（Clientset）
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("initKubeClient: failed creating Clientset")
		return nil, err
	}
	// 返回初始化后的 Kubernetes 客户端集
	return clientset, nil
}

func ConnectWithPods(options *pflag.FlagSet) *corev1.PodList {
	// 调用 initKubeClient() 初始化 Kubernetes 客户
	clientset, err := initKubeClient()
	if err != nil {
		log.Print(err)
	}
	// 使用初始化的客户端获取所有 Pod 的列表
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Print(err)
	}
	// 根据传递的命令行选项中的 exclude 参数，过滤不需要的 Pod
	exclude, err := options.GetString("exclude")

	var excludeList []string
	if exclude != "" { //如果 exclude 参数不为空，则排除指定命名空间的 Pod
		excludeList = strings.Split(exclude, ",")
	}
	if err != nil {
		log.Print(err)
	}
	//返回经过过滤的 Pod 列表
	filteredPods := &corev1.PodList{}
	for _, pod := range pods.Items {
		if len(excludeList) > 0 {
			excluded := false
			for _, s := range excludeList {
				if strings.Contains(pod.Namespace, s) {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}
		filteredPods.Items = append(filteredPods.Items, pod)
	}

	return filteredPods
}
