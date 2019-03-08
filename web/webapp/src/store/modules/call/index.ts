import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {CallState} from './types';
import {getServiceDetails, postServiceRequest} from '@/api/call';


const namespaced: boolean = true;

const state: CallState = {
    requestLoading: false,
    services: [],
    requestResult: {}
}

const mutations: MutationTree<any> = {
    [TYPES.SET_CALL_SERVICES](state: CallState, {services}): void {
        state.services = services
        state.requestLoading = false
    },

    [TYPES.SET_CALL_LOADING](state: CallState, loading: boolean): void {
        state.requestLoading = loading
    },

    [TYPES.SET_CALL_RESULT](state: CallState, result: object): void {
        state.requestResult = result
    },

};

const actions: ActionTree<any, any> = {

    async getServiceDetails({commit}) {

        commit(TYPES.SET_CALL_LOADING, true);

        const res: Ajax.AjaxResponse = await getServiceDetails();
        commit(TYPES.SET_CALL_SERVICES, {
            services: res.data
        });
    },

    async postServiceRequest({commit}, req: object) {

        commit(TYPES.SET_CALL_LOADING, true);
        const res: Ajax.AjaxResponse = await postServiceRequest(req);
        commit(TYPES.SET_CALL_RESULT, [
            {
                "_id": "5c81c4b7d5c7407fbc758cb0",
                "index": 0,
                "guid": "a4653a44-41ff-4d4d-893a-7f42cfef5127",
                "isActive": false,
                "balance": "$3,840.43",
                "picture": "http://placehold.it/32x32",
                "age": 30,
                "eyeColor": "brown",
                "name": "Barrett Chapman",
                "gender": "male",
                "company": "STRALUM",
                "email": "barrettchapman@stralum.com",
                "phone": "+1 (851) 547-3062",
                "address": "105 Bartlett Place, Orovada, Palau, 8945",
                "about": "Qui dolore aute laborum cupidatat occaecat sint in nisi velit tempor eiusmod. Excepteur esse sunt elit tempor. Aute ea occaecat ipsum fugiat Lorem.\r\n",
                "registered": "2017-04-18T11:15:26 -08:00",
                "latitude": -45.141128,
                "longitude": 42.720251,
                "tags": [
                    "do",
                    "ipsum",
                    "minim",
                    "culpa",
                    "nisi",
                    "enim",
                    "cillum"
                ],
                "friends": [
                    {
                        "id": 0,
                        "name": "Delgado Koch"
                    },
                    {
                        "id": 1,
                        "name": "Jaclyn Blackburn"
                    },
                    {
                        "id": 2,
                        "name": "Shepard Guzman"
                    }
                ],
                "greeting": "Hello, Barrett Chapman! You have 2 unread messages.",
                "favoriteFruit": "apple"
            },
            {
                "_id": "5c81c4b75cdc8177af2871be",
                "index": 1,
                "guid": "4a920cf7-245a-4080-874b-810dc7635515",
                "isActive": true,
                "balance": "$2,402.17",
                "picture": "http://placehold.it/32x32",
                "age": 36,
                "eyeColor": "blue",
                "name": "Gwen Peterson",
                "gender": "female",
                "company": "TERRAGEN",
                "email": "gwenpeterson@terragen.com",
                "phone": "+1 (871) 440-3776",
                "address": "793 Stoddard Place, Vaughn, Montana, 8483",
                "about": "Quis cupidatat culpa dolore exercitation adipisicing elit quis ipsum dolore dolor consectetur. Enim magna et mollit cupidatat laborum cillum est consequat reprehenderit do occaecat elit. Fugiat cillum amet do labore nisi duis adipisicing anim fugiat id. Reprehenderit incididunt nisi ut id officia qui. Cupidatat minim duis do esse non in non reprehenderit voluptate ea tempor excepteur aliquip. Consectetur nisi veniam esse amet aliquip non voluptate quis sunt aliquip fugiat adipisicing id.\r\n",
                "registered": "2015-06-02T11:30:42 -08:00",
                "latitude": 52.187568,
                "longitude": 89.729685,
                "tags": [
                    "nulla",
                    "ea",
                    "culpa",
                    "in",
                    "ea",
                    "fugiat",
                    "elit"
                ],
                "friends": [
                    {
                        "id": 0,
                        "name": "John Suarez"
                    },
                    {
                        "id": 1,
                        "name": "Brandy Young"
                    },
                    {
                        "id": 2,
                        "name": "Kramer Travis"
                    }
                ],
                "greeting": "Hello, Gwen Peterson! You have 7 unread messages.",
                "favoriteFruit": "banana"
            },
            {
                "_id": "5c81c4b7d966b8890aefb857",
                "index": 2,
                "guid": "e71ab990-8573-4caf-a3b5-dc67d0898f08",
                "isActive": true,
                "balance": "$2,690.88",
                "picture": "http://placehold.it/32x32",
                "age": 24,
                "eyeColor": "brown",
                "name": "Reeves Collier",
                "gender": "male",
                "company": "DEEPENDS",
                "email": "reevescollier@deepends.com",
                "phone": "+1 (962) 549-2745",
                "address": "249 Hope Street, Sperryville, Arizona, 1039",
                "about": "Est irure irure laboris consequat consequat magna mollit veniam id laboris laboris. Ullamco mollit nisi officia commodo ea exercitation proident adipisicing consequat elit cupidatat eu cupidatat. Eu ullamco id do aliquip non laboris sint duis exercitation excepteur pariatur labore ullamco. Ullamco fugiat nostrud minim ipsum laboris consequat. Laboris dolor ad voluptate deserunt cupidatat mollit voluptate labore minim mollit nisi esse labore sit. Cillum anim excepteur sint velit in dolor excepteur reprehenderit magna consectetur commodo consequat reprehenderit fugiat.\r\n",
                "registered": "2015-10-16T06:12:57 -08:00",
                "latitude": -57.123692,
                "longitude": 45.323274,
                "tags": [
                    "deserunt",
                    "voluptate",
                    "reprehenderit",
                    "et",
                    "in",
                    "ex",
                    "quis"
                ],
                "friends": [
                    {
                        "id": 0,
                        "name": "Lily Stafford"
                    },
                    {
                        "id": 1,
                        "name": "Nunez Shepherd"
                    },
                    {
                        "id": 2,
                        "name": "Berta Forbes"
                    }
                ],
                "greeting": "Hello, Reeves Collier! You have 7 unread messages.",
                "favoriteFruit": "banana"
            },
            {
                "_id": "5c81c4b7a86d96fe6e481d76",
                "index": 3,
                "guid": "3041f41a-19d1-4375-915b-6b7e46cbdd21",
                "isActive": false,
                "balance": "$1,907.20",
                "picture": "http://placehold.it/32x32",
                "age": 34,
                "eyeColor": "green",
                "name": "Hammond Huber",
                "gender": "male",
                "company": "ZYTREK",
                "email": "hammondhuber@zytrek.com",
                "phone": "+1 (982) 475-2982",
                "address": "846 Seaview Court, Hoagland, Illinois, 6092",
                "about": "Excepteur tempor enim incididunt do et voluptate ipsum. Nulla laboris id exercitation do ullamco occaecat sit do amet amet aute. Cupidatat aute aliqua et fugiat aliquip consectetur eiusmod.\r\n",
                "registered": "2017-04-28T12:53:41 -08:00",
                "latitude": -43.410016,
                "longitude": 55.808627,
                "tags": [
                    "minim",
                    "enim",
                    "id",
                    "duis",
                    "consectetur",
                    "deserunt",
                    "duis"
                ],
                "friends": [
                    {
                        "id": 0,
                        "name": "Carly Hubbard"
                    },
                    {
                        "id": 1,
                        "name": "Alba Sweet"
                    },
                    {
                        "id": 2,
                        "name": "Rosa Farley"
                    }
                ],
                "greeting": "Hello, Hammond Huber! You have 7 unread messages.",
                "favoriteFruit": "apple"
            },
            {
                "_id": "5c81c4b72861e1fec0dd6c9b",
                "index": 4,
                "guid": "ab034132-b8df-4488-8356-4f58b46f0f94",
                "isActive": true,
                "balance": "$3,669.74",
                "picture": "http://placehold.it/32x32",
                "age": 24,
                "eyeColor": "green",
                "name": "Salazar Tucker",
                "gender": "male",
                "company": "QOT",
                "email": "salazartucker@qot.com",
                "phone": "+1 (937) 489-2363",
                "address": "750 Suydam Place, Chumuckla, Louisiana, 6735",
                "about": "Aliqua id elit nisi enim ut amet duis eiusmod consectetur fugiat incididunt eiusmod ea. Cillum proident nostrud fugiat Lorem sit elit ad veniam ullamco anim cillum consequat ea velit. Ad amet in ut aute qui officia amet. Dolor deserunt sint anim culpa mollit fugiat amet. Cupidatat pariatur nulla culpa commodo ad dolore do esse ullamco est adipisicing proident. Quis eiusmod id amet velit in aliquip officia amet labore. Cupidatat consectetur esse incididunt magna ad officia aute eiusmod deserunt sit deserunt exercitation.\r\n",
                "registered": "2016-06-26T02:15:30 -08:00",
                "latitude": 9.867729,
                "longitude": -113.984406,
                "tags": [
                    "commodo",
                    "velit",
                    "incididunt",
                    "nostrud",
                    "esse",
                    "elit",
                    "proident"
                ],
                "friends": [
                    {
                        "id": 0,
                        "name": "Clarissa Barrett"
                    },
                    {
                        "id": 1,
                        "name": "Kendra Riggs"
                    },
                    {
                        "id": 2,
                        "name": "Lacey Stuart"
                    }
                ],
                "greeting": "Hello, Salazar Tucker! You have 8 unread messages.",
                "favoriteFruit": "strawberry"
            }
        ]);
    },
};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
