import sys
import json
from sklearn.linear_model import LinearRegression
from sklearn.preprocessing import StandardScaler
import numpy as np

try:
    # 解析输入数据
    input_data = json.loads(sys.argv[1])
    data = np.array(input_data.get("data", []))               # 自变量数据
    labels = np.array(input_data.get("labels", []))           # 目标值
    predict_data = np.array(input_data.get("predict_data", []))  # 待预测数据

    # 检查输入数据维度
    if data.shape[0] != len(labels):
        raise ValueError("Number of samples in 'data' and 'labels' must match")
    if predict_data.shape[1] != data.shape[1]:
        raise ValueError("Number of features in 'predict_data' must match 'data'")

    # 数据标准化
    scaler = StandardScaler()
    data_scaled = scaler.fit_transform(data)
    predict_data_scaled = scaler.transform(predict_data)

    # 构建线性回归模型
    model = LinearRegression()
    model.fit(data_scaled, labels)

    # 进行预测
    predictions = model.predict(predict_data_scaled).tolist()

    # 返回预测结果
    print(json.dumps({"predictions": predictions}))

except Exception as e:
    # 如果出错，输出错误信息
    print(json.dumps({"error": str(e)}))
