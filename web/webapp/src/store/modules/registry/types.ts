import {Service} from "@/store/basic/types";


export interface RegistryState {
    services: Service[];
    serviceDetail: Service[];
    pageLoading: boolean;
    serviceDetailLoading: boolean;
}