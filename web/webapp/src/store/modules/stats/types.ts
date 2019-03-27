import {Service, Node, Endpoint, Error} from "@/store/basic/types";


export class Stats {
    public started ?: number;
    public uptime ?: number;
    public memory ?: number;
    public threads ?: number;
    public gc_pause  ?: number;
    public counters ?: StatsCounter[];
}

export class StatsCounter {
    public timestamp?: number;
    public status_codes?: Map<string, number>;
    public total_reqs?: number;
}


export interface APIStatsState {
    loaded: boolean;
    services: Service[];
    currentNodeStats: Stats;
    pageLoading: boolean;
    xError: Error;
}


export interface ServicesStatsState {
    services: Service[];
    nodesStatsMap: Map<string, Stats>;
    // key: node.ip+node.port
    cardLoading: Map<string, boolean>;
    cardLoadingChanged: string;
    pageLoading: boolean;
    xError: string;
}


export function mergeNodes(services: Service[]) {

    let nodes: Node[] = [];

    services.forEach(item => {
        nodes.push(...item.nodes)
    })

    return nodes;
}
