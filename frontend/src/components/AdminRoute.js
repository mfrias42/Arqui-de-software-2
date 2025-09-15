import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

const AdminRoute = () => {
    const token = localStorage.getItem('token');
    const userType = token ? JSON.parse(atob(token.split('.')[1])).user_type : null;

    if (!token) {
        return <Navigate to="/login" />;
    }

    if (userType !== 'administrador') {
        return <Navigate to="/home" />;
    }

    return <Outlet />;
};

export default AdminRoute;