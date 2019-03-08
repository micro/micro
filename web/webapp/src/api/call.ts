import $axios from '@/utils/axios';

let $qs = require('qs');

export function getServiceDetails() {
    return $axios.get(`/v1/service-details`);
}

export function postServiceRequest(req: object) {
   // let postData = $qs.stringify(req);
    return $axios.post(`/v1/rpc`, req);
}