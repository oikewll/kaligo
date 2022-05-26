const $config = {}

$config.client = require('./client');

/**
 * client端不能引入server端的配置，防止关键数据泄露
 */
if(process.server){
	$config.server = require('./server');
}

export default $config;