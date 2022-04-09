import { Navigate, useRoutes } from 'react-router-dom';
import LogoOnlyLayout from '../layouts/LogoOnlyLayout';
import Register from "../pages/Register";
import Login from "../pages/Login";

export default function Router() {
    return useRoutes([
        {
            path: '/',
            element: <LogoOnlyLayout />,
            children: [
                { path: '/', element: <Login /> },
                { path: 'login', element: <Login /> },
                { path: 'register', element: <Register /> },
                { path: '*', element: <Navigate to="/login" /> }
            ]
        },
        { path: '*', element: <Navigate to="/login" /> }
    ]);
}
