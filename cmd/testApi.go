/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"getNoPSS/pkg"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// testApiCmd represents the testApi command
var testApiCmd = &cobra.Command{
	Use:   "testApi",
	Short: "测试OpenAI API连接",
	Long:  `测试OpenAI API连接是否正常，用于调试配置问题`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("测试OpenAI API连接...")

		// 获取配置文件路径
		options := cmd.Flags()
		configPath, _ := options.GetString("config")

		// 加载配置文件
		config, err := pkg.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("❌ 加载配置文件失败: %v\n", err)
			fmt.Println("\n💡 你可以使用以下命令生成示例配置文件:")
			fmt.Println("./getNoPSS generateConfig")
			return
		}

		// 创建AI分析器
		analyzer := pkg.NewAIAnalyzer(config)

		// 进行一个简单的API测试
		resp, err := analyzer.GetClient().CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: "gpt-3.5-turbo", // 使用便宜的模型进行测试
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: "Hello, this is a test. Please respond with 'API connection successful'.",
					},
				},
				MaxTokens: 20,
			},
		)

		if err != nil {
			fmt.Printf("❌ API连接失败: %v\n", err)
			fmt.Println("\n可能的解决方案:")
			fmt.Println("1. 检查配置文件中的 api_key 是否正确")
			fmt.Println("2. 检查配置文件中的 base_url 是否正确")
			fmt.Println("3. 确认网络连接正常")
			fmt.Println("4. 验证API服务是否可用")
			return
		}

		if len(resp.Choices) == 0 {
			fmt.Println("❌ API返回了空响应")
			return
		}

		fmt.Println("✅ API连接成功!")
		fmt.Printf("响应: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("使用的模型: %s\n", resp.Model)
		fmt.Printf("Token使用情况: %+v\n", resp.Usage)
		fmt.Printf("配置的Base URL: %s\n", config.OpenAI.BaseURL)
		fmt.Printf("配置的模型: %s\n", config.OpenAI.Model)
	},
}

func init() {
	rootCmd.AddCommand(testApiCmd)
	testApiCmd.Flags().StringP("config", "", "config.yaml", "配置文件路径")
}
