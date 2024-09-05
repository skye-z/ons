import request from './request'

function postJSON(url, data) {
    return request({
        url: url,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        data: JSON.stringify(data)
    })
}

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

export const user = {
    now: () => get('/user')
}

export const device = {
    register: name => post('/nas/register', { name }),
    rename: name => post('/nas/rename', { name }),
    remove: () => post('/nas/' + id, {}),
    getList: () => get('/nas/list'),
    getInfo: id => get('/nas/' + id)
}