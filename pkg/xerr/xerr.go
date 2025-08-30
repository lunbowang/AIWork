package xerr

import (
	"fmt"
	"runtime"
	"strconv"
)

// callerSkipOffset 调用者框架跳过的偏移量，用于正确获取调用栈信息
// 调整此值可以准确定位到实际触发错误的代码位置
const callerSkipOffset = 3

// Cause 定义错误原因接口，用于获取原始错误
// 实现此接口的错误类型可以通过Cause()方法获取底层错误
type Cause interface {
	Cause() error
}

// withMessage 包装错误信息的结构体，包含原始错误和附加信息
// 用于在错误传递过程中添加上下文信息，同时保留原始错误链
type withMessage struct {
	cause error
	msg   string
}

// newWithMessage 创建一个带有附加信息的错误包装器
func newWithMessage(skip int, err error, msg string) error {
	if err == nil {
		return nil
	}

	var path string
	// 获取调用者的框架信息（文件名和行号）
	f, ok := getCallerFrame(skip)
	if ok {
		path = f.File + ":" + strconv.Itoa(f.Line)
	}

	return &withMessage{
		cause: err,
		msg:   path + msg,
	}
}

// WithMessage 为错误添加附加信息，保留原始错误链
func WithMessage(err error, message string) error {
	return newWithMessage(1, err, message)
}

// WithMessagef 为错误添加格式化的附加信息，保留原始错误链
func WithMessagef(err error, format string, v ...any) error {
	return newWithMessage(1, err, fmt.Sprintf(format, v...))
}

// New 包装原始错误，添加调用位置信息但不附加额外消息
func New(err error) error {
	return newWithMessage(1, err, "")
}

// Error 实现error接口，返回错误的完整信息
func (w *withMessage) Error() string {
	return w.msg
}

// Cause 实现Cause接口，返回原始错误
func (w *withMessage) Cause() error {
	return w.cause
}

// getCallerFrame 获取调用者的框架信息（文件名、行号等）
func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	// 用于存储程序计数器的切片，长度1表示只获取一个调用帧
	pc := make([]uintptr, 1)
	// 从调用栈中获取程序计数器，skip+callerSkipOffset控制跳过的帧数
	numFrames := runtime.Callers(skip+callerSkipOffset, pc)
	if numFrames < 1 {
		return
	}

	// 将程序计数器转换为调用帧信息
	frame, _ = runtime.CallersFrames(pc).Next()
	// 判断是否成功获取到有效帧信息（PC不为0）
	return frame, frame.PC != 0
}
