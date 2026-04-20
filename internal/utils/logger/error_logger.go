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
	log.Printf("\n" + strings.Repeat("=", 80))
	log.Printf("❌ ERROR DETAILS")
	log.Printf(strings.Repeat("=", 80))
	log.Printf("📅 Time:       %s", info.Timestamp)
	log.Printf("🆔 Request ID: %s", info.RequestID)
	log.Printf("📍 Function:   %s", info.Function)
	log.Printf("📁 File:       %s:%d", info.File, info.Line)
	log.Printf("🌐 Endpoint:   %s %s", info.Method, info.Path)
	log.Printf("📊 Status:     %d", info.StatusCode)
	log.Printf("💬 Message:    %s", info.Message)
	
	if info.StackTrace != "" && info.StackTrace != "  No relevant stack trace" {
		log.Printf(strings.Repeat("-", 80))
		log.Printf("📚 Stack Trace:")
		log.Printf("%s", info.StackTrace)
	}
	
	log.Printf(strings.Repeat("=", 80) + "\n")
}

func LogInfo(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Printf("[INFO] %s | Request-ID: %s | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		requestID,
		message,
	)
}

func LogDebug(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Printf("[DEBUG] %s | Request-ID: %s | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		requestID,
		message,
	)
}

func LogWarning(c *fiber.Ctx, message string) {
	requestID := ""
	if c != nil {
		if val := c.Locals("request_id"); val != nil {
			requestID = val.(string)
		}
	}

	log.Printf("[WARN] %s | Request-ID: %s | %s",
		time.Now().Format("2006-01-02 15:04:05"),
		requestID,
		message,
	)
}