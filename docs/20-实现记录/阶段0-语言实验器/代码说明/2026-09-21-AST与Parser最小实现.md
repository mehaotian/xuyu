# AST 与 Parser 最小实现说明

## 这一步解决什么

Lexer 已经可以产生 Token，但 Token 还不能表达“这是一条赋值语句”。这一
步把固定顺序的 Token 组成 AST，让源码第一次拥有结构。

## 当前支持的唯一形式

```text
设置 名字 = 1
```

它会变成：

```text
程序[
    设置(名字, 整数(1))
]
```

## 结构边界

- `Program` 表示一个源码文件；
- `Assignment` 表示一条赋值语句；
- `IntegerLiteral` 表示整数值；
- `Statement` 和 `Expression` 接口先保留边界，暂时不加入更多节点；
- Parser 要求 Token 顺序完整，缺少或多出 Token 都返回错误。

## 没有做什么

- 没有执行赋值；
- 没有变量环境；
- 没有条件、循环、函数和对象；
- 没有读取 `.中` 文件；
- 没有把 AST 接入 CLI。
