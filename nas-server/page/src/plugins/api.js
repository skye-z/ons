import request from './request'

function post(url, data) {
    let cache = new URLSearchParams();
    for (let key in data) {
        cache.append(key, data[key]);
    }

    return request({
        url: url,
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        data: cache
    })
}

function get(url) {
    return request({
        url: url,
        method: 'GET'
    })
}

export const setting = {
    all: () => get('/setting'),
    setPassword: () => post('/setting/pwd'),
}

export const device = {
    register: () => get('/register'),
    getState: () => get('/conn/state'),
    openServer: () => get('/conn/open'),
    closeServer: () => get('/conn/close'),
    switchAuto: () => get('/auto/switch'),
}