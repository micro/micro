import $axios from '@/utils/axios';

export function getRegistries(name?: string) {
    return $axios.get(`/v2/registry`);
}
