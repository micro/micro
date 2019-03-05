declare namespace Ajax {
  /**
   * axios 返回数据
   * @export
   * @interface AxiosResponse
   */
  export interface AxiosResponse {
    data: AjaxResponse;
  }

  /**
   * 请求接口数据
   * @export
   * @interface AjaxResponse
   */
  export interface AjaxResponse {
    code: number;
    /**
     * 状态码
     * @type {number}
     */
    success: boolean;

    /**
     * 数据
     * @type {any}
     */
    data: any;

    /**
     * 错误
     * @type {any}
     */
    error: Error;

    /**
     * 消息
     * @type {string}
     */
    message?: string;
  }

  /**
   * 请求接口数据 错误信息
   * @export
   * @class Error
   */
  export class Error {
    code: string;
    detail: string;
  }
}
