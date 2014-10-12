## Epooll v1.0.0 (2014-10-01)

epooll is an open-source, high-performance, modularity, full-stack web framework.

More info [epooll.com](http://www.epooll.com)

* [Changelog](https://github.com/owner888/epooll/blob/master/README.md#changelog)
* [Installing](https://github.com/owner888/epooll/blob/master/README.md#installing)
* [Testing](https://github.com/owner888/epooll/blob/master/README.md#testing)
* [Examples](https://github.com/owner888/epooll/blob/master/README.md#examples)
* [To do](https://github.com/owner888/epooll/blob/master/README.md#to-do)
* [Known bugs](https://github.com/owner888/epooll/blob/master/README.md#known-bugs)
* [Documentation](https://github.com/owner888/epooll/blob/master/README.md#documentation)

## Changelog

v1.0.0: 初始化类库

## Installation

    $ go get github.com/owner888/epooll

## Testing


## Examples

### Example 1

`
├── conf
│   └── app.ini
├── control
│   ├── ctl_demo.go
│   ├── ctl_index.go
├── data
│   ├── cache
│   └── log
├── main.go
├── model
│   ├── mod_common.go
├── static
│   ├── css
│   ├── images
│   └── js
└── template
    └── index.tpl
`

## To do

1. http header处理、cookie、session、isAjax
2. 数据库分页
3. 测试cpu、内存使用率
4. 完善使用文档
5. 基于此框架开发一个CMS
6. 一个服务自定义绑定多个端口

## Known bugs


## Documentation
* [control](http://www.godoc.org/pkg/github.com/owner888/epooll/control)
* [db](http://www.godoc.org/pkg/github.com/owner888/epooll/db)
* [redis](http://www.godoc.org/pkg/github.com/owner888/epooll/redis)

## Contact US
QQ:525773145

## LICENSE

epooll is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
