import axios from 'axios';
import varibles from '../vars';
import {clearLocalStorage, getLocalStorage} from "../helpers/localStorage";

const axiosAgent = async (type, url, data) => {
    let token = getLocalStorage('session_id');

    let params = {}
    let tokenString = 'token=' + token;
    if ( url.indexOf('?') > -1 ){
        tokenString = '&' + tokenString
    }else{
        tokenString = '?' + tokenString
    }

    url = varibles.API_URL + url + tokenString;

    switch (type) {
        case 'post': {
            return axios({
                method: 'post',
                url: url,
                data: data,
                ...params
            }).catch((error) => {
                if (error.response === undefined) {
                    return Promise.reject(error);
                }
                const { status } = error.response;
                if (status === 401) {
                    window.location.href = "/";
                    clearLocalStorage('session_id');
                }

                return Promise.reject(error)
            })
        }
        case 'put': {
            return axios({
                method: 'put',
                url: url,
                data: data,
                ...params
            }).catch((error) => {
                const { status } = error.response;
                if (status === 401) {
                    window.location.href = "/";
                    clearLocalStorage('session_id')
                }

                return Promise.reject(error)
            })
        }
        case 'delete': {
            return axios({
                method: 'delete',
                url: url,
                ...params
            }).catch((error) => {
                const { status } = error.response;
                if (status === 401) {
                    window.location.href = "/";
                    clearLocalStorage('session_id')
                }

                return Promise.reject(error)
            })
        }
        default: {
            if (data) {
                params['data'] = data
            }

            return axios({
                    method: 'get',
                    url: url,
                    ...params
                }
            ).catch((error) => {
                const { status } = error.response;
                if (status === 401) {
                    window.location.href = "/";
                    clearLocalStorage('session_id')
                }
                return Promise.reject(error)
            })
        }
    }
}

export default axiosAgent
