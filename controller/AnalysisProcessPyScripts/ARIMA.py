import sys
import json
import warnings
from statsmodels.tsa.arima.model import ARIMA

# 忽略 warnings
warnings.filterwarnings("ignore")

try:
    # 解析输入数据
    input_data = json.loads(sys.argv[1])
    time_series = input_data.get("time_series", [])  # 时间序列信息（如年份）
    data = input_data.get("data", [])               # 时间序列数据
    p = input_data.get("p", 1)                      # 自回归阶数
    d = input_data.get("d", 1)                      # 差分次数
    q = input_data.get("q", 1)                      # 移动平均阶数
    forecast_steps = input_data.get("forecast_steps", 5)  # 预测步数

    # 检查输入的时间序列和数据长度是否匹配
    if len(time_series) != len(data):
        raise ValueError("Length of time_series and data must match")

    # 创建并拟合 ARIMA 模型
    model = ARIMA(data, order=(p, d, q))
    model_fit = model.fit()
    
    # 预测未来值
    forecast = model_fit.forecast(steps=forecast_steps)
    predictions = forecast.tolist()

    # 推算未来的时间序列
    last_time = time_series[-1]
    time_interval = time_series[-1] - time_series[-2] if len(time_series) > 1 else 1
    future_time_series = [last_time + i * time_interval for i in range(1, forecast_steps + 1)]

    # 输出预测结果和未来时间序列
    print(json.dumps({
        "predictions": predictions,
        "future_time_series": future_time_series
    }))

except Exception as e:
    # 如果出错，输出错误信息
    print(json.dumps({"error": str(e)}))
