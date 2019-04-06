declare namespace Ajax {

    export interface AxiosResponse {
        data: AjaxResponse;
    }

    export interface AjaxResponse {
        code: number;
        success: boolean;
        data: any;
        error: Error;
    }

    export class Error {
        code: string;
        detail: string;
    }
}
