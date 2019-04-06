export class Error {
    public code ?: string
    public detail ?: string

    constructor(code ?: string, detail ?: string) {
        this.detail = detail;
        this.code = code;
    }
}

export class Page<T> {
    public pageSize = 10;
    public pageNo = 1;
    public items: T[] = [];
    public total: number = 0;
}


export class Service {
    public name ?: string;
    public version ?: string;
    public metadata ?: Map<string, string>;
    public endpoints ?: Endpoint [];
    public nodes  ?: Node  [];
}

export class Node {
    public id    ?: string
    public address ?: string
    public port  ?: number
    public metadata  ?: Map<string, string>;
}

export class Endpoint {
    public name !: string;
    public request ?: Value;
    public response ?: Value;
    public metadata ?: Map<string, string>;
}

export class Value {
    public name  ?: string
    public type  ?: string
    public values ?: Value[]
}

export default class Language {
    public flag ?: string;
    public title ?: string;
    public lan ?: string;

    constructor(flag: string, title: string, lan: string) {
        this.flag = flag
        this.title = title
        this.lan = lan
    }
}
