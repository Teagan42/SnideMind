![coverage](https://raw.githubusercontent.com/teagan42/snidemind/badges/.badges/main/coverage.svg)

# 🧠 SnideMind

Runs your AI pipeline. Judges your life choices while doing it.


## 🧬 What Is This Monstrosity?

SnideMind is a highly modular, fully self-hosted AI orchestration system. It exists to automate everything you were too lazy or emotionally broken to do yourself. It integrates local LLMs, handles prompt pipelines, calls memory resources, categorizes inputs, and probably knows your secrets.

It’s like if ChatGPT went to therapy, got more sarcastic, and then decided to run your home.


## 🛠️ Features  
### 🧩 Plugin-based pipeline engine  
Dynamically assembles AI workflows based on configuration. No hardcoding. Unless you’re into that.  
### 🧠 LLM Abstraction Layer  
Supports local LLMs (Ollama, OpenRouter, etc.) — just don’t expect hand-holding.  
### 🗃️ Memory + Tool Routing  
Categorizes prompts, triggers tools, stores notes — and unlike you, it remembers.  
### 🧾 JSON/YAML Configurable  
Define behavior in config files so you can break things predictably.  
### 📡 Docker Swarm Ready  
Because you could just use Docker Compose, but you want pain.  
### 🎯 Observability Hooks  
Logs, metrics, and just enough structured output to guilt you into fixing things.  


## 📦 Quick Start (for masochists)

```shell
git clone https://github.com/Teagan42/SnideMind.git
cd SnideMind
docker stack deploy -c docker-stack.yml snidemind
```

Now go to http://localhost:3000 and bask in its passive-aggressive glory.


## 🧞 Configuration

SnideMind runs off a YAML file named config.yaml, which should live at the root of your server or be specified via CLI:

```shell
./snidemind --config path/to/config.yaml
```

Example config.yaml
```yaml
server:
  bind: 0.0.0.0
  port: 3000

mcp_servers:
  - name: "home_mcp"
    url: "http://localhost:9001"
    type: "sse"

pipeline:
  steps:
    - type: extractTags
    - type: retrieveMemory
    - type: llm
      llm:
        model: "mistral"
        base_url: "http://localhost:11434"
    - type: storeMemory
```

You can mix, match, fork, and combine these steps like a modular disaster sandwich.

## 🤖 Pipeline Engine

Each step in your pipeline is defined by its type. Supported types include:  

* extractTags
* fork
* llm
* reduceTools
* retrieveMemory
* storeMemory

Yes, you can nest forks.

Just because recursion didn’t kill you yet doesn’t mean it won’t.

## 📚 Documentation

Coming soon, maybe...  
In the meantime, just read the code like a real dev.  

## 🤬 Philosophy

SnideMind doesn’t coddle. It doesn’t sugarcoat. It works — and it’s going to remind you that it works better than you.

## 🧟‍♀️ License

MIT. Because even sarcasm has boundaries.

## ✨ Contributing

PRs welcome. If your code sucks, we’ll make fun of it in the next release.

## 🐍 Like This?

Also check out SarcastiStack:  

Runs your life. Reminds you it’s doing it better than you ever did.  
