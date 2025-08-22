package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func SaveAnalysisResultsAsHTML(analyses []AIAnalysis, filename string) error {
	htmlContent := generateHTMLReport(analyses)
	return os.WriteFile(filename, []byte(htmlContent), 0644)
}

func generateHTMLReport(analyses []AIAnalysis) string {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kubernetes Pod 安全分析报告</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 20px; line-height: 1.6; }
        .header { background: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .summary { display: flex; gap: 20px; margin-bottom: 20px; }
        .stat-card { background: white; border: 1px solid #dee2e6; border-radius: 8px; padding: 15px; flex: 1; text-align: center; }
        .safe { border-left: 4px solid #28a745; }
        .moderate { border-left: 4px solid #ffc107; }
        .high-risk { border-left: 4px solid #fd7e14; }
        .critical { border-left: 4px solid #dc3545; }
        .unknown { border-left: 4px solid #6c757d; }
        .pod-card { background: white; border: 1px solid #dee2e6; border-radius: 8px; margin: 10px 0; padding: 20px; }
        .pod-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
        .security-level { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; text-transform: uppercase; }
        .level-safe { background: #d4edda; color: #155724; }
        .level-moderate { background: #fff3cd; color: #856404; }
        .level-high-risk { background: #f8d7da; color: #721c24; }
        .level-critical { background: #f5c6cb; color: #721c24; }
        .level-unknown { background: #e2e3e5; color: #383d41; }
        .issues, .recommendations { margin: 10px 0; }
        .issues ul, .recommendations ul { padding-left: 20px; }
        .issues li { color: #dc3545; margin: 5px 0; }
        .recommendations li { color: #28a745; margin: 5px 0; }
        h1, h2, h3 { color: #333; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔒 Kubernetes Pod 安全分析报告</h1>
        <p>生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
        <p>分析Pod数量: ` + fmt.Sprintf("%d", len(analyses)) + `</p>
    </div>`

	// 添加统计摘要
	stats := make(map[string]int)
	for _, analysis := range analyses {
		stats[analysis.SecurityLevel]++
	}

	html += `<div class="summary">`

	levels := []struct{ name, class, label string }{
		{"SAFE", "safe", "安全"},
		{"MODERATE", "moderate", "中等风险"},
		{"HIGH_RISK", "high-risk", "高风险"},
		{"CRITICAL", "critical", "严重"},
		{"UNKNOWN", "unknown", "未知"},
	}

	for _, level := range levels {
		count := stats[level.name]
		html += fmt.Sprintf(`
        <div class="stat-card %s">
            <h3>%s</h3>
            <h2>%d</h2>
        </div>`, level.class, level.label, count)
	}

	html += `</div>`

	// 添加每个Pod的详细分析
	html += `<h2>详细分析结果</h2>`

	for _, analysis := range analyses {
		levelClass := strings.ToLower(strings.ReplaceAll(analysis.SecurityLevel, "_", "-"))
		html += fmt.Sprintf(`
    <div class="pod-card">
        <div class="pod-header">
            <div>
                <h3>%s / %s</h3>
            </div>
            <span class="security-level level-%s">%s</span>
        </div>`, analysis.Namespace, analysis.Pod, levelClass, analysis.SecurityLevel)

		if len(analysis.Issues) > 0 {
			html += `<div class="issues"><h4>🚨 发现的问题:</h4><ul>`
			for _, issue := range analysis.Issues {
				html += fmt.Sprintf("<li>%s</li>", issue)
			}
			html += `</ul></div>`
		}

		if len(analysis.Recommendations) > 0 {
			html += `<div class="recommendations"><h4>💡 安全建议:</h4><ul>`
			for _, rec := range analysis.Recommendations {
				html += fmt.Sprintf("<li>%s</li>", rec)
			}
			html += `</ul></div>`
		}

		html += `</div>`
	}

	html += `</body></html>`
	return html
}

func PrintAnalysisToConsole(analyses []AIAnalysis) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔒 AI Pod安全分析结果")
	fmt.Println(strings.Repeat("=", 80))

	for i, analysis := range analyses {
		fmt.Printf("\n[%d/%d] Pod: %s/%s\n", i+1, len(analyses), analysis.Namespace, analysis.Pod)

		// 安全等级用不同符号表示
		var levelSymbol string
		switch analysis.SecurityLevel {
		case "SAFE":
			levelSymbol = "✅"
		case "MODERATE":
			levelSymbol = "⚠️"
		case "HIGH_RISK":
			levelSymbol = "🔴"
		case "CRITICAL":
			levelSymbol = "🚨"
		default:
			levelSymbol = "❓"
		}

		fmt.Printf("安全等级: %s %s\n", levelSymbol, analysis.SecurityLevel)

		if len(analysis.Issues) > 0 {
			fmt.Println("\n🚨 发现的安全问题:")
			for j, issue := range analysis.Issues {
				fmt.Printf("  %d. %s\n", j+1, issue)
			}
		}

		if len(analysis.Recommendations) > 0 {
			fmt.Println("\n💡 安全改进建议:")
			for j, rec := range analysis.Recommendations {
				fmt.Printf("  %d. %s\n", j+1, rec)
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}
