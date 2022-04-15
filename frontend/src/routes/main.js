import {Navigate, useRoutes} from 'react-router-dom';

import AppLayout from './../layouts/dashboard';
import LogoOnlyLayout from './../layouts/LogoOnlyLayout';

import Status from './../pages/Status';
import StatusList from './../pages/Status/StatusList';
import Feeling from "../pages/Feeling";
import Profile from './../pages/Profile';
import Logout from './../pages/Logout';
import NotFound from './../pages/Page404';

export default function Router() {
    return useRoutes([
        {
            path: '/status',
            element: <AppLayout/>,
            children: [
                {path: 'create', element: <Status/>},
                {path: 'list', element: <StatusList/>}
            ]
        },
        {
            path: '/profile',
            element: <AppLayout/>,
            children: [
                {path: 'edit', element: <Profile/>}
            ]
        },
        {
            path: '/feeling',
            element: <AppLayout/>,
            children: [
                {path: 'frequency', element: <Feeling/>}
            ]
        },
        {
            path: '/',
            element: <LogoOnlyLayout/>,
            children: [
                {path: '', element: <Navigate to="/dashboard/app"/>},
                {path: '404', element: <NotFound/>},
                {path: 'logout', element: <Logout/>},
            ]
        },
        {path: '*', element: <Navigate to="/status/create" replace/>}
    ]);
}