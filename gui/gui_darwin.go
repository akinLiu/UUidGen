//go:build darwin

package gui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"uuidgen/sysinfo"
)

var (
	// 科技感深色主题配色
	bgColor       = color.NRGBA{R: 12, G: 12, B: 16, A: 255}    // 深黑背景
	cardColor     = color.NRGBA{R: 22, G: 28, B: 38, A: 255}    // 带蓝调的深色卡片
	accentColor   = color.NRGBA{R: 0, G: 220, B: 255, A: 255}   // 霓虹青色
	accentHover   = color.NRGBA{R: 0, G: 180, B: 220, A: 255}   // 悬停色
	textPrimary   = color.NRGBA{R: 0, G: 220, B: 255, A: 255}   // 霓虹青色标题
	textSecondary = color.NRGBA{R: 100, G: 140, B: 180, A: 255} // 科技蓝灰色
	textUUID      = color.NRGBA{R: 0, G: 255, B: 136, A: 255}   // 霓虹绿色UUID
	successColor  = color.NRGBA{R: 0, G: 255, B: 136, A: 255}   // 成功色
)

// Run starts the Fyne GUI displaying the given UUID string.
func Run(uuidStr string, uuidErr error) {
	runWithInfo(uuidStr, uuidErr, nil)
}

// RunWithSystemInfo starts the GUI with full system information
func RunWithSystemInfo(info *sysinfo.SystemInfo) {
	runWithInfo(info.UUID, nil, info)
}

func runWithInfo(uuidStr string, uuidErr error, sysInfo *sysinfo.SystemInfo) {
	a := app.NewWithID("com.uuidgen.app")
	w := a.NewWindow("UUidGen")
	w.Resize(fyne.NewSize(520, 320))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	if uuidErr != nil {
		showError(w, uuidErr)
		w.ShowAndRun()
		return
	}

	showUUID(w, uuidStr)
	w.ShowAndRun()
}

func showUUID(w fyne.Window, uuid string) {
	// 背景
	bg := canvas.NewRectangle(bgColor)

	// 标题
	title := canvas.NewText("Device Identifier", textPrimary)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// 副标题
	subtitle := canvas.NewText("SMBIOS UUID", textSecondary)
	subtitle.TextSize = 13
	subtitle.Alignment = fyne.TextAlignCenter

	// 卡片容器
	card := canvas.NewRectangle(cardColor)
	card.CornerRadius = 12

	// UUID 显示
	uuidText := canvas.NewText(uuid, textUUID)
	uuidText.TextSize = 16
	uuidText.TextStyle = fyne.TextStyle{Monospace: true}
	uuidText.Alignment = fyne.TextAlignCenter

	// UUID 容器
	uuidContainer := container.NewCenter(uuidText)

	// 复制按钮
	copyBtn := widget.NewButtonWithIcon("Copy UUID", theme.ContentCopyIcon(), nil)
	copyBtn.Importance = widget.HighImportance
	copyBtn.Resize(fyne.NewSize(140, 36))
	copyBtn.OnTapped = func() {
		w.Clipboard().SetContent(uuid)
		copyBtn.SetText("Copied!")
		copyBtn.Importance = widget.SuccessImportance
		copyBtn.Refresh()

		time.AfterFunc(2*time.Second, func() {
			copyBtn.SetText("Copy UUID")
			copyBtn.Importance = widget.HighImportance
			copyBtn.Refresh()
		})
	}

	// 卡片内容
	cardContent := container.NewVBox(
		layout.NewSpacer(),
		uuidContainer,
		layout.NewSpacer(),
	)
	cardContent.Resize(fyne.NewSize(440, 60))

	// 卡片层叠
	cardStack := container.NewStack(
		card,
		container.NewPadded(cardContent),
	)

	// 主内容
	content := container.NewStack(
		bg,
		container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(title),
			container.NewCenter(subtitle),
			layout.NewSpacer(),
			container.NewPadded(cardStack),
			layout.NewSpacer(),
			container.NewCenter(copyBtn),
			layout.NewSpacer(),
		),
	)

	w.SetContent(content)
}

func showError(w fyne.Window, err error) {
	bg := canvas.NewRectangle(bgColor)

	title := canvas.NewText("Unable to Retrieve UUID", textPrimary)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// 错误图标背景
	iconBg := canvas.NewCircle(color.NRGBA{R: 255, G: 69, B: 58, A: 50})
	iconText := canvas.NewText("✕", color.NRGBA{R: 255, G: 69, B: 58, A: 255})
	iconText.TextSize = 32
	iconText.TextStyle = fyne.TextStyle{Bold: true}
	iconText.Alignment = fyne.TextAlignCenter

	iconContainer := container.NewStack(
		iconBg,
		iconText,
	)
	iconContainer.Resize(fyne.NewSize(64, 64))

	errLabel := canvas.NewText(err.Error(), textSecondary)
	errLabel.TextSize = 14
	errLabel.Alignment = fyne.TextAlignCenter

	content := container.NewStack(
		bg,
		container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(iconContainer),
			layout.NewSpacer(),
			container.NewCenter(title),
			container.NewCenter(errLabel),
			layout.NewSpacer(),
		),
	)

	w.SetContent(content)
}
