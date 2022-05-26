import Mock from "mockjs";

// 设置ajax跨域允许cookie
Mock.XHR.prototype.withCredentials = true;

Mock.mock('/api/users', {
    code: 0,
    msg: "OK",
    data: {
        "list|0-20": [
            {
                "id|+1": 1,
                create_time: "@cname",
                cate_name: "@cname",
                cate_en_name: "@name",
            },
        ],
        pages: {
            page_size: 20,
            page_no: 1,
            total: 50,
        },
    },
});