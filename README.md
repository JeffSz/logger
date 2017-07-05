# example 1
```
// Write to log.log
logger.Info("This is line 1 ", "and tail")
```

# example 2
```
// Write to test.log
logger.BasicConfig("test.log", logger.INFO, "[{datetime}] {message}", "2006-01-02")
logger.Info("This is line 1", "and tail")
```

# example 3
```
// Write to rotate.log rotate
handler, _ := handlers.NewTimeRotateHandler(handlers.Minute, logger.INFO, "test1.log", 3)
logger.Logger.AddHandler(handler)
logger.Logger.SetFormat("[{level}] {message}")
logger.Info("This is line 1", "and tail")
```
