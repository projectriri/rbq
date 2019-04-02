# RBQ

调试用绒布球。

你需要放置 `goshujinsama.json` 在绒布球的嘴里。

## 命令列表

| 命令 | 参数 | 说明 |
| --- | --- | --- |
| !!rbq::say | ... | 让绒布球说话 |
| !!rbq::who | | 是不是主人呢 |
| !!rbq::listen | 正则表达式，空参数相当于匹配任意字符串 | 让绒布球把当前聊天中的内容转发给自己，只有主人才可以用 |
| !!rbq::leave | | 取消 listen |
| !!rbq::scdo | [ChatID [ChatType [Messenger]]] -- 内容 | 跨聊天给梨梨发送消息，只有主人和主人允许的人才可以用 |
| !!rbq::ntr | [时长，默认为 "15m"] | 允许当前聊天里的人类一段时间内使用绒布球 |
| !!rbq::use | | 热更新 `goshujinsama.json` |
