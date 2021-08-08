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

The output destination of the log changes according to the level.

| output | level                          |
|:-------|:-------------------------------|
| stdout | Debug / Info / Warn            |
| stderr | Error / DPanic / Panic / Fatal |


Environment
----------------------------------------------------------------------

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

|                         | production mode           | development mode                            |
|:------------------------|:--------------------------|:--------------------------------------------|
| ZAP_DCOLOR              | ignored                   | enalbes CapitalColorLevelEncoder            |
| default ZAP_LEVEL       | INFO                      | DEBUG                                       |
| default ZAP_TIME_LAYOUT | `02T15:04:05.000000Z0700` | time.RFC3339Nano                            |
| encoder                 | json encoder              | console encoder                             |
| logger.DPanic           | does not panic            | panics                                      |
