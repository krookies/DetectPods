package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	corev1 "k8s.io/api/core/v1"
)

type AIAnalysis struct {
	Namespace       string    `json:"namespace"`
	Pod             string    `json:"pod"`
	SecurityLevel   string    `json:"security_level"` // "SAFE", "MODERATE", "HIGH_RISK", "CRITICAL"
	Issues          []string  `json:"issues"`
	Recommendations []string  `json:"recommendations"`
	Timestamp       time.Time `json:"timestamp"`
}

type AIAnalyzer struct {
	client *openai.Client
	model  string
}

func NewAIAnalyzer(config *Config) *AIAnalyzer {
	if config == nil {
		log.Fatal().Msg("Config is required for AI analyzer")
	}

	if config.OpenAI.APIKey == "" {
		log.Fatal().Msg("OpenAI API key is required in config")
	}

	// 创建配置
	openaiConfig := openai.DefaultConfig(config.OpenAI.APIKey)

	// 检查是否设置了自定义base URL
	if config.OpenAI.BaseURL != "" {
		// 处理base URL，确保没有重复的路径
		baseURL := config.OpenAI.BaseURL
		// 如果base URL以 /v1/ 或 /v1 结尾，则移除它，因为OpenAI客户端会自动添加
		if strings.HasSuffix(baseURL, "/v1/") {
			baseURL = strings.TrimSuffix(baseURL, "/v1/")
		} else if strings.HasSuffix(baseURL, "/v1") {
			baseURL = strings.TrimSuffix(baseURL, "/v1")
		}

		openaiConfig.BaseURL = baseURL
		log.Info().Msgf("使用配置文件中的OpenAI Base URL: %s (处理后: %s)", config.OpenAI.BaseURL, baseURL)
	}

	return &AIAnalyzer{
		client: openai.NewClientWithConfig(openaiConfig),
		model:  config.OpenAI.Model,
	}
}

// GetClient 返回OpenAI客户端，用于测试连接
func (ai *AIAnalyzer) GetClient() *openai.Client {
	return ai.client
}

// cleanResponseContent 清理AI响应，移除Markdown代码块标记
func cleanResponseContent(content string) string {
	// 移除开头的 ```json 或 ```
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}

	// 移除结尾的 ```
	if strings.HasSuffix(content, "```") {
		content = strings.TrimSuffix(content, "```")
	}

	return strings.TrimSpace(content)
}

func (ai *AIAnalyzer) AnalyzePod(pod *corev1.Pod) (*AIAnalysis, error) {
	// 将Pod对象转换为JSON格式
	podJSON, err := json.MarshalIndent(pod, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pod to JSON: %w", err)
	}

	prompt := fmt.Sprintf(`作为Kubernetes安全专家，请分析以下Pod的安全配置。请重点关注以下安全问题：

1. 特权容器 (privileged containers)
2. hostNetwork, hostPID, hostIPC的使用
3. 不安全的卷挂载 (hostPath volumes)
4. 权限提升 (allowPrivilegeEscalation)
5. Linux capabilities的添加和删除
6. 安全上下文配置
7. 资源限制和请求
8. 镜像安全 (latest标签, 非官方镜像等)
9. seccomp和AppArmor配置
10. 网络策略和端口暴露

Pod配置：
%s

请以JSON格式返回分析结果，包含以下字段：
- security_level: "SAFE", "MODERATE", "HIGH_RISK", "CRITICAL"
- issues: 发现的安全问题列表
- recommendations: 安全改进建议列表

只返回JSON，不要包含其他文本。`, string(podJSON))

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: ai.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // 降低随机性，提高一致性
			MaxTokens:   2000,
		},
	)

	if err != nil {
		// 记录更详细的错误信息
		log.Error().Err(err).
			Str("pod", pod.Name).
			Str("namespace", pod.Namespace).
			Str("model", ai.model).
			Msg("OpenAI API调用失败")
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI API")
	}

	responseContent := resp.Choices[0].Message.Content

	// 清理响应内容，移除可能的Markdown代码块标记
	cleanedContent := cleanResponseContent(responseContent)

	// 解析AI返回的JSON响应
	var aiResult struct {
		SecurityLevel   string   `json:"security_level"`
		Issues          []string `json:"issues"`
		Recommendations []string `json:"recommendations"`
	}

	err = json.Unmarshal([]byte(cleanedContent), &aiResult)
	if err != nil {
		log.Warn().
			Str("response", responseContent).
			Str("pod", pod.Name).
			Str("namespace", pod.Namespace).
			Err(err).
			Msg("Failed to parse AI response as JSON, creating fallback analysis")

		// 如果解析失败，创建一个基础分析结果
		return &AIAnalysis{
			Namespace:       pod.Namespace,
			Pod:             pod.Name,
			SecurityLevel:   "UNKNOWN",
			Issues:          []string{fmt.Sprintf("AI analysis failed to parse response: %v", err)},
			Recommendations: []string{"Manual security review recommended", "Check API response format"},
			Timestamp:       time.Now(),
		}, nil
	}

	return &AIAnalysis{
		Namespace:       pod.Namespace,
		Pod:             pod.Name,
		SecurityLevel:   aiResult.SecurityLevel,
		Issues:          aiResult.Issues,
		Recommendations: aiResult.Recommendations,
		Timestamp:       time.Now(),
	}, nil
}

func (ai *AIAnalyzer) AnalyzePods(pods *corev1.PodList) ([]AIAnalysis, error) {
	var analyses []AIAnalysis

	log.Info().Msgf("开始AI分析 %d 个Pods", len(pods.Items))

	for i, pod := range pods.Items {
		log.Info().Msgf("分析Pod %d/%d: %s/%s", i+1, len(pods.Items), pod.Namespace, pod.Name)

		analysis, err := ai.AnalyzePod(&pod)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to analyze pod %s/%s", pod.Namespace, pod.Name)
			// 继续处理其他Pod，不因为一个失败而停止
			continue
		}

		analyses = append(analyses, *analysis)

		// 添加延迟以避免API限流
		time.Sleep(500 * time.Millisecond)
	}

	log.Info().Msgf("AI分析完成，共分析了 %d 个Pods", len(analyses))
	return analyses, nil
}
