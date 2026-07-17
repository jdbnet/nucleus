package mailer

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"net/smtp"
	"os"
	"sort"

	"nucleus/internal/models"
)

func SendReport(target models.Target, findings []models.Finding) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if host == "" || port == "" {
		log.Println("SMTP_HOST or SMTP_PORT not set, skipping email report.")
		return
	}

	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	to := os.Getenv("SMTP_TO")

	if from == "" || to == "" {
		log.Println("SMTP_FROM or SMTP_TO not set, skipping email report.")
		return
	}

	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	subject := fmt.Sprintf("Nucleus Scan Report: %s", target.Name)

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("To: %s\r\n", to))
	body.WriteString(fmt.Sprintf("From: %s\r\n", from))
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	body.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")

	body.WriteString("<html><body style='font-family: sans-serif;'>")
	body.WriteString(fmt.Sprintf("<h2>Scan Report for %s (%s)</h2>", target.Name, target.Address))
	body.WriteString("<p>The scheduled Nuclei scan has completed.</p>")
	
	if len(findings) == 0 {
		body.WriteString("<p>No findings were detected.</p>")
	} else {
		// Calculate summary
		counts := map[string]int{"critical": 0, "high": 0, "medium": 0, "low": 0, "info": 0}
		for _, f := range findings {
			counts[f.Severity]++
		}

		body.WriteString("<h3>Summary</h3><ul>")
		if counts["critical"] > 0 { body.WriteString(fmt.Sprintf("<li>Critical: %d</li>", counts["critical"])) }
		if counts["high"] > 0 { body.WriteString(fmt.Sprintf("<li>High: %d</li>", counts["high"])) }
		if counts["medium"] > 0 { body.WriteString(fmt.Sprintf("<li>Medium: %d</li>", counts["medium"])) }
		if counts["low"] > 0 { body.WriteString(fmt.Sprintf("<li>Low: %d</li>", counts["low"])) }
		if counts["info"] > 0 { body.WriteString(fmt.Sprintf("<li>Info: %d</li>", counts["info"])) }
		body.WriteString("</ul>")

		// Sort findings by severity
		severityScore := map[string]int{"critical": 5, "high": 4, "medium": 3, "low": 2, "info": 1}
		sort.Slice(findings, func(i, j int) bool {
			return severityScore[findings[i].Severity] > severityScore[findings[j].Severity]
		})

		body.WriteString("<h3>Details</h3>")
		body.WriteString("<table border='1' cellpadding='5' style='border-collapse: collapse; width: 100%;'>")
		body.WriteString("<tr style='background-color: #f2f2f2; text-align: left;'><th>Severity</th><th>Name</th><th>Host</th><th>Template</th></tr>")
		for _, f := range findings {
			color := "#ffffff"
			switch f.Severity {
			case "critical": color = "#ffcccc"
			case "high": color = "#ffe6cc"
			case "medium": color = "#ffffcc"
			case "low": color = "#e6f2ff"
			case "info": color = "#f2f2f2"
			}
			body.WriteString(fmt.Sprintf("<tr style='background-color: %s;'><td><strong>%s</strong></td><td>%s</td><td>%s</td><td>%s</td></tr>", 
				color, 
				html.EscapeString(f.Severity), 
				html.EscapeString(f.Name), 
				html.EscapeString(f.Host), 
				html.EscapeString(f.TemplateID)))
		}
		body.WriteString("</table>")
	}

	body.WriteString("</body></html>")

	addr := fmt.Sprintf("%s:%s", host, port)
	err := smtp.SendMail(addr, auth, from, []string{to}, body.Bytes())
	if err != nil {
		log.Printf("Failed to send email report: %v", err)
	} else {
		log.Println("Successfully sent email report for target", target.Name)
	}
}
