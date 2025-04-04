package cryobot

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"strconv"
	"strings"
)

// RGB 得到一个RGB颜色的ANSI颜色代码
func RGB(rgb string) string {
	// 如果RGB字符串以#开头则去掉
	if strings.HasPrefix(rgb, "#") {
		rgb = rgb[1:]
	}

	// 将RGB字符串转换为整数
	r, err := strconv.ParseInt(rgb[0:2], 16, 64)
	if err != nil {
		r = 255
	}
	g, err := strconv.ParseInt(rgb[2:4], 16, 64)
	if err != nil {
		g = 255
	}
	b, err := strconv.ParseInt(rgb[4:6], 16, 64)
	if err != nil {
		b = 255
	}
	// 返回ANSI颜色代码
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func RandomDeviceNumber() int {
	return rand.Intn(9999999-1000000+1) + 1000000
}

// NewUUID 生成一个新的UUID
func NewUUID() string {
	return uuid.NewV4().String()
}

// newNickname 生成一个新的昵称
func newNickname() string {
	botCount := len(connectedClients)
	return fmt.Sprintf("Bot%d", botCount)
}

// GetQRCodeString 生成二维码字符串
//
// 基于 https://github.com/Baozisoftware/qrcode-terminal-go 修改而来
func GetQRCodeString(content string) (result *string) {
	var qr *qrcode.QRCode
	var err error
	qr, err = qrcode.New(content, qrcode.Low)
	if err != nil {
		return nil
	}
	data := qr.Bitmap()

	str := ""
	for ir, row := range data {
		lr := len(row)
		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}
		for ic, col := range row {
			lc := len(data)
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				str += fmt.Sprint("\033[48;5;0m  \033[0m") // 前景色
			} else {
				str += fmt.Sprint("\033[48;5;7m  \033[0m") // 背景色
			}
		}
		str += fmt.Sprintln()
	}
	return &str
}
