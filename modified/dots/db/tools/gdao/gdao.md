
- 数据库调用代码生成工具
- 针对每个表单独生成一个组件。
- 支持乐观锁，实现数据库的并发访问。

备注：
> 基于版本号的形式实现乐观锁。所以使用乐观锁请务必保证表中含有version字段，该字段将在并发检查时使用，不能含有其他实际意义。<br>
已经经过简单测试，更多具体问题bug请大家相互讨论提交issues.

> 生成dao指令执行需要项目工作在 GOPATH 之下。

```sh
# dev debug use
os.Setenv("GOPACKAGE", "model")
os.Setenv("GOFILE", "/abs/pathto/models.go")
```

```sh
# Generate model file.
go run dots/db/tools/gmodel/gmodel.go -typeName Notice

# Generate model's dao file.
# go run /abs/pathto/gdao.go -typeName Notice -daoPackage pgs
go run dots/db/tools/gdao/gdao.go -typeName Notice -daoPackage pgs
```