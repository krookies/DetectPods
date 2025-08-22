/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"getNoPSS/pkg"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// aiAnalysisCmd represents the aiAnalysis command
var aiAnalysisCmd = &cobra.Command{
	Use:   "aiAnalysis",
	Short: "使用AI大模型分析Pod安全性",
	Long:  `使用OpenAI API对所有Pod进行深度安全分析，并将结果保存到本地文件`,
	Run: func(cmd *cobra.Command, args []string) {
		options := cmd.Flags()

		// 获取配置文件路径
		configPath, _ := options.GetString("config")

		// 加载配置文件
		config, err := pkg.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("❌ 加载配置文件失败: %v\n", err)
			fmt.Println("\n💡 你可以使用以下命令生成示例配置文件:")
			fmt.Println("./getNoPSS generateConfig")
			return
		}

		// 获取输出文件路径和格式
		outputFile, _ := options.GetString("output")
		outputFormat, _ := options.GetString("format")
		consoleOutput, _ := options.GetBool("console")

		if outputFile == "" {
			timestamp := time.Now().Format("20060102_150405")
			switch outputFormat {
			case "html":
				outputFile = fmt.Sprintf("pod_security_analysis_%s.html", timestamp)
			default:
				outputFile = fmt.Sprintf("pod_security_analysis_%s.json", timestamp)
			}
		}

		// 创建AI分析器
		analyzer := pkg.NewAIAnalyzer(config)

		// 获取过滤后的Pod列表
		pods := pkg.ConnectWithPods(options)

		if len(pods.Items) == 0 {
			fmt.Println("没有找到符合条件的Pod")
			return
		}

		fmt.Printf("开始AI安全分析，共 %d 个Pod...\n", len(pods.Items))
		fmt.Printf("使用模型: %s\n", config.OpenAI.Model)

		// 使用AI分析Pod
		analyses, err := analyzer.AnalyzePods(pods)
		if err != nil {
			fmt.Printf("AI分析失败: %v\n", err)
			return
		}

		// 如果启用控制台输出
		if consoleOutput {
			pkg.PrintAnalysisToConsole(analyses)
		}

		// 保存结果到文件
		switch outputFormat {
		case "html":
			err = pkg.SaveAnalysisResultsAsHTML(analyses, outputFile)
		default:
			err = saveAnalysisResults(analyses, outputFile)
		}

		if err != nil {
			fmt.Printf("保存分析结果失败: %v\n", err)
			return
		}

		// 打印统计信息
		printAnalysisStats(analyses)
		fmt.Printf("\n分析结果已保存到: %s\n", outputFile)
	},
}

func saveAnalysisResults(analyses []pkg.AIAnalysis, filename string) error {
	// 创建包含总结信息的完整报告
	report := struct {
		GeneratedAt time.Time        `json:"generated_at"`
		TotalPods   int              `json:"total_pods"`
		Summary     map[string]int   `json:"summary"`
		Analyses    []pkg.AIAnalysis `json:"analyses"`
	}{
		GeneratedAt: time.Now(),
		TotalPods:   len(analyses),
		Summary:     make(map[string]int),
		Analyses:    analyses,
	}

	// 计算统计信息
	for _, analysis := range analyses {
		report.Summary[analysis.SecurityLevel]++
	}

	// 将报告序列化为JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化分析结果失败: %w", err)
	}

	// 写入文件
	return os.WriteFile(filename, data, 0644)
}

func printAnalysisStats(analyses []pkg.AIAnalysis) {
	stats := make(map[string]int)
	totalIssues := 0

	for _, analysis := range analyses {
		stats[analysis.SecurityLevel]++
		totalIssues += len(analysis.Issues)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 AI安全分析统计")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("总Pod数量: %d\n", len(analyses))
	fmt.Printf("发现问题总数: %d\n", totalIssues)
	fmt.Println("\n安全等级分布:")

	levels := []struct{ name, symbol string }{
		{"SAFE", "✅"},
		{"MODERATE", "⚠️"},
		{"HIGH_RISK", "🔴"},
		{"CRITICAL", "🚨"},
		{"UNKNOWN", "❓"},
	}

	for _, level := range levels {
		if count, exists := stats[level.name]; exists && count > 0 {
			fmt.Printf("  %s %s: %d\n", level.symbol, level.name, count)
		}
	}
}

func init() {
	rootCmd.AddCommand(aiAnalysisCmd)

	// 添加参数
	aiAnalysisCmd.Flags().StringP("config", "", "config.yaml", "配置文件路径")
	aiAnalysisCmd.Flags().StringP("output", "o", "", "输出文件路径")
	aiAnalysisCmd.Flags().StringP("format", "f", "json", "输出格式 (json|html)")
	aiAnalysisCmd.Flags().BoolP("console", "c", false, "在控制台显示详细结果")
}
