Ниже формат `SKILL.md` для Nexflow, совместимый с Moltbot/Clawdbot и Claude Agent Skills.[^1][^2]

## Общая структура файла

`SKILL.md` всегда состоит из двух частей:[^3][^4]

1) YAML‑frontmatter сверху (метаданные, конфигурация).
2) Markdown‑тело с инструкциями и примерами.
```markdown
---
name: shell-run
description: Safely run shell commands on the local machine and return their output.
emoji: 🐚
version: 1.0.0
author: your-name
homepage: https://github.com/youruser/nexflow-skills
location: ./run.sh
tags: [system, shell, cli]
category: system
permissions: [shell, filesystem]
env_required: false
metadata: {"timeoutSec": 30, "maxOutputKb": 64}
requirements:
  binaries: [bash]
  files: [./run.sh]
  env: []
---

# Shell Run Skill

## Purpose

This skill lets you safely execute **simple** shell commands on the local Nexflow host and capture their stdout/stderr.
It is designed for short, non-interactive commands like listing files, checking disk usage, or running existing scripts.

## When to use this skill

- User explicitly asks to "run a shell command", "execute in terminal", "list files", "check disk space", etc.
- You need fresh information from the system that is not already in memory.
- The task can be completed by a single non-interactive command.

Do **not** use this skill for:
- Long-running background jobs.
- Commands that can destroy data (`rm -rf`, destructive database commands, etc.).
- Anything requiring interactive input.

## How to use this skill

1. Restate the user's goal in your own words.
2. Propose a safe shell command that accomplishes the goal.
3. Ask the user for confirmation if the command modifies files or services.
4. Once confirmed, call this skill with:
   - `command`: the exact shell command to run.
   - `cwd` (optional): directory to run in; default is the Nexflow workspace.

### Input schema

- `command` (string, required): Shell command to execute.
- `cwd` (string, optional): Working directory; if omitted, use the default workspace.

### Output schema

- `exit_code` (integer): Process exit code.
- `stdout` (string): Captured standard output (truncated if too long).
- `stderr` (string): Captured standard error (truncated if too long).

## Examples

### Example 1: List project files

User: "Покажи, какие файлы есть в текущем проекте."

Good plan:
- Run `ls -la` in the project directory to inspect the file structure.

Call this skill with:

```json
{
  "command": "ls -la",
  "cwd": "{baseDir}"
}
```


### Example 2: Check disk space

User: "Сколько места свободно на диске?"

Call this skill with:

```json
{
  "command": "df -h",
  "cwd": "{baseDir}"
}
```


## Implementation details (for the runtime)

- The Nexflow runtime executes `location` (`./run.sh`) with `command` and `cwd` passed as environment variables or arguments.
- The script must:
    - Run the command with a timeout of `metadata.timeoutSec` seconds.
    - Capture stdout/stderr and exit code.
    - Return JSON matching the output schema.
- Use `{baseDir}` placeholder to refer to the skill folder path at runtime.

```

## Обязательные поля frontmatter

Минимальный набор (по AgentSkills/Moltbot‑спеку):[^4][^1]
- `name`: уникальное имя skill (ключ в реестре).  
- `description`: короткое описание, по которому агент матчится и решает, подходит ли skill.[^3]

Пример минимального frontmatter:

```markdown
---
name: summarize-notes
description: Summarize long markdown notes into concise bullet points.
---
```


## Рекомендуемые поля frontmatter

Эти поля сильно улучшают работу Nexflow и совместимость с Moltbot/Claude:[^2][^1][^3]

- `emoji`: иконка для UI.
- `version`: семвер‑версия skill.
- `author`: автор/ник.
- `homepage`: ссылка на репозиторий/доку.
- `location`: путь к основному скрипту/бинарю внутри папки skill (`./run.sh`, `./main.py`, `./skill.bin`).
- `tags`: список тегов (например, `[system, shell, cli]`).
- `category`: логическая категория (`system`, `dev`, `productivity`, `search`, и т.п.).
- `permissions`: список прав, которые требует skill (`shell`, `filesystem`, `network`, `secrets`).[^5][^6]
- `env_required`: флаг, нужны ли обязательные переменные окружения (например, API‑ключ).
- `metadata`: однострочный JSON‑объект с произвольными настройками (таймауты, лимиты и т.п.).[^1]
- `requirements`: вложенный объект:
    - `binaries`: список бинарей, которые должны быть доступны (например, `["bash", "python"]`).
    - `files`: обязательные файлы внутри skill‑директории.
    - `env`: список переменных окружения, которые желательно заданы (например, `["OPENAI_API_KEY"]`).[^5][^1]

Важно: как в Moltbot, ключ `metadata` должен быть **одной строкой JSON**, без переноса.[^1]

## Markdown‑тело: что обязательно

Тело `SKILL.md` — это “how‑to” для агента, а не для пользователя, и должно включать:[^4][^3]

1. **Purpose**: 1–2 абзаца, что делает skill и какую проблему решает.
2. **When to use**: буллеты “когда этот skill уместен” и “когда нельзя/не нужно использовать”.[^3]
3. **How to use**: пошаговая схема действий агента (переформулировать цель, предложить план, спросить подтверждение и т.п.).
4. **Input schema**: перечисление параметров, типов и обязательности.
5. **Output schema**: что вернёт скрипт, какие поля и как агент должен их использовать.
6. **Examples**: 2–3 реальных примера запросов и JSON‑вызовов (как в примере выше).[^7][^3]
7. **Implementation details (optional)**: подсказки для рантайма (таймауты, особенности запуска, ограничения).

## Правила для Nexflow‑рантайма

Чтобы этот формат реально работал:[^2][^1]

- При старте Nexflow:
    - Найти все `SKILL.md` в configured директориях.
    - Прочитать frontmatter, построить реестр skills (name → location/tags/permissions/metadata).
    - В системный промпт подмешивать только name/description/location (как делает Moltbot/AgentSkills, чтобы не раздувать контекст).[^2][^1]
- При выборе skill:
    - Агент сначала матчится по `name`/`description`/tags (discovery), без загрузки всего тела.[^2]
    - После решения “я хочу использовать этот skill” Nexflow загружает целиком `SKILL.md` и даёт его агенту (activation).[^2]
- При исполнении:
    - Рантайм использует `location`, `{baseDir}`, `metadata`, `permissions` и `requirements` для запуска скрипта.
    - Результат скрипта должен соответствовать описанному Output schema; если нет — агент получает ошибку и описание проблемы (для self‑healing).[^7][^3]

Если хочешь, дальше могу сделать ещё один пример `SKILL.md` для чисто HTTP‑интеграции (например, “GitHub Issues Skill”) в этом же формате.
<span style="display:none">[^10][^11][^12][^13][^14][^15][^16][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://docs.molt.bot/tools/skills

[^2]: https://claudecn.com/en/blog/claude-skills-landing-guide/

[^3]: https://zenvanriel.nl/ai-engineer-blog/moltbot-custom-skill-creation-guide/

[^4]: https://www.aifreeapi.com/en/posts/how-to-create-clawdbot-skills

[^5]: https://composio.dev/blog/secure-moltbot-clawdbot-setup-composio

[^6]: https://socprime.com/active-threats/the-moltbot-clawdbots-epidemic/

[^7]: https://www.datacamp.com/tutorial/moltbot-clawdbot-tutorial

[^8]: https://www.reddit.com/r/LocalLLM/comments/1qq32b4/resource_500_formatted_skills_for_moltbotclawdbot/

[^9]: https://www.youtube.com/watch?v=Fh-aBKrG5CI

[^10]: https://www.datacamp.com/de/tutorial/moltbot-clawdbot-tutorial

[^11]: https://www.digitalocean.com/community/tutorials/moltbot-quickstart-guide

[^12]: https://docs.clawd.bot/tools/skills

[^13]: https://github.com/clawdbot/clawdhub

[^14]: https://www.youtube.com/watch?v=mDsyFrQPPfg

[^15]: https://vertu.com/lifestyle/complete-clawdbot-tutorial-deploy-with-caution/

[^16]: https://www.linkedin.com/posts/juliangoldieseo_moltbotclawdbot-ai-seo-how-i-ranked-1-activity-7422305673054535681-h-95

