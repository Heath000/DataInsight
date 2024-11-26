import sys
import json
import matplotlib.pyplot as plt
import numpy as np

def draw_line_bar_mixed_chart(title, line_data, bar_data):
    fig, ax1 = plt.subplots()

    # Plot bar chart
    bar_labels = [item['label'] for item in bar_data]
    bar_values = [item['value'] for item in bar_data]
    x = np.arange(len(bar_labels))  # the label locations
    width = 0.35  # the width of the bars

    ax1.bar(x, bar_values, width, label='Bar', color='b', alpha=0.6)
    ax1.set_xlabel('Labels')
    ax1.set_ylabel('Bar Values', color='b')

    # Plot line chart
    line_x = [item['x'] for item in line_data]
    line_y = [item['y'] for item in line_data]
    ax2 = ax1.twinx()
    ax2.plot(x, line_y, color='r', marker='o', label='Line')
    ax2.set_ylabel('Line Values', color='r')

    # Set x-ticks and x-tick labels for both charts
    ax1.set_xticks(x)
    ax1.set_xticklabels(bar_labels)
    ax2.set_xticks(x)
    ax2.set_xticklabels(bar_labels)

    plt.title(title)
    plt.savefig('line_bar_mixed_chart.png')

if __name__ == "__main__":
    params = json.loads(sys.argv[1])
    draw_line_bar_mixed_chart(params['title'], params['line_data'], params['bar_data'])