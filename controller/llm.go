package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//前端传入参数
//前端传入参数
//前端传入参数
/*
		聊天模块：
		POST请求：localhost:7077/llm/chat
		body raw json
		前端传入参数：
		{
    		"Prompt": "你好，请介绍一下你自己"
		}
		后端返回参数：
		{
    		"data": "您好，我是科大讯飞研发的认知智能大模型，我的名字叫讯飞星火认知大模型。我可以和人类进行自然交流，解答问题，高效完成各领域认知智能需求。",
    		"message": "success"
		}
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		生成报告模块：
		POST请求：localhost:7077/llm/report
		body raw json
		前端传入参数：
		{
	    	Algorithm string   // 算法名称
			Table Table        // 表格数据
		}
		Table struct {
			Title       string         // 表格标题
			Description string         // 表格描述
    		Data        [][]float64
    		Labels      []float64
    		PredictData [][]float64
		}
		例如你可以传入：
		{
			"algorithm": "线性回归预测",
    		"table": {
                "title": "GDP线性回归预测",
                "Description": "这个是描述福州市2024年GDP线性回归预测的结果",
                "data": [[1.0, 2.0, 3.0], [2.0, 3.0, 4.0], [3.0, 4.0, 5.0], [4.0, 5.0, 6.0]],
                "labels": [6.0, 9.0, 12.0, 15.0],
                "predict_data": [[5.0, 6.0, 7.0], [6.0, 7.0, 8.0]]
            }
		}
		后端返回参数：
		{
    		"data": "# GDP线性回归预测分析报告\n\n## 背景\n本次分析针对的是福州市2024年GDP的线性回归预测。通过收集过去几年的GDP数据，并应用线性回归模型进行未来一年的经济预测。\n\n## 数据集描述\n本报告使用了一组简化的数据来模拟福州市过去四年的GDP情况以及预测第五和第六年的GDP值。数据格式如下：\n\n- **年份**: 表示数据的年份。\n- **X1, X2, X3**: 这些为输入变量，可能代表影响GDP的其他经济指标。\n- **Y (GDP)**: 每年的总产出，即GDP。\n- **Predict_data**: 未来两年用于预测的输入变量。\n\n## 数据\n| 年份 | X1 | X2 | X3 | Y (GDP) | Predict_data |\n|------|----|----|----|---------|--------------|\n| 1    | 1  | 2  | 3  | 6       |              |\n| 2    | 2  | 3  | 4  | 9       |              |\n| 3    | 3  | 4  | 5  | 12      |              |\n| 4    | 4  | 5  | 6  | 15      |              |\n\n## 预测数据\n对于未来两年的预测数据如下：\n\n| 年份 | X1 | X2 | X3 |\n|------|----|----|----|\n| 5    | 5  | 6  | 7  |\n| 6    | 6  | 7  | 8  |\n\n## 线性回归模型\n使用以上数据，我们构建了一个线性回归模型来预测未来两年的GDP。线性回归是一种统计方法，它试图通过建立一个自变量和因变量之间的线性关系来对因变量做出预测。在本例中，自变量是X1, X2, X3，而因变量是GDP。\n\n### 模型公式\n假设模型形式为 \\( y = \\beta_0 + \\beta_1x_1 + \\beta_2x_2 + \\beta_3x_3 + \\epsilon \\)，其中\\( \\beta_0, \\beta_1, \\beta_2, \\beta_3 \\)是系数，\\( \\epsilon \\)是误差项。\n\n### 训练结果\n经过训练得到以下系数（假设）：\n- \\(\\beta_0 = 1\\)\n- \\(\\beta_1 = 0.5\\)\n- \\(\\beta_2 = 1\\)\n- \\(\\beta_3 = 1.5\\)\n\n### 预测计算\n根据上述模型和系数，对未来两年的预测计算如下：\n- 第5年: \\( y = 1 + 0.5*5 + 1*6 + 1.5*7 = 20.5 \\)\n- 第6年: \\( y = 1 + 0.5*6 + 1*7 + 1.5*8 = 23.5 \\)\n\n### 预测结果\n因此，预测福州市第五年的GDP为20.5亿元，第六年的GDP为23.5亿元。\n\n## 结论与建议\n从预测结果来看，福州市的经济增长呈现稳定上升趋势。建议市政府继续关注影响GDP的关键因素，如提高投资效率、促进消费等，以持续推动经济的健康发展。同时，应注意潜在的风险因素，制定相应的应对策略，确保经济增长的稳定性和可持续性。",
    		"message": "success"
		}
		请注意里面所有content内容加一起的tokens需要控制在8192以内否则会报错(1 token约等于1.5个中文汉字 或者 0.8个英文单词)
		请注意返回的数据需要是markdown格式
*/
type LlmController struct{}

var (
	hostUrl   = "wss://spark-api.xf-yun.com/v4.0/chat"
	appid     = "7e568c90"
	apiSecret = "MmMxYzg5NTY5YTk5MjYxYTUzNmQxMDJj"
	apiKey    = "72c36b595faac7d22548aa1cbc3c93d9"
)

type Table struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Data        [][]float64 `json:"data"`
	Labels      []float64   `json:"labels"`
	PredictData [][]float64 `json:"predict_data"`
}

type user_questions struct {
	Algorithm string `json:"algorithm"`
	Table     Table  `json:"table"`
}

type user_chats struct {
	Prompt string `json:"prompt"`
}

func (ctrl *LlmController) GetChat(c *gin.Context) {

	time.Sleep(1 * time.Second)
	var user_chat user_chats
	if err := c.ShouldBindJSON(&user_chat); err != nil {
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
	//user_prompt := string(tableJSON) + "以上表格是"+ user_question.Table.Title +"运用"+ user_question.Algorithm + "算法的数据，" + "请根据以上表格帮我生成一份数据分析报告,并结合标题的语境进行数据解读,以markdown格式返回"
	go func() {

		data := genParams1(appid, user_chat.Prompt)
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

// GetUser gets the user info
func (ctrl *LlmController) GetReport(c *gin.Context) {
	var user_question user_questions
	if err := c.ShouldBindJSON(&user_question); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 将 Table 转换为 JSON 字符串
	tableJSON, err := json.Marshal(user_question.Table)
	if err != nil {
		fmt.Println("Error marshalling table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing table data"})
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
	user_prompt := string(tableJSON) + "以上表格是" + user_question.Table.Title + "运用" + user_question.Algorithm + "算法的数据，" + "请根据以上表格帮我生成一份数据分析报告,并结合标题的语境进行数据解读,以markdown格式返回"
	go func() {

		data := genParams1(appid, user_prompt)
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
