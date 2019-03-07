export class Service {
    name!: string;
    version?: string;
    metadata?: string;
    endpoints?: Endpoint[];
    nodes?: Node[]
}

export class Value {
    name!: string;
    type?: string;
    values?: Value[];
}

export class Endpoint {
    name!: string;
    request?: Value;
    response?: Value;
    metadata?: Map<string, string>;
}

export class Node {
    id!: string;
    address?: string
    port?: number;
    metadata?: Map<string, string>;
}

export interface RegistryState {
    registries: Service[];
    pageLoading: boolean;
}