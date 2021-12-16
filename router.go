package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber"
)

/*
{
  "dashboardId":1,
  "evalMatches":[
    {
      "value":1,
      "metric":"Count",
      "tags":{}
    }
  ],
  "imageUrl":"https://grafana.com/assets/img/blog/mixed_styles.png",
  "message":"Notification Message",
  "orgId":1,
  "panelId":2,
  "ruleId":1,
  "ruleName":"Panel Title alert",
  "ruleUrl":"http://localhost:3000/d/hZ7BuVbWz/test-dashboard?fullscreen\u0026edit\u0026tab=alert\u0026panelId=2\u0026orgId=1",
  "state":"alerting",
  "tags":{
    "tag name":"tag value"
  },
  "title":"[Alerting] Panel Title alert"
}
Reference: https://grafana.com/docs/grafana/latest/alerting/notifications/
*/

// Hook webhook body
type Hook struct {
	DashboardID int64         `json:"dashboardId"`
	EvalMatches []interface{} `json:"evalMatches"`
	ImageURL    string        `json:"imageUrl"`
	Message     string        `json:"message"`
	OrgID       int64         `json:"orgId"`
	PanelID     int64         `json:"panelId"`
	RuleID      int64         `json:"ruleId"`
	RuleName    string        `json:"ruleName"`
	RuleURL     string        `json:"ruleUrl"`
	State       string        `json:"state"`
	Tags        interface{}   `json:"tags"`
	Title       string        `json:"title"`
}

var sentCount int = 0

// GwStat stats check
func GwStat() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		statMsg := "G2WW Server created by Nova Kwok is running! \nParsed & forwarded " + strconv.Itoa(sentCount) + " messages to WeChat Work!"
		c.Send(statMsg)
		return
	}
}

// GwWorker Send to WeChat Work
func GwWorker() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		h := new(Hook)
		if err := c.BodyParser(h); err != nil {
			fmt.Println(err)
			c.Send("Error on JSON format")
			return
		}

		// Send to WeChat Work

		// {
		// 	"msgtype": "news",
		// 	"news": {
		// 	  "articles": [
		// 		{
		// 		  "title": "%s",
		// 		  "description": "%s",
		// 		  "url": "%s",
		// 		  "picurl": "%s"
		// 		}
		// 	  ]
		// 	}
		//   }
		// {
		// 	"msgtype": "text",
		// 	"text": {
		// 		"content": "%s\n%s\n查看详情:%s",
		// 		"mentioned_list":["@all"],
		// 	}
		//   }

		url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + c.Params("key")
		now := time.Now().Format("2006-01-02 15:04:05")

		// msgStr := fmt.Sprintf(`
		// {
		// 	"msgtype": "news",
		// 	"news": {
		// 	  "articles": [
		// 		{
		// 		  "title": "%s",
		// 		  "description": "%s",
		// 		  "url": "%s",
		// 		  "picurl": "%s"
		// 		}
		// 	  ]
		// 	}
		//   }
		// `, h.Title, h.Message, h.RuleURL, h.ImageURL)

		color := "warning"
		if h.State == "ok" {
			color = "info"
		}
		alertItem := fmt.Sprintf(`<font color=\"%s\">%s</font>`, color, h.Message)

		// if h.ImageURL == "" {
		msgStr := fmt.Sprintf(`
			{
				"msgtype": "markdown",
				"markdown": {
					"content": "**%s**\n
					>%s
					>时间: %s
					>[点击查看详情](%s)",
				}
			  }
			`, h.Title, alertItem, now, h.RuleURL)
		// }
		fmt.Println(msgStr)
		jsonStr := []byte(msgStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Send("Error sending to WeChat Work API")
			return
		}
		defer resp.Body.Close()
		c.Send(resp)
		sentCount++
	}
}
