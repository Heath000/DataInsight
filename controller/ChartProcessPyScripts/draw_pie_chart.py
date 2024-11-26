# draw_pie_chart.py
import sys
import json
import matplotlib.pyplot as plt
import matplotlib

# 设置字体以支持汉字
matplotlib.rcParams['font.sans-serif'] = ['SimHei']  # 使用黑体
matplotlib.rcParams['axes.unicode_minus'] = False  # 解决保存图像时负号 '-' 显示为方块的问题

def draw_pie_chart(title, data):
    labels = [item['name'] for item in data]
    sizes = [item['value'] for item in data]

    fig, ax = plt.subplots()
    ax.pie(sizes, labels=labels, autopct='%1.1f%%', startangle=90)
    ax.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.

    plt.title(title)
    plt.savefig('pie_chart.png')

if __name__ == "__main__":
    params = json.loads(sys.argv[1])
    draw_pie_chart(params['title'], params['data'])