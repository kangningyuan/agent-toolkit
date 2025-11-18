English | [中文](README_zh.md)
# Agent Toolkit

A powerful AI agent toolkit built on the [trpc-agent-go](https://github.com/trpc-group/trpc-agent-go) framework.  
It integrates practical utilities (calculator, time query, Bash & Python execution), runs inside secure containers, and is fully compatible with the [Google Agent Development Kit Web UI](https://github.com/google/adk-web).

<img src="./image/cover.png">

## Features

### 🛠️ Core Tools
- **Calculator** – add, subtract, multiply, divide  
- **Time Query** – current time, date and weekday for any timezone  
- **Bash Executor** – run shell commands in a locked-down Docker container  
- **Python Executor** – run Python code in the same isolated environment  

### 🔒 Security
- **Container isolation** – every execution happens inside its own Docker container  
- **Resource limits** – memory & CPU caps to prevent abuse  
- **Network isolation** – no outbound network access  
- **Read-only filesystem** – stops any permanent modification  


## Quick Start

### Prerequisites
- Go 1.24  
- Docker  
- An accessible LLM API (OpenAI-compatible endpoint)

### Installation

1. **Get the trpc-agent-go framework**
```bash
git clone https://github.com/trpc-group/trpc-agent-go.git
cd trpc-agent-go
```

2. **Install Agent Toolkit**
```bash
# Copy Agent Toolkit files under trpc-agent-go, e.g.
cd examples/agent-toolkit

# Pull the runtime images
chmod +x pull-images.sh && ./pull-images.sh

# Smoke-test the setup
chmod +x test-setup.sh && ./test-setup.sh
```

3. **Start the Agent Toolkit service**
```bash
# Set up your LLM (in an OpenAI-compatible format)
export OPENAI_API_KEY="your-api-key-here"
export OPENAI_BASE_URL="your-base-url-here"

# run
go run examples/agent-toolkit/main.go examples/agent-toolkit/tools.go
```

4. **Launch the ADK Web UI**
```bash
# Leave trpc-agent-go directory
git clone https://github.com/google/adk-web.git
cd adk-web
npm install
npm run serve --backend=http://localhost:8080
```

5. **Open the UI**
Browse to `http://localhost:4200` and start chatting.

### CLI Options
```bash
# Custom model & port
go run examples/agent-toolkit/main.go examples/agent-toolkit/tools.go -model gpt-4 -addr :9090

# Defaults: deepseek-chat on :8080
go run examples/agent-toolkit/main.go examples/agent-toolkit/tools.go
```

## Project Layout
```
trpc-agent-go/
├── examples/
│   └── agent-toolkit/
│       ├── main.go          # entry point & agent config
│       ├── tools.go         # tool implementations
│       ├── pull-images.sh   # docker pull helper
│       ├── test-setup.sh    # env sanity check
│       └── README.md        # this file
├── ...                      # other framework code
└── go.mod
```

## Developer Guide

### Adding a New Tool

1. Define argument & result structs in `tools.go`  
2. Implement the handler function  
3. Register it in `main.go`

Template:
```go
type newToolArgs struct {
    Param string `json:"param" jsonschema:"description=What to do,required"`
}

type newToolResult struct {
    Result string `json:"result"`
}

func executeNewTool(ctx context.Context, args newToolArgs) (newToolResult, error) {
    // your logic here
    return newToolResult{Result: "done"}, nil
}
```

### Tuning the LLM
Edit `main.go`:
```go
genConfig := model.GenerationConfig{
    MaxTokens:   intPtr(2000),
    Temperature: floatPtr(0.7),
    Stream:      true,
}
```

## License
Apache License 2.0

## Credits
- [trpc-agent-go](https://github.com/trpc-group/trpc-agent-go) – Tencent’s open-source agent framework  
- [Google ADK Web](https://github.com/google/adk-web) – Google’s open-source agent web interface  

**NOTE:** Run only in trusted environments; tighten security before production use.