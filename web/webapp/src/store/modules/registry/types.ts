import {Service} from "@/store/basic/types";


export interface RegistryState {
    services: Service[];
    webServices: Service[];
    serviceDetail: Service[];
    pageLoading: boolean;
    serviceDetailLoading: boolean;
}