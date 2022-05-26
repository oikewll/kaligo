/**
 * gulpfile一些工具，不能用于生产环境的 dependencies 依赖
 * 只能在nuxt build后执行一些文件处理，譬如ssr部署和字体文件本地化
 * @author Dawsonliu
 * @date 2021/11/09
 */

const fs = require('fs');

/**
 * 解析配置文件
 * @param {*} uri 
 */
const parseConfig = (uri) => {
    try {
        var content = fs.readFileSync(uri, 'UTF-8');
        var regexjing = /\s*(#+)/; //去除注释行的正则
        var regexkong = /\s*=\s*/; //去除=号前后的空格的正则
        var keyvalue = {}; //存储键值对

        var arr_case = null;
        var regexline = /.+/g; //匹配换行符以外的所有字符的正则
        while (arr_case = regexline.exec(content)) { //过滤掉空行
            if (!regexjing.test(arr_case)) { //去除注释行
                keyvalue[arr_case.toString().split(regexkong)[0]] = arr_case.toString().split(regexkong)[1]; //存储键值对
                console.log('parse uri ====>', arr_case.toString());
            }
        }
    } catch (e) {
        throw e;
    }
    return keyvalue;
}

/**
 * 检查路径是否存在 如果不存在则创建路径
 * @param {string} folderpath 文件路径
 */
const createPath = (folderpath) => {
    const pathArr = folderpath.split('/');
    let _path = '.';
    for (let i = 0; i < pathArr.length; i++) {
        if (pathArr[i]) {
            _path += `/${pathArr[i]}`;
            if (!fs.existsSync(_path)) {
                fs.mkdirSync(_path);
            }
        }
    }
}

module.exports = {
    parseConfig,
    createPath
}