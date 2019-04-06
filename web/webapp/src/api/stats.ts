import $axios from '@/utils/axios';

export function getStats(name: string, address: string) {
    return $axios.get(`/v1/stats?service=${name}&address=${address}`);
}

export function getAPIStats(name: string, address: string) {
    return $axios.get(`/v1/api-stats?name=${name}&address=${address}`);
}
