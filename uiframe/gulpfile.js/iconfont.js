/**
 * 会根据根目录下的.iconfont 路径远程读取css文件
 * 并把里面的字体文件下载到对应目录
 * 注意要http或者https请求头，否则会报错
 * 
 * @author Dawsonliu
 * @date 2021/11/09
 */

const request = require("request");
const fs = require('fs');
const path = require('path')
const FONT_DIR = path.resolve(__dirname, '../static/font')

module.exports = (iconsrc) => {
	const regString = (s) => {
		// let reg = /url\s*\(('\s*[A-Za-z0-9\-\_\.\/\:]+\s*')\)\s*;?/gi
		let reg = /url\s*\(('\s*[A-Za-z0-9\-\_\.\/\:]+\s*)?/gi;
		s = s.match(reg);
		return(s)
	}
	const reqEleAndFsWrite = (url, filename) => {
		if(url.indexOf('http') === -1){
            throw new Error("🙅‍♂ .iconfont字体路径需要填写主机头")
        };
		request(url).pipe(fs.createWriteStream(filename ? path.resolve(FONT_DIR, filename) : ''))
	}
	const mkdirsSync = (dirname) => {
		if (fs.existsSync(dirname)) {  
			return true;
		} else {  
			if (mkdirsSync(path.dirname(dirname))) {  
				fs.mkdirSync(dirname);  
				return true;  
			}  
		}  
	} 
	mkdirsSync(FONT_DIR);

	return request(iconsrc, (err, res, body)=>{
		let urlArr = regString(body);
		// 先去重
		urlArr = Array.from(new Set(urlArr))

		// 删掉不是url的项
		urlArr.splice(urlArr.findIndex(item => item.indexOf('//at.alicdn') === -1), 1)

		// 替换掉前面的字符提取出最终的pure url, 多线程request到本地
		urlArr.forEach(item=>{
			item = item.replace("url('", 'http:');
			console.log(item);
			let suffix = item.substring(item.lastIndexOf('.') + 1);
			reqEleAndFsWrite(item, `iconfont.${suffix}`);
		});

		const fontname = /.*?com\/t\/(.*)/.exec(urlArr[0])[1].split('.')[0];
		body = body.replace(fontname, 'iconfont').replace(/\/\/at.alicdn.com\/t/gi, '.');
		// 最后写入字体目录
		fs.writeFileSync(path.resolve(FONT_DIR, 'iconfont.css'), body)
	})
}