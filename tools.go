//
// Tencent is pleased to support the open source community by making trpc-agent-go available.
//
// Copyright (C) 2025 Tencent.  All rights reserved.
//
// trpc-agent-go is licensed under the Apache License Version 2.0.
//
//

package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"bytes"
	//"os"
	"encoding/base64"
)

// Constants for supported calculator operations.
const (
	opAdd      = "add"
	opSubtract = "subtract"
	opMultiply = "multiply"
	opDivide   = "divide"
)

// calculatorArgs holds the input for the calculator tool.
type calculatorArgs struct {
	Operation string  `json:"operation" jsonschema:"description=The operation to perform,enum=add,enum=subtract,enum=multiply,enum=divide,required"`
	A         float64 `json:"a" jsonschema:"description=First number operand,required"`
	B         float64 `json:"b" jsonschema:"description=Second number operand,required"`
}

// calculatorResult holds the output for the calculator tool.
type calculatorResult struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
}

// timeArgs holds the input for the time tool.
type timeArgs struct {
	Timezone string `json:"timezone" jsonschema:"description=Timezone or leave empty for local,required"`
}

// timeResult holds the output for the time tool.
type timeResult struct {
	Timezone string `json:"timezone"`
	Time     string `json:"time"`
	Date     string `json:"date"`
	Weekday  string `json:"weekday"`
}

// bashArgs holds the input for the bash tool.
type bashArgs struct {
	Command string `json:"command" jsonschema:"description=Bash command to execute,required"`
	Timeout int    `json:"timeout" jsonschema:"description=Timeout in seconds (default: 30),default=30"`
}

// bashResult holds the output for the bash tool.
type bashResult struct {
	Command   string `json:"command"`
	Output    string `json:"output"`
	ExitCode  int    `json:"exit_code"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
	Container string `json:"container,omitempty"`
}

// pythonArgs holds the input for the python tool.
type pythonArgs struct {
	Code    string `json:"code" jsonschema:"description=Python code to execute,required"`
	Timeout int    `json:"timeout" jsonschema:"description=Timeout in seconds (default: 30),default=30"`
}

// pythonResult holds the output for the python tool.
type pythonResult struct {
	Code      string `json:"code"`
	Output    string `json:"output"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
	Container string `json:"container,omitempty"`
}

// Calculator tool implementation.
// calculate performs the requested mathematical operation.
// It supports add, subtract, multiply, and divide operations.
func calculate(ctx context.Context, args calculatorArgs) (calculatorResult, error) {
	var result float64
	// Select operation based on input.
	switch strings.ToLower(args.Operation) {
	case opAdd:
		result = args.A + args.B
	case opSubtract:
		result = args.A - args.B
	case opMultiply:
		result = args.A * args.B
	case opDivide:
		if args.B != 0 {
			result = args.A / args.B
		}
	}
	return calculatorResult{
		Operation: args.Operation,
		A:         args.A,
		B:         args.B,
		Result:    result,
	}, nil
}

// Time tool implementation.
// getCurrentTime returns the current time for the specified timezone.
// If the timezone is invalid or empty, it defaults to local time.
func getCurrentTime(ctx context.Context, args timeArgs) (timeResult, error) {
	loc := time.Local
	zone := args.Timezone
	// Attempt to load the specified timezone.
	if zone != "" {
		var err error
		loc, err = time.LoadLocation(zone)
		if err != nil {
			loc = time.Local
		}
	}
	now := time.Now().In(loc)
	return timeResult{
		Timezone: loc.String(),
		Time:     now.Format("15:04:05"),
		Date:     now.Format("2006-01-02"),
		Weekday:  now.Weekday().String(),
	}, nil
}

// checkDockerAvailable checks if Docker is available on the system.
func checkDockerAvailable() error {
	cmd := exec.Command("docker", "--version")
	return cmd.Run()
}

// // executeInContainer executes a command in a Docker container with security restrictions.
// func executeInContainer(ctx context.Context, image string, command []string, input string, timeout int) (string, string, int, error) {
// 	if timeout == 0 {
// 		timeout = 20
// 	}

// 	// Create context with timeout
// 	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
// 	defer cancel()

// 	// Generate unique container name
// 	containerName := fmt.Sprintf("agent_exec_%d", time.Now().UnixNano())

// 	// Prepare Docker command
// 	dockerArgs := []string{
// 		"run",
// 		"--rm",           // Remove container after execution
// 		"--name", containerName,
// 		"--network", "none", // No network access for security
// 		"--memory", "200m", // Memory limit
// 		"--cpus", "1", // CPU limit
// 		"--read-only", // Read-only filesystem
// 		"--tmpfs", "/tmp:rw,size=50m", // Temporary writeable directory
// 		"--user", "1000:1000", // Non-root user
// 		"--workdir", "/workspace",
// 		"-i", // Keep STDIN open
// 	}

// 	// Add image and command
// 	dockerArgs = append(dockerArgs, image)
// 	dockerArgs = append(dockerArgs, command...)

// 	cmd := exec.CommandContext(timeoutCtx, "docker", dockerArgs...)
	
// 	// Set input
// 	cmd.Stdin = strings.NewReader(input)
	
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr

// 	err := cmd.Run()
	
// 	var exitCode int
// 	if err != nil {
// 		if exitErr, ok := err.(*exec.ExitError); ok {
// 			exitCode = exitErr.ExitCode()
// 		} else {
// 			exitCode = -1
// 		}
// 	} else {
// 		exitCode = 0
// 	}

// 	return stdout.String(), stderr.String(), exitCode, err
// }

// Bash tool implementation with container isolation.
// executeBash executes a bash command in a secure Docker container.
func executeBash(ctx context.Context, args bashArgs) (bashResult, error) {
    // Check if Docker is available
    if err := checkDockerAvailable(); err != nil {
        return bashResult{}, fmt.Errorf("Docker is not available: %v", err)
    }

    if args.Timeout == 0 {
        args.Timeout = 30
    }

    // Create context with timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(args.Timeout)*time.Second)
    defer cancel()

    // Generate unique container name
    containerName := fmt.Sprintf("bash_exec_%d", time.Now().UnixNano())

    // 使用本地镜像
    image := "alpine:latest"

    // 使用更安全的方式传递命令
    command := []string{"sh", "-c", args.Command}

    // 准备 Docker 命令参数
    dockerArgs := []string{
        "run",
        "--rm",
        "--name", containerName,
        "--network", "none",
        "--memory", "200m",
        "--cpus", "1",
        "--read-only",
        "--tmpfs", "/tmp:rw,size=50m,uid=1000,gid=1000",
        "--user", "1000:1000",
        image,
    }
    dockerArgs = append(dockerArgs, command...)

    cmd := exec.CommandContext(timeoutCtx, "docker", dockerArgs...)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    execErr := cmd.Run()
    
    var exitCode int
    if execErr != nil {
        if exitErr, ok := execErr.(*exec.ExitError); ok {
            exitCode = exitErr.ExitCode()
        } else {
            exitCode = -1
        }
    } else {
        exitCode = 0
    }

    result := bashResult{
        Command:   args.Command,
        Output:    strings.TrimSpace(stdout.String()),
        ExitCode:  exitCode,
        Timestamp: time.Now().Format(time.RFC3339),
        Container: image,
    }

    if execErr != nil {
        errorMsg := strings.TrimSpace(stderr.String())
        if errorMsg == "" {
            errorMsg = execErr.Error()
        }
        result.Error = errorMsg
    }

    return result, nil
}

// Python tool implementation with container isolation.
// executePython executes Python code in a secure Docker container.
func executePython(ctx context.Context, args pythonArgs) (pythonResult, error) {
    // Check if Docker is available
    if err := checkDockerAvailable(); err != nil {
        return pythonResult{}, fmt.Errorf("Docker is not available: %v", err)
    }

    // 使用本地镜像
    image := "python:3-alpine"
    
    if args.Timeout == 0 {
        args.Timeout = 30
    }

    // 预处理代码：尝试恢复可能丢失的换行符和缩进
    processedCode := preprocessPythonCode(args.Code)

    // Create context with timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(args.Timeout)*time.Second)
    defer cancel()

    // Generate unique container name
    containerName := fmt.Sprintf("python_exec_%d", time.Now().UnixNano())

    // 使用 Base64 编码避免特殊字符问题
    encodedCode := base64.StdEncoding.EncodeToString([]byte(processedCode))
    command := []string{"sh", "-c"}
    
    // 创建 Python 脚本并执行
    fullCommand := fmt.Sprintf(`
echo %s | base64 -d > /tmp/script.py && 
cd /tmp && python script.py
`, encodedCode)

    // 准备 Docker 命令参数
    dockerArgs := []string{
        "run",
        "--rm",
        "--name", containerName,
        "--network", "none",
        "--memory", "200m",
        "--cpus", "1",
        "--read-only",
        "--tmpfs", "/tmp:rw,size=50m,uid=1000,gid=1000",
        "--user", "1000:1000",
        image,
    }
    dockerArgs = append(dockerArgs, command...)
    dockerArgs = append(dockerArgs, fullCommand)

    cmd := exec.CommandContext(timeoutCtx, "docker", dockerArgs...)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    execErr := cmd.Run()
    
    result := pythonResult{
        Code:      args.Code, // 保留原始代码用于显示
        Output:    strings.TrimSpace(stdout.String()),
        Timestamp: time.Now().Format(time.RFC3339),
        Container: image,
    }

    if execErr != nil {
        errorMsg := strings.TrimSpace(stderr.String())
        if errorMsg == "" {
            errorMsg = execErr.Error()
        }
        result.Error = errorMsg
    }

    return result, nil
}

// preprocessPythonCode 尝试恢复可能丢失的换行符和缩进
func preprocessPythonCode(code string) string {
    // 如果代码看起来已经被压缩成一行（包含多个语句但没有换行）
    if strings.Contains(code, ":") && !strings.Contains(code, "\n") {
        // 尝试基于 Python 语法关键词添加换行符
        keywords := []string{"for ", "if ", "def ", "class ", "while ", "elif ", "else:"}
        
        processed := code
        for _, keyword := range keywords {
            // 在关键词前添加换行符（除了第一个）
            if strings.Count(processed, keyword) > 1 {
                // 找到除了第一个之外的所有关键词位置
                parts := strings.Split(processed, keyword)
                for i := 1; i < len(parts); i++ {
                    parts[i] = "\n" + keyword + parts[i]
                }
                processed = strings.Join(parts, "")
            }
        }
        
        // 添加基本的缩进
        lines := strings.Split(processed, "\n")
        for i := 1; i < len(lines); i++ {
            if strings.HasPrefix(lines[i], "    ") || strings.HasPrefix(lines[i], "\t") {
                continue // 已经有缩进
            }
            // 如果上一行以冒号结束，当前行应该缩进
            if i > 0 && strings.HasSuffix(strings.TrimSpace(lines[i-1]), ":") {
                lines[i] = "    " + lines[i]
            }
        }
        
        return strings.Join(lines, "\n")
    }
    
    return code
}

// intPtr returns a pointer to the given int value.
func intPtr(i int) *int {
	return &i
}

// floatPtr returns a pointer to the given float64 value.
func floatPtr(f float64) *float64 {
	return &f
}

// This example demonstrates how to integrate tRPC agent orchestration
// with LLM-based tools, providing a simple HTTP server for manual
// testing. It is intended as a reference for developers looking to build
// custom LLM agents with tool support in Go.
//
// The calculator tool supports basic arithmetic operations, while the
// time tool provides current time information for a given timezone.
// The bash tool executes shell commands in a secure container, and the 
// python tool executes Python code snippets in a secure container.
//
// The code is structured for clarity and ease of extension.