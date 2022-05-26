// 加载mock 下的所有js文件，不包括子目录
const modules = require.context('~/mock', false, /\.js$/)

// 遍历文件，返回对象
const requireAll = (context) => {
    context.keys().map(context)
}

export default function ({ app }, inject) {
    requireAll(modules)
}