## 接口封装说明

#### 约定
> 目录下的所有js文件都会读取挂载到this.$api，文件名为$api调用的键名

用法： `this.$api.{NAMESPACE}.{MODULE}()`


## 规范

- api目录下的 {NANME}.js 即为api的命名空间，以业务划分模块(module)
- 所有封装的接口，都有默认参数
- 返回Promise

```javascript
// app为Vue 实例
export default app => ({
	login: () => app.$axios.get('/login')
})

this.$api.common.login()
```