/**
 * 用法：
 * this.$api.NAMESPACE.MODULE()
 */
const modules = require.context('~/apis', false, /\.js$/)
const $api = {}

// 遍历文件，返回对象
const requireAll = (context) => {
    const files = []
    context.keys().forEach((key) => {
        const filename = key.replace(/(.*\/)*([^.]+).*/gi, '$2')
        files.push({
            [filename]: context(key).default
        })
    })
    return files
}

export default function({ app }, inject) {
    // 挂载到$api
    requireAll(modules).forEach((module) => {
        Object.keys(module).forEach((key) => {
            if (typeof module[key] === 'function') {
                $api[key] = module[key](app);
            }
        })
    });

    inject('api', $api)
}