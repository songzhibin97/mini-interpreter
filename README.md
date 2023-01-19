# Mini-interpreter

一个mini解释器,具有词法分析、语法分析、求值等简单功能.
解析表达式使用的是递归下降普特拉解法

```
.
├── README.md
├── ast // 抽象语法树定义
│   ├── ast.go 
│   ├── ast_test.go
│   ├── modify.go
│   └── modify_test.go
├── eval // 解析表达式
│   ├── builtins.go
│   ├── eval.go
│   └── eval_test.go
├── go.mod
├── go.sum
├── lexer // 词法解析器
│   ├── lexer.go
│   └── lexer_test.go
├── main.go
├── object // 抽象对象类型
│   ├── env.go
│   ├── object.go
│   └── object_test.go
├── parser // 词法分析器
│   ├── parse.go
│   └── parse_test.go
├── repl
│   └── repl.go
└── token // 词法单元
    └── token.go

```

## Demo



```

```shell
# go run main.go

Welcome to Mini-interpreter
>>>print("hello world")
hello world
>>>
```


```shell
>>>var i = 1
>>>var b = true
>>>var s = "string"
>>>var array = [1,2,3]
>>>var mp = {1:1,2:2}
>>>print(i,b,s,array,mp)
1
true
string
[1, 2, 3]
{1:1, 2:2}

>>>print(array[1])
2

>>>print(mp[2])
2

>>>func add(a,b){return a + b}
>>>print(add(1,2))
3

>>>print( 1 + 2 * 3 )
7

>>>func max(a, b) {if (a > b) { return a } return b }
>>>print(max(1,2))
2
>>>print(max(2,1))
2

>>>func max(a, b) {if (a > b) { return a } else { return b }}
>>>print(max(1,2))
2
>>>print(max(2,1))
2

>>>print(!true)
false

>>>print(-10)
-10
```


