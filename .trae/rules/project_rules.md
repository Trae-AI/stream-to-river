1. 在编写文件之前，使用 `list dir` 检查文件是否存在。
2. 编辑文件之前，一定要先看文件，一定不要破坏原有的代码逻辑，不要随意修改原有的代码。
3. 确保在提供的代码中增加对应的 `import` 语句（例如 `time`，`math/rand` 等等）
4. 一定不要修改原有的 `import` 语句，如果要新增，一定用 `...existing code...` 来放原有 `import` 语句。
5. 这个目录使用 linter 检查实现。
6. `rpcservice` 的接口需要在 `rpcservice/dal` 中实现数据访问逻辑，在 `rpcservice/biz` 中实现业务逻辑，不要混淆。
7. `apiservice` 的路由需要在 `apiservice/router.go` 中实现，对应的 api 处理逻辑需要在 `apiservice/biz/handler` 中实现
8. `apiservice` 使用的 `rpcclient`，不要做任何修改，直接使用即可。