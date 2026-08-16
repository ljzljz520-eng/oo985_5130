# 工艺培训课程管理

这是一个纯 Go 的命令行示例项目，用内存仓储管理木工、陶艺、皮具和金工培训课程。管理员维护分类、审核课程并归档停课课程，老师维护材料、工具和安全提示，在审核通过后公开课程。

## 环境

- Go 1.26.6
- `GOTOOLCHAIN=local`
- 无外部服务或第三方依赖
- 支持 `CGO_ENABLED=0`

## 运行

从模块根目录运行完整业务演示：

```sh
GOTOOLCHAIN=local go run ./cmd/craftschool
```

程序会创建四个工艺分类，完成一门皮具课程的编辑、送审、审核、公开与停课归档，并输出最终状态。

## 测试

从模块根目录运行：

```sh
GOTOOLCHAIN=local go test -count=1 ./...
```

## 构建

```sh
CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...
```
