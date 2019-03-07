import $axios from '@/utils/axios';

export function getServices() {
    return $axios.get(`/v1/services`);
}

export function getService(name: string) {
    return $axios.get(`/v1/service/${name}`);
}