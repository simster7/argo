import * as H from 'history';
import * as React from 'react';
import {match} from 'react-router';
import {NavigationApi, NotificationsApi, PopupApi} from './components';

export interface AppContext {
    router: {
        history: H.History;
        route: {
            location: H.Location;
            match: match<any>;
        };
    };
    apis: {
        popup: PopupApi;
        notifications: NotificationsApi;
        navigation: NavigationApi;
    };
    history: H.History;
}

export interface ContextApis {
    popup: PopupApi;
    notifications: NotificationsApi;
    navigation: NavigationApi;
    history: H.History;
}

export const {Provider, Consumer} = React.createContext<ContextApis>(null);
