## 消息分发服务

### 路由配置规则
`message-dispatch` 使用一份 YAML 配置文件（示例：`config/development.yml`）来定义两类路由规则：

1. 接收路由（`receiver.rules`）：决定“来自 onebot 事件的消息”应该转发到哪些下游 Agent / 服务。
2. 发送路由（`sender.rules`）：决定“Agent 返回的发送动作”应通过哪个 onebot 实现 / 账号发送出去。

这两部分规则都按配置文件中的顺序自上而下匹配，命中后可根据 `is_end` 决定是否继续向下匹配。

---
#### 1. 接收路由（receiver.rules）

用于把来自不同群 / 私聊 / 用户的消息分发到对应的业务服务。

字段说明（对应 `config.go` 中 `ReceiveRule`）：

| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| name | string | 是 | "ICPC南京站大群" | 规则名称，日志追踪用，建议唯一。 |
| from_type | string | 是 | group / private / all | 消息来源类型过滤；`all` 表示不过滤类型。 |
| group_id | string | 否 (group 时建议) | "1070161797" / `.*` | 群号正则，仅当消息为群聊时参与匹配；支持正则。 |
| user_id | string | 否 | `.*` / `2731190200` | 发送者用户 ID 正则；`.*` 表示任意用户。 |
| to_addr | string | 是 | http://icpc-nanjing-agent:3000/message | 目标下游服务基础地址。 |
| is_end | bool | 是 | true / false | 命中后是否终止继续匹配；false 允许命中后继续匹配后续规则（实现多播）。 |

匹配逻辑（见 `receive.go:IsReceiveMatchRule`）：
1. 检查 `from_type`：`from_type` 需等于消息类型（`group`/`private`），或配置为 `all`。
2. 若消息为 `group`，使用正则匹配 `group_id`（配置为空或不匹配即失败）。
3. 使用正则匹配 `user_id`。
4. 全部通过则视为命中；根据 `is_end` 决定是否停止。

分发行为：命中规则后调用 `to_addr` 进行 HTTP POST，Body 为原始消息 JSON（未修改）。

示例：
```yaml
receiver:
  rules:
    - name: "测试用群"
      from_type: group
      group_id: "515785850"
      user_id: ".*"
      to_addr: http://test-agent:8000/message
      is_end: true
    - name: "ICPC南京站大群"
      from_type: group
      group_id: "1070161797"
      user_id: ".*"
      to_addr: http://icpc-nanjing-agent:3000/message
      is_end: true
    - name: "ICPC南京站技术群"
      from_type: group
      group_id: "217221751"
      user_id: ".*"
      to_addr: http://icpc-nanjing-agent:3000/message
      is_end: false   # 允许继续匹配后续规则，实现同一消息多播
```

多播场景：若某条规则 `is_end=false`，且后续规则也命中，则消息会被转发多次（每次独立 HTTP 请求）。请确保下游幂等或可接受重复处理。

---
#### 2. 发送路由（sender.rules）

用于把 Agent 产生的“发送消息动作”转交给指定的 onebot 实现（某个账号）。

字段说明（对应 `config.go` 中 `SendRule` & 处理逻辑见 `send.go:IsSendMatchRule`）：

| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| name | string | 是 | "aibot的QQ" | 规则名称。 |
| onebot_type | string | 是 | lagrange | onebot 实现类型；当前代码中仅对 `lagrange` 做了发送支持。 |
| self_id | string | 是 | "2731190200" | 机器人账号（消息发送来源身份）。匹配 `msg.params.self_id`。 |
| to_addr | string | 是 | http://lagrange-onebot:8000 | 目标 onebot 服务基础地址。 |
| is_end | bool | 是 | true / false | 命中后是否停止继续匹配。当前发送逻辑一般配置为 true。 |

匹配逻辑：
1. 仅当消息 Action 为 `send_message`（`HandleSendDispatchMessage` 中硬编码）。
2. 逐条规则比较 `self_id == msg.params.self_id`，相等即命中。
3. 命中后调用具体 Sender（目前实现：`LagrangeSender`）。若 `onebot_type` 非支持类型则忽略并输出警告。

示例：
```yaml
sender:
  rules:
    - name: "aibot的QQ"
      onebot_type: lagrange
      self_id: "2731190200"
      to_addr: http://lagrange-onebot:8000
      is_end: true
```

---
#### 3. 匹配顺序与执行策略
* 两类规则互不影响：接收链路与发送链路分别匹配。
* 都是“顺序扫描 + 可短路终止”模型：自上而下，命中后若 `is_end=true` 立即停止；否则继续尝试后续规则。
* 匹配计数会在日志中输出，方便调试（`匹配到规则数量: n`）。
* 未命中时会输出提示（接收侧：数量 0；发送侧：打印“未匹配到任何规则”）。

### 消息接受服务说明

`message-dispatch` （下称 Dispatch）默认监听 `8080` 端口（若使用容器编排请在部署层暴露 / 映射），提供以下 HTTP 接口：

| 方法 | 路径 | 用途 | 说明 |
|------|------|------|------|
| GET | `/ping` | 探活 | 返回 `{ "message": "pong" }`，可用于健康检查。 |
| POST | `/` | 接收上游 OneBot 事件 | OneBot 实现（如 lagrange.onebot）将原始事件 JSON 推送到此入口；匹配 `receiver.rules` 后转发到下游 Agent。 |
| POST | `/send` | 接收下游 Agent 的发送指令 | Agent / 工具将标准发送动作 JSON 推入；匹配 `sender.rules` 后调用对应 OneBot Adapter 发送真实消息。 |

#### 1. POST `/` （接收上游事件）
上游（onebot 实现）向该接口推送原始事件。当前代码期望字段与 `types.ReveiceMessage` 对应。核心字段：

| 字段 | 类型 | 示例 | 说明 |
|------|------|------|------|
| message_type | string | `group` / `private` | 消息来源类型（用于接收规则匹配 `from_type`）。 |
| sub_type | string | - | 子类型（按需扩展）。 |
| message_id | int64 | 123456 | 上游消息 ID。 |
| group_id | int64 | 1070161797 | 群聊场景存在；用于接收规则的 `group_id` 正则匹配。 |
| user_id | int64 | 2731190200 | 发送者用户；用于接收规则的 `user_id` 正则匹配。 |
| message | array | 见下 | 结构化消息片段数组（每段含 `type`、`data.text`）。 |
| raw_message | string | "你好" | 原始纯文本（如果有）。 |
| sender.user_id | int64 | 2731190200 | 发送者 ID（冗余）。 |
| sender.nickname | string | "Alice" | 发送者昵称。 |
| time | int64 | 1712345678 | 时间戳。 |
| post_type | string | `message` | 上游事件类型（可拓展）。 |
| self_id | string/number | 2731190200 | 机器人自身账号（支持字符串或数字，内部通过 `StringOrNumber` 兼容）。 |

`message` 字段元素示例：
```json
{
  "type": "text",
  "data": { "text": "Hello World" }
}
```

完整上游事件示例：
```json
{
  "message_type": "group",
  "sub_type": "normal",
  "message_id": 123456,
  "group_id": 1070161797,
  "user_id": 2731190200,
  "message": [
    {"type": "text", "data": {"text": "你好，Bot"}}
  ],
  "raw_message": "你好，Bot",
  "sender": {"user_id": 2731190200, "nickname": "Alice"},
  "time": 1712345678,
  "post_type": "message",
  "self_id": 2731190200
}
```

处理流程：
1. Dispatch 读取 Body → 反序列化为 `ReveiceMessage`。
2. 按 `receiver.rules` 顺序匹配（规则详见上文）。
3. 每命中一条规则：向该规则的 `to_addr` 做一次 `POST`，Body 为原始完整 JSON（不做修改）。
4. 返回 `{ "message": "pong" }` 给上游（无论是否命中规则，除非解析失败）。

下游 Agent 需实现：`POST`（或你自定义路由，但需与配置里 `to_addr` 一致）来接收这份 JSON。其内部可：
* 抽取文本拼接上下文 → 调用 LLM
* 调用工具 / MCP
* 最终决定是否回复 → 构造发送动作（下一节格式）调用 Dispatch `/send`

#### 2. POST `/send` （接收下游发送指令）
下游 Agent 向该接口推送“我要发送一条消息”的结构，字段映射 `types.SendMessage`：

| 字段 | 层级 | 类型 | 示例 | 说明 |
|------|------|------|------|------|
| action | 根 | string | "send_message" | 目前仅处理此值；否则直接忽略。 |
| params.detail_type | params | string | `group` / `private` | 指定发送目标类型。 |
| params.user_id | params | string | "2731190200" | 私聊时必填，群聊可为空。 |
| params.group_id | params | string | "1070161797" | 群聊时必填，私聊可为空。 |
| params.message | params | array | 见示例 | 消息分段（与接收侧结构兼容，可放 text / 其它类型）。 |
| params.self_id | params | string | "2731190200" | 选择哪一个机器人账号发送，用于 `sender.rules` 匹配。 |

发送动作示例（群消息）：
```json
{
  "action": "send_message",
  "params": {
    "detail_type": "group",
    "group_id": "1070161797",
    "user_id": "",
    "self_id": "2731190200",
    "message": [
      {"type": "text", "data": {"text": "收到，正在处理"}}
    ]
  }
}
```

处理流程：
1. Dispatch 读取 Body → 反序列化为 `SendMessage`。
2. 检查 `action == send_message`，否则直接返回 `pong` 不处理。
3. 按 `sender.rules` 顺序匹配 `self_id`。
4. 命中后调用具体 OneBot Sender（当前实现：Lagrange）。
5. 响应 `{ "message": "pong" }`（不包含发送结果细节，未来可扩展返回 message_id）。

#### 3. 交互时序概览
```
OneBot(Lagrange) --> POST /        (原始事件)
Dispatch         --> POST to Agent (/send on agent side with same JSON)
Agent            --> 解析/决策/生成回复结构
Agent            --> POST /send    (发送动作 SendMessage)
Dispatch         --> 规则匹配 sender.rules -> 调用 OneBot Adapter -> 平台真实发送
```

#### 4. 常见注意点

* 请求发送消息时，必须携带`self_id` 且必须与发送路由中配置的 `self_id` 一致，否则不会发送。 
* 未命中任何发送规则仅打印日志，不返回错误。 
* 目前的路由规则仅对groupid和user_id做了正则匹配，因不同平台消息结构可能不一致，因此例如是否被@提及等规则暂不考虑做适配。
