# Module 初始化测试
## Output
```
var init in 'm'
module init in 'm'
var init in 'n'
module init in 'n'
var init in main
module init in main
```
## 结论
1. 初始化顺序depency var -> dependency init() -> 当前 var -> 当前init()
2. 同级的依赖按照包引入顺序决定（字母排序）