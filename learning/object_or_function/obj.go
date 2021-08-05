package obj

import (
    "fmt"
    "errors"
)

// PayChannel 支付通道
type PayChannel int

const (
    // AliPay 支付宝支付通道
    AliPay    PayChannel = 1
    // WechatPay 微信支付通道
    WechatPay PayChannel = 2
)

// 支付宝支付方式实现
func payWithAli(price float64) error {
    fmt.Printf("Pay with Alipay %f\n", price)
    return nil
}

// 微信支付方式实现
func payWithWechat(price float64) error {
    fmt.Printf("Pay with Wechat %f\n", price)
    return nil
}

// PayWith 面向对象编程（对象为一等公民），多一种支付方式，就需要增加一个方法，并且这里要修改一次 case 
func PayWith(channel PayChannel, price float64) error {
    switch channel {
    case AliPay:
        return payWithAli(price)
    case WechatPay:
        return payWithWechat(price)
    default:
        return errors.New("not support the channel")
    }
}


// PayFunc 函数式编程（函数为一等公民），通过传递函数作为参数，return 函数调用来实现
type PayFunc func(price float64) error
// Pay is ...
func Pay(payMethod PayFunc, price float64) error {
    return payMethod(price)
}
