enzap
======================================================================

Getting zap.Logger out-of-the-box

How to use
----------------------------------------------------------------------

```
logger := enzap.New()
defer logger.Sync()
zap.ReplaceGlobals(logger)
zap.L().Info("hello world")
```

Stdout and Stderr
----------------------------------------------------------------------

The output destination of logs is vary according to the level.

| output | level                          |
|:-------|:-------------------------------|
| stdout | Debug / Info / Warn            |
| stderr | Error / DPanic / Panic / Fatal |


Environment
----------------------------------------------------------------------

A few environment variables changes the logger behavior.

| key             | type   | description                  |
|:----------------|:-------|:-----------------------------|
| ZAP_CALLER      | bool   | include caller               |
| ZAP_DCOLOR      | bool   | use CapitalColorLevelEncoder |
| ZAP_DEVELOPMENT | bool   | enable development mode      |
| ZAP_LEVEL       | level  | filter level                 |
| ZAP_STACK_TRACE | level  | include stack trace          |
| ZAP_TIME_LAYOUT | string | time stamp layout            |


Mode
----------------------------------------------------------------------

Production mode is default.
Defining the environment variable `ZAP_DEVELOPMENT` to `true` enables development mode.

|                         | production mode           | development mode                            |
|:------------------------|:--------------------------|:--------------------------------------------|
| ZAP_DCOLOR              | ignored                   | enables CapitalColorLevelEncoder            |
| default ZAP_LEVEL       | INFO                      | DEBUG                                       |
| default ZAP_TIME_LAYOUT | `02T15:04:05.000000Z0700` | time.RFC3339Nano                            |
| encoder                 | json encoder              | console encoder                             |
| logger.DPanic           | does not panic            | panics                                      |


Example
----------------------------------------------------------------------

```
bash-5.1$ go build -trimpath ./examples/simple
bash-5.1$ ./simple
{"level":"INFO","ts":"2021-08-08T21:41:08.185466+09:00","caller":"simple/main.go:37","msg":"info"}
{"level":"WARN","ts":"2021-08-08T21:41:08.185635+09:00","caller":"simple/main.go:39","msg":"warn"}
{"level":"ERROR","ts":"2021-08-08T21:41:08.185653+09:00","caller":"simple/main.go:41","msg":"error","stacktrace":"main.main\n\tgithub.com/takumakei/go-enzap/examples/simple/main.go:41\nruntime.main\n\truntime/proc.go:225"}
{"level":"FATAL","ts":"2021-08-08T21:41:08.185677+09:00","caller":"simple/main.go:47","msg":"fatal","stacktrace":"main.main\n\tgithub.com/takumakei/go-enzap/examples/simple/main.go:47\nruntime.main\n\truntime/proc.go:225"}
bash-5.1$ env ZAP_DEVELOPMENT=true ./simple
08T21:41:11.409757+0900 DEBUG   simple/main.go:31       level.Set       {"arg": "debug"}
08T21:41:11.409932+0900 DEBUG   simple/main.go:35       debug
08T21:41:11.409941+0900 DEBUG   simple/main.go:31       level.Set       {"arg": "info"}
08T21:41:11.409946+0900 INFO    simple/main.go:37       info
08T21:41:11.409950+0900 DEBUG   simple/main.go:31       level.Set       {"arg": "warn"}
08T21:41:11.409955+0900 WARN    simple/main.go:39       warn
08T21:41:11.409960+0900 DEBUG   simple/main.go:31       level.Set       {"arg": "error"}
08T21:41:11.409964+0900 ERROR   simple/main.go:41       error
main.main
        github.com/takumakei/go-enzap/examples/simple/main.go:41
runtime.main
        runtime/proc.go:225
08T21:41:11.409989+0900 DEBUG   simple/main.go:31       level.Set       {"arg": "fatal"}
08T21:41:11.409993+0900 FATAL   simple/main.go:47       fatal
main.main
        github.com/takumakei/go-enzap/examples/simple/main.go:47
runtime.main
        runtime/proc.go:225
bash-5.1$
```
