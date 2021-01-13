package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber"
)

type Hook struct {
	dashboardId int64         `json:"dashboardId"`
	evalMatches []interface{} `json:"evalMatches"`
	ImageUrl    string        `json:"imageUrl"`
	Message     string        `json:"message"`
	orgId       int64         `json:"orgId"`
	panelId     int64         `json:"panelId"`
	ruleId      int64         `json:"ruleId"`
	ruleName    string        `json:"ruleName"`
	RuleUrl     string        `json:"ruleUrl"`
	state       string        `json:"state"`
	tags        string        `json:"tags"`
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
		`, h.Title, h.Message, h.RuleUrl, h.ImageUrl)

		color := "warning"
		if h.state == "ok" {
			color = "info"
		}
		alertItem := fmt.Sprintf(`<font color=\"%s\">%s</font>`, color, h.Message)

		if h.ImageUrl == "" {
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
			`, h.Title, alertItem, now, h.RuleUrl)
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
