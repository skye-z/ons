import axios from 'axios'
import router from './router'

const request = axios.create({
    baseURL: '/api',
    timeout: 60000
})

request.interceptors.request.use(
    config => {
        let token = localStorage.getItem("access:token");
        if (token) config.headers['Authorization'] = 'Bearer '+token;
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

request.interceptors.response.use(
    response => {
        if (response.config.responseType === 'blob') {
          return response;
        }
        if (response?.data?.code){
            let code = parseInt(response.data.code);
            // if (code >= 10100 && code <= 10103) {
            //     window.$message.warning(response.data.message);
            //     localStorage.removeItem("access:token");
            //     location.href = '/app'
            // } else 
            return response.data
        } else return response.data
    }, () => {
        window.$message.error('网络异常')
        throw "网络异常";
    }
)

export default request
