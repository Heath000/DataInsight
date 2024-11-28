import sys
import json
import numpy as np
from scipy.stats import kurtosis, skew

def calculate_statistics(data):
    stats = {}
    data = np.array(data)

    stats["count"] = int(len(data))
    stats["max"] = float(np.max(data))
    stats["min"] = float(np.min(data))
    stats["mean"] = float(np.mean(data))
    stats["std_dev"] = float(np.std(data, ddof=1) ) # 样本标准差
    stats["median"] = float(np.median(data))
    stats["variance"] = float(np.var(data, ddof=1) ) # 样本方差
    stats["kurtosis"] = float(kurtosis(data))
    stats["skewness"] = float(skew(data))
    stats["cv"] = stats["std_dev"] / stats["mean"] if stats["mean"] != 0 else float('nan')  # 变异系数
    stats["cv"] = float(stats["cv"])
    return stats

try:
    # 从命令行接收 JSON 数据
    input_data = json.loads(sys.argv[1])

    # 获取输入数据
    data = input_data.get("data", [])

    # 计算统计信息
    stats = calculate_statistics(data)
    
    # 输出统计结果
    print(json.dumps(stats))
except Exception as e:
    # 如果出错，输出错误信息
    print(json.dumps({"error": str(e)}))
