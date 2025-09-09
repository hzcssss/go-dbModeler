package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go-DBmodeler/internal/db/connector"
)

// TableView 表示表结构视图
type TableView struct {
	widget.BaseWidget
	container *fyne.Container
	metadata  *connector.TableMetadata
}

// NewTableView 创建一个新的表结构视图
func NewTableView(metadata *connector.TableMetadata) *TableView {
	view := &TableView{
		metadata: metadata,
	}
	view.ExtendBaseWidget(view)
	view.container = view.buildUI()
	return view
}

// CreateRenderer 创建渲染器
func (v *TableView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.container)
}

// buildUI 构建UI
func (v *TableView) buildUI() *fyne.Container {
	// 创建表头
	header := container.NewHBox(
		widget.NewLabel("字段名"),
		widget.NewLabel("类型"),
		widget.NewLabel("可空"),
		widget.NewLabel("主键"),
		widget.NewLabel("默认值"),
		widget.NewLabel("注释"),
	)

	// 创建字段列表
	rows := make([]fyne.CanvasObject, 0, len(v.metadata.Fields))
	for _, field := range v.metadata.Fields {
		// 创建字段行
		row := container.NewHBox(
			widget.NewLabel(field.Name),
			widget.NewLabel(field.Type),
			widget.NewCheck("", nil),
			widget.NewCheck("", nil),
			widget.NewLabel(field.Default),
			widget.NewLabel(field.Comment),
		)

		// 设置复选框状态
		row.Objects[2].(*widget.Check).SetChecked(field.IsNullable)
		row.Objects[3].(*widget.Check).SetChecked(field.IsPrimary)

		// 禁用复选框（只读）
		row.Objects[2].(*widget.Check).Disable()
		row.Objects[3].(*widget.Check).Disable()

		rows = append(rows, row)
	}

	// 创建表格
	table := container.NewVBox(append([]fyne.CanvasObject{header}, rows...)...)

	// 创建滚动容器
	scroll := container.NewScroll(table)

	return container.NewBorder(
		widget.NewLabel("表名: "+v.metadata.Name),
		nil,
		nil,
		nil,
		scroll,
	)
}

// SetMetadata 设置表元数据
func (v *TableView) SetMetadata(metadata *connector.TableMetadata) {
	v.metadata = metadata
	// 重新构建UI
	newContainer := v.buildUI()

	// 更新容器引用
	v.container = newContainer

	// 强制刷新整个widget
	v.Refresh()
}
