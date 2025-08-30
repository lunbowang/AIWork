package main

import (
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	// 创建一个新的图
	p := plot.New()

	// 设置标题和标签
	p.Title.Text = "bar test"
	p.X.Label.Text = "sex"
	p.Y.Label.Text = "value"

	// 数据
	values := []float64{10, 20, 15, 25}
	labels := []string{"a", "b", "c", "d"}

	// 创建柱状图
	bars, err := plotter.NewBarChart(plotter.Values(values), 10)
	if err != nil {
		log.Fatalf("could not create bar chart: %v", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	// 将柱状图添加到图表中
	p.Add(bars)
	// 设置X轴标签
	p.NominalX(labels...)
	// 保存图形到文件
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "barchart.png"); err != nil {
		log.Fatalf("could not save plot to file: %v", err)
	}
}
