import sys
import json
from sklearn.ensemble import RandomForestClassifier
from sklearn.preprocessing import StandardScaler, LabelEncoder
import numpy as np

try:
    # 解析输入数据
    input_data = json.loads(sys.argv[1])
    data = np.array(input_data.get("data", []))               # 自变量数据
    labels = np.array(input_data.get("labels", []))           # 类别名称
    predict_data = np.array(input_data.get("predict_data", []))  # 待分类数据
    n_estimators = input_data.get("n_estimators", 100)        # 决策树数量
    max_depth = input_data.get("max_depth", None)             # 树的最大深度
    max_leaf_nodes = input_data.get("max_leaf_nodes", None)   # 叶节点的最大数量
    min_samples_split = input_data.get("min_samples_split", 2)  # 内部节点分裂的最小样本数
    min_samples_leaf = input_data.get("min_samples_leaf", 1)  # 叶节点的最小样本数

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

    # 构建随机森林分类模型
    model = RandomForestClassifier(
        n_estimators=n_estimators,
        max_depth=max_depth,
        max_leaf_nodes=max_leaf_nodes,
        min_samples_split=min_samples_split,
        min_samples_leaf=min_samples_leaf,
        random_state=42
    )
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
