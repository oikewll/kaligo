### 组件目录

#### 约定 
1 目录大驼峰命名，文件小写下横线
2 公共组件放在`Common`目录下, 页面组件放在`Pages`下

#### TypeScript版本的组件开发 [可选]

```javascript
import Vue from 'vue';

export default Vue.extend({
    name: 'MY_COMPONENT',
    mouted() {
        ...具体业务操作
    },
    computed: {
        MY_VARIABLE(): 返回值 {
        ...具体业务操作
        }
    }
})
```

- 使用 Vue.extend() 包裹
- computed 中的函数需要有具体的返回值以便 TypeScript 做类型推导
- 所有挂载在Vue实例上的自定义属性如this.$api，需要在 {PROJECT_ROOT}/typing/index.d.ts 中扩展 vue 的接口模块