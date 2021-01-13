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
	"evalMatches":[],
	"message":"service is down",
	"orgId":1,
	"panelId":47,
	"ruleId":6,
	"ruleName":"xxx alert",
	"ruleUrl":"http://xxxx:3000/d/NUg5SDJMk/quark-dashboard?tab=alert\u0026viewPanel=47\u0026orgId=1",
	"state":"ok",
	"tags":{},
	"title":"[OK] xxx alert"
}

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

var sent_count int = 0

func GwStat() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		stat_msg := "G2WW Server created by Nova Kwok is running! \nParsed & forwarded " + strconv.Itoa(sent_count) + " messages to WeChat Work!"
		c.Send(stat_msg)
		return
	}
}

func GwWorker() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		h := new(Hook)
		fmt.Println(c.Body())
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

		msgStr := fmt.Sprintf(`
		{
			"msgtype": "news",
			"news": {
			  "articles": [
				{
				  "title": "%s",
				  "description": "%s",
				  "url": "%s",
				  "picurl": "%s"
				}
			  ]
			}
		  }
		`, h.Title, h.Message, h.RuleURL, h.ImageURL)

		color := "warning"
		if h.State == "ok" {
			color = "info"
		}
		alertItem := fmt.Sprintf(`<font color=\"%s\">%s</font>`, color, h.Message)

		if h.ImageURL == "" {
			msgStr = fmt.Sprintf(`
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
		}
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
		sent_count++
	}
}
