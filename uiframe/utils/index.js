// 判断浏览器函数
const isMobile = () => {
    if (
        window.navigator.userAgent.match(
            /(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i
        )
    ) {
        return true;
    } else {
        return false;
    }
};

// 取得时区
const getTimezone = () => {
    let offset = new Date().getTimezoneOffset();
    const o = Math.abs(offset);
    offset =
        (offset < 0 ? "+" : "-") +
        ("00" + Math.floor(o / 60)).slice(-2) +
        ":" +
        ("00" + (o % 60)).slice(-2);
    return offset;
};

// 格式化数字
const toThousands = (number = 0) => {
    let num = number.toString(),
        result = "";
    while (num.length > 3) {
        result = "," + num.slice(-3) + result;
        num = num.slice(0, num.length - 3);
    }
    if (num) {
        result = num + result;
    }
    return result;
};

const isUrl = (url) => {
    let reg =
        /^(?:(http|https|ftp):\/\/)?((?:[\w-]+\.)+[a-z0-9]+)((?:\/[^/?#]*)+)?(\?[^#]+)?(#.+)?$/i;
    return reg.exec(url);
};

export default {
    isMobile,
    getTimezone,
    toThousands,
    isUrl,
};
