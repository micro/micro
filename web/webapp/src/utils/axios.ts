/* eslint-disable */
import axios from "axios";
import config from "@/config";
import {Error} from '@/store/basic/types'

const baseURL = config.url.basicUrl;
const $axios = axios.create({
    baseURL,
    withCredentials: true,
    timeout: 10000
});

$axios.interceptors.request.use(
    (config: any) => {
        return config;
    },
    (error: any) => {
        if (error.error.code) {

            console.log(error.error.code);
        }
        return Promise.reject(error);
    }
);

$axios.interceptors.response.use(
    (response: any) => {
        // store.commit(TYPES.SET_USER_LOADING_GET_DONE);

        if (
            response.status >= 200 &&
            response.status < 300 &&
            response.data.success == true || response.status == 200
        ) {
            if (response.data) {
                return response.data;
            } else {
                return response;
            }
        } else if (response.data && !response.data.success) {
            const error = new Error(response.data.error);
            // throw error;
            return error;
        } else {
            const error = new Error(response.statusText);
            // throw error;
            return error;
        }
    },
    (error: any) => {

        if (
            error.code === "ECONNABORTED" &&
            error.message.indexOf("timeout") != -1
        ) {
            const config = error.config;
            // If config does not exist or the retry option is not set, reject
            if (!config || !config.retry) return Promise.reject(error);

            // Set the variable for keeping track of the retry count
            config.__retryCount = config.__retryCount || 1;

            // Check if we've maxed out the total number of retries
            if (config.__retryCount >= config.retry) {
                return Promise.reject(error);
            }

            // Increase the retry count
            config.__retryCount += 1;

            // Create new promise to handle exponential backoff
            let backoff = new Promise(resolve => {
                setTimeout(() => {
                    // console.log("resolve");
                    resolve();
                }, config.retryDelay || 1);
            });

            return backoff.then(() => $axios(config));
        } else {

            if (error.response.data) {
                return {
                    success: false,
                    error: new Error(error.response.data.code, error.response.data.detail)
                }
            }

            if (error.code) {
                return {
                    success: false,
                    error: new Error(error.code, error.detail)
                }
            }

            if (error.response.status >= 499) {
                return {
                    success: false,
                    error: new Error(error.response.status, error.response.statusText)
                }
            }

            return Promise.reject(error);
        }
    }
);

export default $axios;
