package logger

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ErrorInfo struct {
	Timestamp  string
	RequestID  string
	Function   string
	File       string
	Line       int
	Message    string
	StackTrace string
	Path       string
	Method     string
	StatusCode int
}

// LogError logs detailed error information
func LogError(c *fiber.Ctx, err error, functionName string, skipFrames ...int) {
	skip := 2
	if len(skipFrames) > 0 {
		skip = skipFrames[0]
	}

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "unknown"
		line = 0
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	shortFuncName := parts[len(parts)-1]

	if functionName != "" {
		shortFuncName = functionName
	}

	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	stackTrace := getShortStackTrace()

	errorInfo := ErrorInfo{
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		RequestID:  requestID,
		Function:   shortFuncName,
		File:       getFileName(file),
		Line:       line,
		Message:    err.Error(),
		StackTrace: stackTrace,
	}

	if c != nil {
		errorInfo.Path = c.Path()
		errorInfo.Method = c.Method()
		errorInfo.StatusCode = c.Response().StatusCode()
	}

	printErrorLog(errorInfo)
}

func getShortStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	lines := strings.Split(stackTrace, "\n")
	var relevantLines []string

	relevantKeywords := []string{
		"/internal/",
		"/delivery/",
		"/usecase/",
		"/repository/",
		"/handlers/",
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.Contains(line, "github.com/gofiber/fiber") {
			break
		}
		if strings.Contains(line, "runtime/") {
			break
		}

		for _, keyword := range relevantKeywords {
			if strings.Contains(line, keyword) {
				cleanLine := cleanStackTraceLine(line)
				relevantLines = append(relevantLines, cleanLine)
				break
			}
		}
	}

	if len(relevantLines) == 0 {
		return "  No relevant stack trace"
	}

	return strings.Join(relevantLines, "\n")
}

func cleanStackTraceLine(line string) string {
	parts := strings.Split(line, "/github.com/savanyv/cartify/")
	if len(parts) > 1 {
		return "  " + parts[1]
	}

	if strings.TrimSpace(line) != "" {
		return "  " + line
	}
	return ""
}

func getFileName(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	return parts[len(parts)-1]
}

func printErrorLog(info ErrorInfo) {
	log.Print("\n" + strings.Repeat("=", 80))
	log.Print("❌ ERROR DETAILS")
	log.Print(strings.Repeat("=", 80))
	
	log.Print("📅 Time:       " + info.Timestamp)
	log.Print("🆔 Request ID: " + info.RequestID)
	log.Print("📍 Function:   " + info.Function)
	log.Print("📁 File:       " + info.File + ":" + intToString(info.Line))
	log.Print("🌐 Endpoint:   " + info.Method + " " + info.Path)
	log.Print("📊 Status:     " + intToString(info.StatusCode))
	log.Print("💬 Message:    " + info.Message)
	
	if info.StackTrace != "" && info.StackTrace != "  No relevant stack trace" {
		log.Print(strings.Repeat("-", 80))
		log.Print("📚 Stack Trace:")
		log.Print(info.StackTrace)
	}
	
	log.Print(strings.Repeat("=", 80) + "\n")
}

func intToString(n int) string {
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}

func LogInfo(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Print("[INFO] " + time.Now().Format("2006-01-02 15:04:05") + " | Request-ID: " + requestID + " | " + message)
}

func LogDebug(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Print("[DEBUG] " + time.Now().Format("2006-01-02 15:04:05") + " | Request-ID: " + requestID + " | " + message)
}

func LogWarning(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Print("[WARN] " + time.Now().Format("2006-01-02 15:04:05") + " | Request-ID: " + requestID + " | " + message)
}