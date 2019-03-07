import $axios from '@/utils/axios';

export function getServiceDetails() {
    return $axios.get(`/v1/service-details`);
}