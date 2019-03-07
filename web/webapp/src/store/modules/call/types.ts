import {Service} from "@/store/basic/types";


export interface CallState {
    services: Service[]
    requestLoading: boolean
    requestResult: Map<string, Object>
}