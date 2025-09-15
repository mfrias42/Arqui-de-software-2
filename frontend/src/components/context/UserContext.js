import React, { createContext, useState, useEffect } from 'react';
import jwtDecode from 'jwt-decode';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Verificar si hay un token al cargar la aplicación
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decoded = jwtDecode(token);
                setUser({
                    ...decoded,
                    token
                });
            } catch (error) {
                // Si el token es inválido, limpiar el localStorage
                localStorage.removeItem('token');
                localStorage.removeItem('userId');
                localStorage.removeItem('usertype');
            }
        }
        setLoading(false);
    }, []);

    if (loading) return <div>Cargando...</div>;

    return (
        <UserContext.Provider value={{ user, setUser }}>
            {children}
        </UserContext.Provider>
    );
};
