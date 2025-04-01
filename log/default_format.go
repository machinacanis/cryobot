/*

default_format.go

默认的日志格式化样式，包含一个暗色的和一个亮色的格式化样式，用于适配不同的终端背景色
你可以按照这个格式化样式编写修改你自己的日志样式，或者直接使用logrus的默认格式化样式
总之是通过实现logrus.Formatter接口来实现的

*/

package log

import (
	"fmt"
	"github.com/machinacanis/cryobot/utils"
	"github.com/sirupsen/logrus"
)

var white = utils.RGB("#FFFFFF")      // 白色
var gray = utils.RGB("#808080")       // 灰色
var aquamarine = utils.RGB("#7FFFD4") // 浅蓝色
var lightcyan = utils.RGB("#E0FFFF")  // 浅青色
var turquoise = utils.RGB("#40E0D0")  // 浅绿色
var cyan = utils.RGB("#00FFFF")       // 青色
var red = utils.RGB("#FF0000")        // 红色
var deepred = utils.RGB("#8B0000")    // 深红色
var yellow = utils.RGB("#FFFF00")     // 黄色
var purple = utils.RGB("#800080")     // 紫色
var reset = "\033[0m"

// DefaultDarkFormatter 默认的暗色格式化样式
type DefaultDarkFormatter struct{}

func (f *DefaultDarkFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置不同的颜色
	var textColor string
	var levelColor string
	var levelText string

	switch entry.Level {
	case logrus.TraceLevel: // Trace等级，标识色为浅蓝色，文本色为浅蓝色
		levelColor = aquamarine
		textColor = aquamarine
		levelText = "🔨TRACE"
	case logrus.DebugLevel: // Debug等级，标识色为浅青色，文本色为浅绿色
		levelColor = lightcyan
		textColor = turquoise
		levelText = "🔍DEBUG"
	case logrus.InfoLevel: // Info等级，标识色为青色，文本色为白色
		levelColor = cyan
		textColor = white
		levelText = "🧊INFO_"
	case logrus.WarnLevel: // Warn等级，标识色为黄色，文本色为黄色
		levelColor = yellow
		textColor = yellow
		levelText = "⚠️WARN_"
	case logrus.ErrorLevel: // Error等级，标识色为红色，文本色为红色
		levelColor = red
		textColor = red
		levelText = "⛔ERROR"
	case logrus.FatalLevel: // Fatal等级，标识色为深红色，文本色为深红色
		levelColor = deepred
		textColor = deepred
		levelText = "💀FATAL"
	case logrus.PanicLevel: // Panic等级，标识色为紫色，文本色为紫色
		levelColor = purple
		textColor = purple
		levelText = "🏴‍☠️PANIC"
	default:
		levelColor = white
		textColor = white
		levelText = "UNKNOWN"
	}

	// 构建日志格式,可以按需修改
	logMsg := fmt.Sprintf(
		"%s• %s [%s%s%s] %s%s%s\n",
		gray,
		entry.Time.Format("2006-01-02 15:04:05"),
		levelColor,
		levelText,
		gray,
		textColor,
		entry.Message,
		reset,
	)

	return []byte(logMsg), nil
}

var black = utils.RGB("#000000")       // 黑色
var darkGray = utils.RGB("#A9A9A9")    // 深灰色
var lightBlue = utils.RGB("#ADD8E6")   // 浅蓝色
var lightGreen = utils.RGB("#90EE90")  // 浅绿色
var darkYellow = utils.RGB("#FF7F00")  // 浅黄色
var lightRed = utils.RGB("#FFA07A")    // 浅红色
var lightPurple = utils.RGB("#DDA0DD") // 浅紫色

// DefaultLightFormatter 默认的亮色格式化样式
type DefaultLightFormatter struct{}

func (f *DefaultLightFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据日志级别设置不同的颜色
	var textColor string
	var levelColor string
	var levelText string

	switch entry.Level {
	case logrus.TraceLevel: // Trace等级，标识色为浅蓝色，文本色为浅蓝色
		levelColor = lightBlue
		textColor = lightBlue
		levelText = "🔨TRACE"
	case logrus.DebugLevel: // Debug等级，标识色为浅绿色，文本色为浅绿色
		levelColor = lightGreen
		textColor = lightGreen
		levelText = "🔍DEBUG"
	case logrus.InfoLevel: // Info等级，标识色为黑色，文本色为黑色
		levelColor = black
		textColor = black
		levelText = "🧊INFO_"
	case logrus.WarnLevel: // Warn等级，标识色为浅黄色，文本色为浅黄色
		levelColor = darkYellow
		textColor = darkYellow
		levelText = "⚠️WARN_"
	case logrus.ErrorLevel: // Error等级，标识色为浅红色，文本色为浅红色
		levelColor = lightRed
		textColor = lightRed
		levelText = "⛔ERROR"
	case logrus.FatalLevel: // Fatal等级，标识色为深灰色，文本色为深灰色
		levelColor = darkGray
		textColor = darkGray
		levelText = "💀FATAL"
	case logrus.PanicLevel: // Panic等级，标识色为浅紫色，文本色为浅紫色
		levelColor = lightPurple
		textColor = lightPurple
		levelText = "🏴‍☠️PANIC"
	default:
		levelColor = black
		textColor = black
		levelText = "UNKNOWN"
	}

	// 构建日志格式,可以按需修改
	logMsg := fmt.Sprintf(
		"%s• %s [%s%s%s] %s%s%s\n",
		gray,
		entry.Time.Format("2006-01-02 15:04:05"),
		levelColor,
		levelText,
		gray,
		textColor,
		entry.Message,
		reset,
	)

	return []byte(logMsg), nil
}
