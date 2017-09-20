Go语言如何在没有实现功能的情况下写出完善的单元测试代码
==============================================

## 背景

最近在研究用Go写一个自己的解释型语言，有一本书叫《Writing An Interpreter In Go》, 作者在讲解如何编写解释器的时候，都是从写一个`_test.go`开始的，也就是说作者习惯于先写单元测试，以测试驱动开发，其实这是一个非常好的习惯，不过，作者在写`_test.go`文件的时候，都是先假设这个结构体、函数已经存在了，并且没有把关键的对象抽象成接口，因此，作者在运行`go test`的时候，是无法完成测试的，因为连编译都过不了，必须一边完善代码，一边重复运行`go test`，一直到完成开发。

基于这种开发模式下，其实我更期望能有一个Mock实现，写测试代码的时候畅通无阻，即使是没有实现，也能把各个测试用例覆盖到，当真实的实现完成后，我们只需要把mock实现替换成真实的实现就好了。

这么做还带来另一个好处，如果公司有SDET岗位，则可以直接让测试人员编写单元测试，开发任务和测试任务可以并行。

---

## gomock 框架

昨天闲着没事逛了逛 https://github.com/golang, 发现了一个非常有意思的框架: [gomock](https://github.com/golang/mock), 官方的描述是，这是一个mocking framework, 在使用上也很简单，大致的步骤如下：

1、定义一个待实现的接口.

```go
  type MyInterface interface {
    SomeMethod(x int64, y string)
    GetSomething() string
  }
```

2、使用`mockgen`生成mock代码.

3、测试:

```
func TestMyThing(t *testing.T) {
    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()
    
    mockObj := something.NewMockMyInterface(mockCtrl)
    mockObj.EXPECT().SomeMethod(4, "blah")
    mockObj.EXPECT().GetSomething.Return("haha")
}
```
 
看了这个步骤，想必大家应该猜到了，mocking framework 使用的是根据你的接口定义来自动生成一个Mock实现，我们还可以往这个实现里注入数据。

### 期望接口调用参数

我们可以通过 `.EXPECT()` 为这个mock对象注入期望值

```go
mockObj.EXPECT().SomeMethod(4, "blah")
```

- 测试通过:

```go
mockObj.SomeMethod(4, "blah")
```

- 测试不通过:

```go
mockObj.SomeMethod(5, "bldah") // 此处调用时直接会抛出错误
```

对于参数类型的期望，在调用这个Mock函数的时候会直接抛错异常


### 期望返回值

`.EXPECT()` 同样可以为某个函数注入返回值


```go
mockObj.EXPECT().GetSomething().Return("haha")
```

- 测试通过:

```go
if "haha" == mockObj.GetSomething() {
    // -> 执行到这里
    // ...
}
// ...
```

- 测试不通过:

```go
if "haha" == mockObj.GetSomething() {
    // -> 不会执行到这里
    // ...
}
// ...
```


功能强的还不止这个，如果我们在测试一个循环，希望的是每次调用 `GetSomething()` 都返回不同的值，该怎么办？

答案很简单，依次调用

```go
gomock.InOrder(
    mockObj.EXPECT().GetSomething().Return("A"),
    mockObj.EXPECT().GetSomething().Return("B"),
    mockObj.EXPECT().GetSomething().Return("C"),
)
```

接下来，让我们来实战一下吧。

---

## gomock 实战

我们以《Writing An Interpreter In Go》这本书中的 `monkey` 语言的 `lexer` 作为例子

我们看一下 `monkey` 的目录结构：

```bash
monkey> tree .
.
├── lexer
│   ├── lexer.go
│   ├── lexer_test.go
│   └── mock_lexer
│       └── mock_lexer.go
└── token
    └── token.go
```

### Lexer和Token

```
type Lexer interface {
	NextToken() token.Token
}
```

lexer的功能很简单，每次调用`NextToken()`，都是返回下一个`Token`
`Token`的结构
```
type TokenType string

type Token struct {
	Type    TokenType   // 类型
	Literal string      // 内容
}
```

比如下面的`go`语句

```go
var a = 1
```

lexter 在调用三次`NextToken()`后会得到三个Token, 依次是:

```
Token{VAR, var}
Token{IDENT, a}
Token{INT, 1}
```


### 测试思路

其实测试方法就是：给定一段代码，用Lexer解析后，能得到指定顺序的Token，而`gomock`是完全可以实现的。


### 使用gomock

#### 安装gomock

```bash
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
```

#### 生成mock代码

```bash
mockgen -source lexer.go -destination  mock_lexer/mock_lexer.go
```

### 编写lexer_test.go


#### 测试数据

input: 输入的语句
tokens: 期望的Token

```go
func getTestData() (input string, tokens []token.Token) {
	input = `
let five = 5;
let ten = 10;
`

	tokens = []token.Token{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	return
}

```

#### 生成一个真实的MonkeyLexer实例

当然，这里我们没有实现，所以返回是`nil`

```go
func newMonkeyLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {
	return nil, func() {}
	// return NewMonkeyLexer(input), func() {}
}
```

#### 构建MockLexer实例

由于没写完真正的lexer, 那么我们就开始Mock吧

```go
func newMockLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {

	ctrl := gomock.NewController(t)

	// 生成一个Mock实例
	mockLexter := mock_lexer.NewMockLexer(ctrl)

    // 将期望值一次传递给 NextToken()
    // 每次调用 NextToken() 也会依次获得期望值
	for i := 0; i < len(excepts); i++ {
		mockLexter.EXPECT().NextToken().Return(excepts[i])
	}

	l = mockLexter
	
	// 用于清理
	deferFN = func() { ctrl.Finish() }

	return
}
```


为了方便在mock实例和真实实例之间进行切换，我们可以通过环境变量来控制当前的测试实例是什么，如果要使用mock进行测试，我们只需要在运行 `go test` 前执行:

```go
> export GO_MOCK_TEST=1
```

或

```go
> GO_MOCK_TEST=1 go test -v
```

```go
func newLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {
	env := os.Getenv("GO_MOCK_TEST")
	if env == "1" {
		t.Log("MOCK TEST ENABLED!!!")
		return newMockLexer(input, excepts, t)
	}

	return newMonkeyLexer(input, excepts, t)
}
```



以下是真正的测试代码，在没真实实现`monkey lexer`的情况下，我们可以写测试代码了，而且如果运行 `go test -v` 也是能通过的。

```go
func TestNextToken(t *testing.T) {

	input, excepts := getTestData()

	l, fn := newLexer(input, excepts, t) // 生成 Lexer 对象

	defer fn() // 清理

	for i, tt := range excepts {
		tok := l.NextToken()  // 获取下一个Token

		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.Type, tok.Type)
		}

		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.Literal, tok.Literal)
		}
	}
}
```

> 下载完整的代码：
https://github.com/xujinzheng/monkey


### 测试

#### 使用mock实例进行测试

```bash
monkey> cd lexer
lexer> GO_MOCK_TEST=1 go test -v

=== RUN   TestNextToken
--- PASS: TestNextToken (0.00s)
PASS
ok  	github.com/xujinzheng/monkey/lexer	0.007s
```


#### 使用真实实例测试

```go
func newMonkeyLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {
	return NewMonkeyLexer(input), func() {}
}
```

我们将测试数据修改一下，假设 `ten=666`, 但不修改期望值，让测试报错

```go
input = `
let five = 5;
let ten = 666;
`
```

再次运行测试

```bash
monkey> cd lexer
lexer> go test -v

=== RUN   TestNextToken
--- FAIL: TestNextToken (0.00s)
	lexer_test.go:52: tests[8] - literal wrong. expected="10", got="666"
FAIL
exit status 1
FAIL	github.com/xujinzheng/monkey/lexer	0.008s
```

这里就报错了，说明我们的真实实例实现得有问题，需要修复这个BUG

`gomock` 的使用到这里就结束了，除了上面介绍到的一些功能，`gomock` 还有很多其他丰富的方法，大家可以去 [GoDoc](https://godoc.org/github.com/golang/mock/gomock) 获取更详细的接口信息。





