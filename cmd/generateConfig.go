/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"getNoPSS/pkg"

	"github.com/spf13/cobra"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "generateConfig",
	Short: "生成示例配置文件",
	Long:  `生成一个示例的config.yaml配置文件，包含OpenAI API的配置模板`,
	Run: func(cmd *cobra.Command, args []string) {
		options := cmd.Flags()
		configPath, _ := options.GetString("output")

		if configPath == "" {
			configPath = "config.yaml"
		}

		err := pkg.SaveExampleConfig(configPath)
		if err != nil {
			fmt.Printf("❌ 生成配置文件失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 示例配置文件已生成: %s\n", configPath)
		fmt.Println("\n📝 请编辑配置文件，设置你的OpenAI API密钥和其他配置:")
		fmt.Printf("1. 打开文件: %s\n", configPath)
		fmt.Println("2. 设置 openai.api_key 为你的实际API密钥")
		fmt.Println("3. 根据需要修改 base_url 和 model")
		fmt.Println("\n🔍 之后你可以使用以下命令测试配置:")
		fmt.Printf("./getNoPSS testApi --config %s\n", configPath)
	},
}

func init() {
	rootCmd.AddCommand(generateConfigCmd)
	generateConfigCmd.Flags().StringP("output", "o", "config.yaml", "输出配置文件路径")
}
