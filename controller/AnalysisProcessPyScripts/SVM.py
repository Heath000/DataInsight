import sys
import json
from sklearn.svm import SVC
from sklearn.preprocessing import StandardScaler, LabelEncoder
import numpy as np

try:
    # 解析输入数据
    input_data = json.loads(sys.argv[1])
    data = np.array(input_data.get("data", []))               # 自变量数据
    labels = np.array(input_data.get("labels", []))           # 类别名称
    predict_data = np.array(input_data.get("predict_data", []))  # 待分类数据
    C = input_data.get("C", 1.0)                              # 惩罚系数
    tol = input_data.get("tol", 0.001)                        # 误差收敛条件
    max_iter = input_data.get("max_iter", 1000)               # 最大迭代次数

    # 检查输入数据维度
    if data.shape[0] != len(labels):
        raise ValueError("Number of samples in 'data' and 'labels' must match")
    if predict_data.shape[1] != data.shape[1]:
        raise ValueError("Number of features in 'predict_data' must match 'data'")

    # 编码类别名称
    label_encoder = LabelEncoder()
    labels_encoded = label_encoder.fit_transform(labels)

    # 数据标准化
    scaler = StandardScaler()
    data_scaled = scaler.fit_transform(data)
    predict_data_scaled = scaler.transform(predict_data)

    # 构建 SVM 分类模型
    model = SVC(C=C, tol=tol, max_iter=max_iter, probability=True, random_state=42)
    model.fit(data_scaled, labels_encoded)

    # 预测类别概率
    probabilities = model.predict_proba(predict_data_scaled)
    class_order = label_encoder.classes_.tolist()

    # 确定预测类别
    predictions_encoded = np.argmax(probabilities, axis=1)
    predictions = label_encoder.inverse_transform(predictions_encoded)

    # 格式化结果
    result_probabilities = [class_order] + probabilities.tolist()

    # 返回预测结果
    print(json.dumps({
        "predictions": predictions.tolist(),
        "probabilities": result_probabilities
    }))

except Exception as e:
    # 如果出错，输出错误信息
    print(json.dumps({"error": str(e)}))
