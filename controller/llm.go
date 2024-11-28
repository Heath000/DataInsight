package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
		POST请求：localhost:7077/llm/report
		body raw json
		前端传入参数：
		{
	    	"prompt": 用户的提示词
		}
		例如你可以传入：
		{
	    "prompt": "你好，请介绍一下你自己"
		}
		后端返回参数：
		{
	    "data": "您好，我是科大讯飞研发的认知智能大模型，我的名字叫讯飞星火认知大模型。我可以和人类进行自然交流，解答问题，高效完成各领域认知智能需求。",
	    "message": "success"
		}
		请注意prompt必须是string类型且里面的所有content内容加一起的tokens需要控制在8192以内否则会报错(1 token约等于1.5个中文汉字 或者 0.8个英文单词)
		如果要实现固定的提问词，之间在前端以字符串的形式拼接到prompt上即可，例如：前端先把用户的表格数据拼接成字符串，然后加上固定的提问词，
		例如，请根据以上表格帮我生成一份数据分析报告，然后把表格和固定的提问词拼接起来传入prompt即可
*/
type LlmController struct{}

var (
	hostUrl   = "wss://spark-api.xf-yun.com/v4.0/chat"
	appid     = "7e568c90"
	apiSecret = "MmMxYzg5NTY5YTk5MjYxYTUzNmQxMDJj"
	apiKey    = "72c36b595faac7d22548aa1cbc3c93d9"
)

type user_questions struct {
	Prompt string `json:"prompt"`
}

// GetUser gets the user info
func (ctrl *LlmController) GetReport(c *gin.Context) {
	var user_question user_questions
	if err := c.ShouldBindJSON(&user_question); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/*
		response := map[string]string{
			"response": "This is a response to the prompt: " + user_question.Prompt,
		}
		c.JSON(http.StatusOK, response)
	*/
	// fmt.Println(HmacWithShaTobase64("hmac-sha256", "hello\nhello", "hello"))
	// st := time.Now()
	//创建一个新的 websocket.Dialer 实例，并设置了握手超时时间。
	//websocket.Dialer 用于配置和建立 WebSocket 连接。
	//HandshakeTimeout 字段设置了 WebSocket 握手完成的最长持续时间。
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		panic(readResp(resp) + err.Error())
		return
	} else if resp.StatusCode != 101 {
		panic(readResp(resp) + err.Error())
	}
	//data := genParams1(appid, "你是谁，可以干什么？"): 调用 genParams1 函数，传入 appid 和一个问题字符串，生成要发送的数据。
	//conn.WriteJSON(data): 通过 WebSocket 连接 conn 发送生成的 JSON 数据。
	//这段代码的目的是异步发送一个消息，询问 "你是谁，可以干什么？" 给 WebSocket 服务器。
	go func() {

		data := genParams1(appid, user_question.Prompt)
		conn.WriteJSON(data)

	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}

		var data map[string]interface{}
		/*
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
				return
			}
		*/
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err1)
			return
		}
		//fmt.Println(string(msg))
		//解析数据
		payload, ok := data["payload"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			fmt.Println(data["payload"])
			return
		}
		status := choices["status"].(float64)
		//fmt.Println(status)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			fmt.Println("收到最终结果")
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	//输出返回结果
	fmt.Println(answer)

	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    answer,
	})
}

// 生成参数
func genParams1(appid, question string) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      "4.0Ultra",   // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8), // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),     // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),  // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",    // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
