import $config from '~/config'

export default function ({ app }, inject) {
    inject('config', $config);
}